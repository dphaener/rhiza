package api

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/genesis/internal/canon"
)

// Humaner is implemented by every response type: Human writes the response
// as markdown for people. Callers print it unless JSON was asked for, in
// which case they print canon.Marshal(response) instead.
type Humaner interface {
	Human(w io.Writer)
}

// DescribeRequest asks for one operation, or every one if Op is empty.
type DescribeRequest struct {
	Op string `json:"op,omitempty" arg:"0" desc:"The operation to describe, such as change.new; every operation if empty"`
}

// DescribeResponse lists operations with their schemas.
type DescribeResponse struct {
	Operations []OperationInfo `json:"operations" desc:"The operations, sorted by name"`
}

// OperationInfo describes one operation.
type OperationInfo struct {
	Name     string          `json:"name" desc:"The operation name, its command path with dots"`
	Summary  string          `json:"summary" desc:"One line describing the operation"`
	Write    bool            `json:"write" desc:"Whether the operation appends to the log"`
	Request  json.RawMessage `json:"request" desc:"The JSON Schema of the request"`
	Response json.RawMessage `json:"response" desc:"The JSON Schema of the response"`
	Input    json.RawMessage `json:"input,omitempty" desc:"The JSON Schema of a file or stdin document the operation reads, if any"`
}

func init() {
	register(Operation{
		Name:     "describe",
		Summary:  "Describe operations: their arguments and request and response schemas",
		Request:  reflect.TypeFor[DescribeRequest](),
		Response: reflect.TypeFor[DescribeResponse](),
		Run:      runDescribe,
	})
}

func runDescribe(_ *Ctx, req any) (any, error) {
	r := req.(DescribeRequest)
	ops := Operations()
	if r.Op != "" {
		o, ok := Lookup(r.Op)
		if !ok {
			return nil, NotFoundf("unknown operation %q", r.Op).WithNext("genesis describe")
		}
		ops = []Operation{o}
	}
	resp := DescribeResponse{Operations: make([]OperationInfo, 0, len(ops))}
	for _, o := range ops {
		reqSchema, err := canon.Marshal(Schema(o.Request))
		if err != nil {
			return nil, fmt.Errorf("describe %s: %w", o.Name, err)
		}
		respSchema, err := canon.Marshal(Schema(o.Response))
		if err != nil {
			return nil, fmt.Errorf("describe %s: %w", o.Name, err)
		}
		info := OperationInfo{
			Name: o.Name, Summary: o.Summary, Write: o.Write,
			Request: reqSchema, Response: respSchema,
		}
		if o.Input != nil {
			if info.Input, err = canon.Marshal(o.Input.Schema()); err != nil {
				return nil, fmt.Errorf("describe %s: %w", o.Name, err)
			}
		}
		resp.Operations = append(resp.Operations, info)
	}
	return resp, nil
}

// Human lists each operation with its synopsis and parameters.
func (r DescribeResponse) Human(w io.Writer) {
	for i, info := range r.Operations {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "## %s\n\n%s\n", info.Name, info.Summary)
		o, ok := Lookup(info.Name)
		if !ok {
			continue
		}
		fmt.Fprintf(w, "\n    %s\n", UsageLine(o))
		if o.Input != nil {
			fmt.Fprintf(w, "\nInput: %s. Its schema: genesis describe %s --json\n", o.Input.Summary, info.Name)
		}
		ps := Params(o.Request)
		if len(ps) == 0 {
			continue
		}
		fmt.Fprintln(w)
		for _, p := range ps {
			name := "--" + p.Flag
			if p.Arg >= 0 {
				name = "<" + p.Name + ">"
			}
			var notes []string
			if p.Required && p.Kind != ParamBool {
				notes = append(notes, "required")
			}
			if p.Kind == ParamList && p.Arg < 0 {
				notes = append(notes, "repeatable")
			}
			if len(p.Enum) > 0 {
				notes = append(notes, "one of "+strings.Join(p.Enum, ", "))
			}
			line := "- `" + name + "`"
			if p.Desc != "" {
				line += ": " + p.Desc
			}
			if len(notes) > 0 {
				line += " (" + strings.Join(notes, "; ") + ")"
			}
			fmt.Fprintln(w, line)
		}
	}
}
