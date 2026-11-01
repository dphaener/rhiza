package index

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// modelSentinels name the errors model problems wrap; each wraps exactly
// one.
var modelSentinels = []sentinel{
	{model.ErrConflict, "conflict"},
	{model.ErrMissing, "missing"},
	{model.ErrBase, "base"},
	{model.ErrExists, "exists"},
	{model.ErrRetired, "retired"},
	{model.ErrParent, "parent"},
	{model.ErrCycle, "cycle"},
	{model.ErrHandle, "handle"},
	{model.ErrState, "state"},
}

func modelSentinel(err error) string { return nameOf(err, modelSentinels) }

// writeViews recomputes the views asked for over sets and writes each by
// difference, adding the rows written per view to writes.
func writeViews(tx *sql.Tx, sets []changeset.ChangeSet, asIs, toBe bool, writes map[string]int) error {
	if asIs {
		m, ps, st := model.AsIsState(sets)
		if err := writeAsIsState(tx, sets, m, st); err != nil {
			return err
		}
		n, err := writeView(tx, viewAsIs, m, ps, &st)
		if err != nil {
			return err
		}
		if n > 0 {
			writes[viewAsIs] += n
		}
	}
	if toBe {
		m, ps := model.ToBe(sets)
		n, err := writeView(tx, viewToBe, m, ps, nil)
		if err != nil {
			return err
		}
		if n > 0 {
			writes[viewToBe] += n
		}
	}
	return nil
}

// writeView makes the stored rows of view equal m and its problems equal
// ps: rows whose hash changed are upserted, rows not in m deleted. It
// returns the number of node and link rows written or deleted.
func writeView(tx *sql.Tx, view string, m model.Model, ps []model.Problem, st *model.State) (int, error) {
	written := 0

	nodes := m.Nodes()
	rows := make([]row, len(nodes))
	for i, n := range nodes {
		rows[i] = nodeRow(n, st)
	}
	n, err := diffRows(tx, "nodes", nodeColumns, view, rows)
	if err != nil {
		return 0, err
	}
	written += n

	if err := syncFTS(tx, view, nodes); err != nil {
		return 0, err
	}

	links := m.Links()
	rows = make([]row, len(links))
	for i, l := range links {
		rows[i] = linkRow(l, st)
	}
	if n, err = diffRows(tx, "links", linkColumns, view, rows); err != nil {
		return 0, err
	}
	written += n

	out := make([]storedProblem, 0, len(ps))
	for _, p := range ps {
		out = append(out, storedProblem{ChangeSet: idText(p.ChangeSet), Op: idText(p.Op),
			Sentinel: modelSentinel(p.Err), Message: p.Err.Error()})
	}
	return written, replaceProblems(tx, view, out)
}

// row is one node or link as stored: its ID, its column values after view
// and id, and the hash of those values.
type row struct {
	id   string
	vals []any
	hash []byte
}

// The columns after view and id, in the order row.vals holds them.
var (
	nodeColumns = []string{"kind", "handle", "parent", "key", "fields", "retired", "reason", "version", "created", "history", "aspects", "placed"}
	linkColumns = []string{"type", "from_id", "to_id", "note", "decision", "retired", "version", "created", "history", "aspects"}
)

// newRow builds a row and hashes its values. Values are strings, ints or
// nil (SQL NULL).
func newRow(rid string, vals ...any) row {
	h := sha256.New()
	for _, v := range vals {
		switch v := v.(type) {
		case nil:
			h.Write([]byte{0})
		case string:
			fmt.Fprintf(h, "s%d:%s", len(v), v)
		case int:
			fmt.Fprintf(h, "i%d;", v)
		}
	}
	return row{id: rid, vals: vals, hash: h.Sum(nil)}
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// nodeRow is n as stored.
// st, when not nil, is the state of the view, whose aspects and placement
// the row keeps.
func nodeRow(n model.Node, st *model.State) row {
	var aspects, placed any
	if st != nil {
		aspects = encodeAspects(st.Aspects[n.ID])
		m := st.Placed[n.ID]
		placed = fmt.Sprintf("%d %s %s", m.Seq, idText(m.Ref.ChangeSet), idText(m.Ref.Op))
	}
	return newRow(idText(n.ID), string(n.Kind), n.Handle, idText(n.Parent), n.Key, encodeFields(n.Fields),
		boolInt(n.Retired), n.Reason, idText(n.Version), idText(n.Created), encodeHistory(n.History), aspects, placed)
}

// linkRow is l as stored.
func linkRow(l model.Link, st *model.State) row {
	var aspects any
	if st != nil {
		aspects = encodeAspects(st.Aspects[l.ID])
	}
	return newRow(idText(l.ID), string(l.Type), idText(l.From), idText(l.To), l.Note, idText(l.Decision),
		boolInt(l.Retired), idText(l.Version), idText(l.Created), encodeHistory(l.History), aspects)
}

// encodeFields spells a field map as a JSON object with sorted keys and
// each value written verbatim, so decoding gives back the same bytes; nil
// (SQL NULL) when empty.
func encodeFields(m map[string]json.RawMessage) any {
	if len(m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	var b bytes.Buffer
	b.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			b.WriteByte(',')
		}
		kb, _ := json.Marshal(k)
		b.Write(kb)
		b.WriteByte(':')
		b.Write(m[k])
	}
	b.WriteByte('}')
	return b.String()
}

// decodeFields reverses encodeFields.
func decodeFields(s sql.NullString) (map[string]json.RawMessage, error) {
	if !s.Valid {
		return nil, nil
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s.String), &m); err != nil {
		return nil, fmt.Errorf("stored fields: %w", err)
	}
	if len(m) == 0 {
		return nil, nil
	}
	return m, nil
}

// encodeHistory spells a history as a JSON array of [change set, op]
// pairs; nil (SQL NULL) when empty.
func encodeHistory(h []model.Ref) any {
	if len(h) == 0 {
		return nil
	}
	pairs := make([][2]string, len(h))
	for i, r := range h {
		pairs[i] = [2]string{idText(r.ChangeSet), idText(r.Op)}
	}
	b, _ := json.Marshal(pairs)
	return string(b)
}

// decodeHistory reverses encodeHistory.
func decodeHistory(s sql.NullString) ([]model.Ref, error) {
	if !s.Valid {
		return nil, nil
	}
	var pairs [][2]string
	if err := json.Unmarshal([]byte(s.String), &pairs); err != nil {
		return nil, fmt.Errorf("stored history: %w", err)
	}
	if len(pairs) == 0 {
		return nil, nil
	}
	out := make([]model.Ref, len(pairs))
	for i, p := range pairs {
		out[i] = model.Ref{ChangeSet: parseID(p[0]), Op: parseID(p[1])}
	}
	return out, nil
}

// diffRows makes table's rows for view equal rows, by hash.
func diffRows(tx *sql.Tx, table string, cols []string, view string, rows []row) (int, error) {
	have := map[string][]byte{}
	q, err := tx.Query(`SELECT id, hash FROM `+table+` WHERE view = ?`, view)
	if err != nil {
		return 0, err
	}
	for q.Next() {
		var rid string
		var h []byte
		if err := q.Scan(&rid, &h); err != nil {
			q.Close()
			return 0, err
		}
		have[rid] = h
	}
	q.Close()
	if err := q.Err(); err != nil {
		return 0, err
	}

	written := 0
	var upsert *sql.Stmt
	for _, r := range rows {
		h, ok := have[r.id]
		delete(have, r.id)
		if ok && bytes.Equal(h, r.hash) {
			continue
		}
		if upsert == nil {
			marks := "?, ?, ?"
			for range cols {
				marks += ", ?"
			}
			upsert, err = tx.Prepare(fmt.Sprintf(`INSERT OR REPLACE INTO %s (view, id, hash, %s) VALUES (%s)`,
				table, strings.Join(cols, ", "), marks))
			if err != nil {
				return 0, err
			}
			defer upsert.Close()
		}
		args := append([]any{view, r.id, r.hash}, r.vals...)
		if _, err := upsert.Exec(args...); err != nil {
			return 0, err
		}
		written++
	}
	if len(have) > 0 {
		del, err := tx.Prepare(`DELETE FROM ` + table + ` WHERE view = ? AND id = ?`)
		if err != nil {
			return 0, err
		}
		defer del.Close()
		for rid := range have {
			if _, err := del.Exec(view, rid); err != nil {
				return 0, err
			}
			written++
		}
	}
	return written, nil
}

// nodeSelect selects the node columns scanNode reads, from nodes.
const nodeSelect = `SELECT id, kind, handle, parent, key, fields, retired, reason, version, created, history FROM nodes`

// linkSelect selects the link columns scanLink reads, from links.
const linkSelect = `SELECT id, type, from_id, to_id, note, decision, retired, version, created, history FROM links`

// scanner is what *sql.Row and *sql.Rows share.
type scanner interface {
	Scan(dest ...any) error
}

// scanNode reads one row of nodeSelect back into a model.Node.
func scanNode(s scanner, extra ...any) (model.Node, error) {
	var (
		n                                   model.Node
		nid, kind, parent, version, created string
		fields, history                     sql.NullString
		retired                             int
	)
	if err := s.Scan(append([]any{&nid, &kind, &n.Handle, &parent, &n.Key, &fields, &retired, &n.Reason, &version, &created, &history}, extra...)...); err != nil {
		return model.Node{}, err
	}
	var err error
	if n.Fields, err = decodeFields(fields); err != nil {
		return model.Node{}, err
	}
	if n.History, err = decodeHistory(history); err != nil {
		return model.Node{}, err
	}
	n.ID, n.Kind, n.Parent = parseID(nid), op.Kind(kind), parseID(parent)
	n.Retired, n.Version, n.Created = retired != 0, parseID(version), parseID(created)
	return n, nil
}

// scanLink reads one row of linkSelect back into a model.Link.
func scanLink(s scanner, extra ...any) (model.Link, error) {
	var (
		l                                              model.Link
		lid, typ, from, to, decision, version, created string
		history                                        sql.NullString
		retired                                        int
	)
	if err := s.Scan(append([]any{&lid, &typ, &from, &to, &l.Note, &decision, &retired, &version, &created, &history}, extra...)...); err != nil {
		return model.Link{}, err
	}
	var err error
	if l.History, err = decodeHistory(history); err != nil {
		return model.Link{}, err
	}
	l.ID, l.Type, l.From, l.To, l.Decision = parseID(lid), op.LinkType(typ), parseID(from), parseID(to), parseID(decision)
	l.Retired, l.Version, l.Created = retired != 0, parseID(version), parseID(created)
	return l, nil
}
