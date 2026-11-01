package index

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"unicode"

	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// Tokens splits s into the words text search matches: maximal runs of
// letters and digits, lower-cased. It is the rule of the FTS table's
// tokenizer (unicode61 without diacritic removal), so a scan in memory and
// the table agree.
func Tokens(s string) []string {
	return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

// FieldText is the text of one field value: a JSON string unquoted, any
// other value as written.
func FieldText(raw json.RawMessage) string {
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return string(raw)
}

// FieldNames returns n's field names, sorted.
func FieldNames(n model.Node) []string {
	names := make([]string, 0, len(n.Fields))
	for k := range n.Fields {
		names = append(names, k)
	}
	slices.Sort(names)
	return names
}

// ftsBody is the text of n's fields, in field-name order, one per line.
func ftsBody(n model.Node) string {
	var b strings.Builder
	for _, k := range FieldNames(n) {
		b.WriteString(FieldText(n.Fields[k]))
		b.WriteByte('\n')
	}
	return b.String()
}

// syncFTS makes the nodes_fts rows of view equal the live nodes: rows
// whose text changed are rewritten, rows of nodes gone or retired deleted.
func syncFTS(tx *sql.Tx, view string, nodes []model.Node) error {
	have := map[string]string{}
	q, err := tx.Query(`SELECT id, sig FROM nodes_fts WHERE view = ?`, view)
	if err != nil {
		return err
	}
	for q.Next() {
		var i, sig string
		if err := q.Scan(&i, &sig); err != nil {
			q.Close()
			return err
		}
		have[i] = sig
	}
	q.Close()
	if err := q.Err(); err != nil {
		return err
	}
	for _, n := range nodes {
		if n.Retired {
			continue
		}
		i, body := idText(n.ID), ftsBody(n)
		sum := sha256.Sum256([]byte(n.Handle + "\x00" + body))
		sig := hex.EncodeToString(sum[:])
		old, ok := have[i]
		delete(have, i)
		if ok && old == sig {
			continue
		}
		if ok {
			if _, err := tx.Exec(`DELETE FROM nodes_fts WHERE view = ? AND id = ?`, view, i); err != nil {
				return err
			}
		}
		if _, err := tx.Exec(`INSERT INTO nodes_fts (view, id, sig, handle, body) VALUES (?, ?, ?, ?, ?)`,
			view, i, sig, n.Handle, body); err != nil {
			return err
		}
	}
	for i := range have {
		if _, err := tx.Exec(`DELETE FROM nodes_fts WHERE view = ? AND id = ?`, view, i); err != nil {
			return err
		}
	}
	return nil
}

// Search returns the live nodes of view v, of kinds when given, whose
// handle or field text holds every one of terms (lower-case Tokens), each
// as a word or the start of one, sorted by ID. It reads the FTS table, so
// a caller that needs exact matching filters the result by Tokens; the
// table can only hold more than that. No terms is an error.
func (x *Index) Search(v View, terms []string, kinds []op.Kind) ([]model.Node, error) {
	if err := v.check(); err != nil {
		return nil, err
	}
	if len(terms) == 0 {
		return nil, fmt.Errorf("index: search needs a term")
	}
	match := make([]string, len(terms))
	for i, t := range terms {
		match[i] = `"` + strings.ReplaceAll(t, `"`, `""`) + `"*`
	}
	query := `SELECT n.id, n.kind, n.handle, n.parent, n.key, n.fields, n.retired, n.reason, n.version, n.created, n.history
		FROM nodes_fts f JOIN nodes n ON n.view = f.view AND n.id = f.id
		WHERE f.view = ? AND nodes_fts MATCH ? AND n.retired = 0`
	args := []any{string(v), strings.Join(match, " ")}
	if len(kinds) > 0 {
		marks := make([]string, len(kinds))
		for i, k := range kinds {
			marks[i] = "?"
			args = append(args, string(k))
		}
		query += ` AND n.kind IN (` + strings.Join(marks, ", ") + `)`
	}
	return x.nodes(query+` ORDER BY n.id`, args...)
}
