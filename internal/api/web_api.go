package api

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/canon"
)

// ReadOnlyHandler serves GET /api/<op> for the web UI: it runs the operation
// the registry marks ReadOnly, with its request built from the query string
// (RequestFromQuery), and answers the bytes `rhiza <cmd> --json` prints.
// Every operation without the marker answers 404, as an unknown one does,
// and every other method 405, so no write is reachable. base is the Ctx the
// daemon runs with; each call gets a copy cut down to a read.
func ReadOnlyHandler(base *Ctx) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			writeAPIError(w, http.StatusMethodNotAllowed, Usagef("only GET is allowed"))
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/api/")
		o, ok := Lookup(name)
		if !ok || !o.ReadOnly {
			writeAPIError(w, http.StatusNotFound, NotFoundf("no read-only operation %q", name))
			return
		}
		req, err := RequestFromQuery(o, r.URL.Query())
		if err != nil {
			writeAPIError(w, http.StatusBadRequest, err)
			return
		}
		ctx := *base
		ctx.Record, ctx.DryRun, ctx.Out, ctx.JSON = false, false, nil, true
		ctx.Stdin, ctx.Rex, ctx.Start, ctx.Detach, ctx.Interrupts = nil, nil, nil, nil, nil
		resp, err := o.Run(&ctx, req)
		if err != nil {
			writeAPIError(w, apiStatus(err), err)
			return
		}
		b, err := canon.Marshal(resp)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
}

func apiStatus(err error) int {
	switch AsError(err).Code {
	case CodeUsage:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func writeAPIError(w http.ResponseWriter, status int, err error) {
	b, merr := canon.Marshal(ErrorJSON(err))
	if merr != nil {
		http.Error(w, err.Error(), status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

// RequestFromQuery builds o's request from a query string: each key is a
// request field's JSON name; a list field repeats its key. It goes through
// JSON and DecodeRequest like every other caller, so an unknown key, a
// repeated scalar or a wrong type is a usage error.
func RequestFromQuery(o Operation, q url.Values) (any, error) {
	byName := map[string]Param{}
	for _, p := range Params(o.Request) {
		byName[p.Name] = p
	}
	obj := map[string]any{}
	for k, vs := range q {
		p, ok := byName[k]
		if !ok {
			return nil, Usagef("%s: unknown parameter %q", o.Name, k).WithNext("rhiza describe " + o.Name)
		}
		if p.Kind == ParamList {
			obj[k] = vs
			continue
		}
		if len(vs) != 1 {
			return nil, Usagef("%s: parameter %q given more than once", o.Name, k)
		}
		switch p.Kind {
		case ParamBool:
			b, err := strconv.ParseBool(vs[0])
			if err != nil {
				return nil, Usagef("%s: parameter %q takes true or false, not %q", o.Name, k, vs[0])
			}
			obj[k] = b
		case ParamInt:
			n, err := strconv.ParseInt(vs[0], 10, 64)
			if err != nil {
				return nil, Usagef("%s: parameter %q takes an integer, not %q", o.Name, k, vs[0])
			}
			obj[k] = n
		default:
			obj[k] = vs[0]
		}
	}
	data, err := canon.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return DecodeRequest(o, data)
}
