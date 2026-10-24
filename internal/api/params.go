package api

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

// ParamKind is the shape of a command-line parameter.
type ParamKind int

// The parameter kinds.
const (
	ParamString ParamKind = iota // one value
	ParamBool                    // present or not; --x=false also works
	ParamInt                     // one integer value
	ParamList                    // a repeated flag, or the rest of the positionals
)

// Param is one command-line parameter of a request type: a positional
// argument or a flag, mapped to one JSON field of the request.
//
// The struct-tag contract for request types:
//
//   - json:"name" is the wire name; Name is that name.
//   - arg:"N" makes the field positional argument N (0-based, contiguous).
//     A []string positional takes the rest of the arguments and must be
//     the last one.
//   - Any other field is a flag. flag:"name" names it; the default is the
//     json name with "_" turned into "-". flag:"-" keeps the field off the
//     command line (it is still in the JSON request).
//   - desc:"..." describes the field; enum:"a,b" lists its allowed values.
//   - A field without omitempty or omitzero is required.
//
// Fields may be strings (named string types too), bools, ints or []string.
type Param struct {
	// Name is the field's JSON name.
	Name string
	// Arg is the positional index, or -1 for a flag.
	Arg int
	// Rest marks a []string positional that takes every remaining argument.
	Rest bool
	// Flag is the flag name without dashes; empty for a positional.
	Flag string
	// Kind is the parameter's shape.
	Kind ParamKind
	// Required reports a field without omitempty or omitzero.
	Required bool
	// Desc and Enum come from the desc and enum tags.
	Desc string
	Enum []string
}

// Params returns the command-line parameters of a request struct type:
// positionals in order, then flags in field order. It panics on a request
// type that breaks the tag contract (see Param), since request types are
// fixed by the code.
func Params(t reflect.Type) []Param {
	var args, flags []Param
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		tag, ok := f.Tag.Lookup("json")
		if !ok || tag == "-" {
			continue
		}
		name, opts, _ := strings.Cut(tag, ",")
		if f.Tag.Get("flag") == "-" {
			continue
		}
		p := Param{
			Name: name,
			Arg:  -1,
			Kind: paramKind(t, f),
			Desc: f.Tag.Get("desc"),
			Required: !slices.ContainsFunc(strings.Split(opts, ","), func(o string) bool {
				return o == "omitempty" || o == "omitzero"
			}),
		}
		if e, ok := f.Tag.Lookup("enum"); ok {
			p.Enum = strings.Split(e, ",")
		}
		if a, ok := f.Tag.Lookup("arg"); ok {
			n, err := strconv.Atoi(a)
			if err != nil || n < 0 {
				panic(fmt.Sprintf("api: params: %s.%s: bad arg tag %q", t, f.Name, a))
			}
			p.Arg = n
			p.Rest = p.Kind == ParamList
			args = append(args, p)
			continue
		}
		p.Flag = f.Tag.Get("flag")
		if p.Flag == "" {
			p.Flag = strings.ReplaceAll(name, "_", "-")
		}
		flags = append(flags, p)
	}
	slices.SortStableFunc(args, func(a, b Param) int { return a.Arg - b.Arg })
	for i, p := range args {
		if p.Arg != i {
			panic(fmt.Sprintf("api: params: %s: positional args are not 0..%d", t, len(args)-1))
		}
		if p.Rest && i != len(args)-1 {
			panic(fmt.Sprintf("api: params: %s: list positional %q is not last", t, p.Name))
		}
	}
	return append(args, flags...)
}

func paramKind(t reflect.Type, f reflect.StructField) ParamKind {
	switch f.Type.Kind() {
	case reflect.String:
		return ParamString
	case reflect.Bool:
		return ParamBool
	case reflect.Int, reflect.Int64:
		return ParamInt
	case reflect.Slice:
		if f.Type.Elem().Kind() == reflect.String {
			return ParamList
		}
	}
	panic(fmt.Sprintf("api: params: %s.%s: type %s cannot be a command-line parameter (tag it flag:\"-\")", t, f.Name, f.Type))
}

// UsageLine returns the command-line synopsis of o, such as
// "genesis change new <intent> [--kind <kind>]": positionals, then the
// required flags, then the optional ones.
func UsageLine(o Operation) string {
	parts := []string{"genesis", CommandPath(o.Name)}
	ps := Params(o.Request)
	var required, optional []string
	for _, p := range ps {
		switch {
		case p.Arg >= 0 && p.Rest && p.Required:
			parts = append(parts, "<"+p.Name+">...")
		case p.Arg >= 0 && p.Rest:
			parts = append(parts, "[<"+p.Name+">...]")
		case p.Arg >= 0 && p.Required:
			parts = append(parts, "<"+p.Name+">")
		case p.Arg >= 0:
			parts = append(parts, "[<"+p.Name+">]")
		default:
			s := "--" + p.Flag
			if p.Kind != ParamBool {
				s += " <" + p.Flag + ">"
			}
			if p.Required && p.Kind != ParamBool {
				required = append(required, s)
				continue
			}
			s = "[" + s + "]"
			if p.Kind == ParamList {
				s += "..."
			}
			optional = append(optional, s)
		}
	}
	parts = append(parts, required...)
	parts = append(parts, optional...)
	return strings.Join(parts, " ")
}

// CommandPath turns an operation name into its command words:
// "change.new" becomes "change new".
func CommandPath(name string) string {
	return strings.ReplaceAll(name, ".", " ")
}
