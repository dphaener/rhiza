package main

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/api"
)

// globals are the flags accepted anywhere on the command line.
type globals struct {
	json, dryRun, record, help bool
	// change, actor and forRef are the raw refs of --change, --actor and
	// --for; the *Set fields record that the flag was given, so it wins
	// over the environment even when empty.
	change, actor, forRef       string
	changeSet, actorSet, forSet bool
}

// parseGlobals takes the global flags out of args, wherever they are, and
// returns the rest in order. Everything from "--" on is left alone.
func parseGlobals(args []string) (globals, []string, error) {
	var g globals
	var rest []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			rest = append(rest, args[i:]...)
			break
		}
		if a == "-h" {
			g.help = true
			continue
		}
		if !strings.HasPrefix(a, "--") {
			rest = append(rest, a)
			continue
		}
		name, val, hasVal := strings.Cut(a[2:], "=")
		var b *bool
		var s *string
		var set *bool
		switch name {
		case "json":
			b = &g.json
		case "dry-run":
			b = &g.dryRun
		case "record":
			b = &g.record
		case "help":
			b = &g.help
		case "change":
			s, set = &g.change, &g.changeSet
		case "actor":
			s, set = &g.actor, &g.actorSet
		case "for":
			s, set = &g.forRef, &g.forSet
		default:
			rest = append(rest, a)
			continue
		}
		if b != nil {
			v, err := boolValue(name, val, hasVal)
			if err != nil {
				return g, nil, err
			}
			*b = v
			continue
		}
		if !hasVal {
			if i+1 >= len(args) {
				return g, nil, api.Usagef("flag --%s needs a value", name)
			}
			i++
			val = args[i]
		}
		*s, *set = val, true
	}
	return g, rest, nil
}

func boolValue(name, val string, hasVal bool) (bool, error) {
	if !hasVal {
		return true, nil
	}
	v, err := strconv.ParseBool(val)
	if err != nil {
		return false, api.Usagef("flag --%s takes true or false, not %q", name, val)
	}
	return v, nil
}

// commandWords returns the leading words of args that are not flags: the
// candidates for a command path.
func commandWords(args []string) []string {
	n := 0
	for n < len(args) && !strings.HasPrefix(args[n], "-") {
		n++
	}
	return args[:n]
}

// matchCommand finds the operation named by the longest prefix of words
// joined with dots, and how many words it took.
func matchCommand(words []string) (api.Operation, int, bool) {
	for n := len(words); n > 0; n-- {
		if o, ok := api.Lookup(strings.Join(words[:n], ".")); ok {
			return o, n, true
		}
	}
	return api.Operation{}, 0, false
}

// subcommands finds the longest prefix of words that names a group of
// operations rather than one, such as "change" for change.ls and
// change.new. It returns that command path and the group's next words,
// sorted and deduplicated; none is false.
func subcommands(words []string) (string, []string, bool) {
	for n := len(words); n > 0; n-- {
		prefix := strings.Join(words[:n], ".") + "."
		var subs []string
		for _, o := range api.Operations() {
			if rest, ok := strings.CutPrefix(o.Name, prefix); ok {
				sub, _, _ := strings.Cut(rest, ".")
				subs = append(subs, sub)
			}
		}
		if len(subs) > 0 {
			slices.Sort(subs)
			return strings.Join(words[:n], " "), slices.Compact(subs), true
		}
	}
	return "", nil, false
}

// requestFromArgs builds o's request from the arguments after its command
// path, through the request type's arg and flag tags (see api.Param). It
// goes through JSON and api.DecodeRequest, so the command line can express
// nothing `api call` cannot. Every failure is a usage error.
func requestFromArgs(o api.Operation, args []string) (any, error) {
	ps := api.Params(o.Request)
	flags := map[string]api.Param{}
	var positionals []api.Param
	for _, p := range ps {
		if p.Arg >= 0 {
			positionals = append(positionals, p)
		} else {
			flags[p.Flag] = p
		}
	}
	next := "rhiza describe " + o.Name
	usage := func(format string, a ...any) error {
		return api.Usagef("%s: "+format, append([]any{api.CommandPath(o.Name)}, a...)...).WithNext(next)
	}

	obj := map[string]any{}
	var values []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			values = append(values, args[i+1:]...)
			break
		}
		if a == "-" || !strings.HasPrefix(a, "-") {
			values = append(values, a)
			continue
		}
		if !strings.HasPrefix(a, "--") {
			return nil, usage("unknown flag %s", a)
		}
		name, val, hasVal := strings.Cut(a[2:], "=")
		p, ok := flags[name]
		if !ok {
			return nil, usage("unknown flag --%s", name)
		}
		if p.Kind == api.ParamBool {
			v, err := boolValue(name, val, hasVal)
			if err != nil {
				return nil, usage("%v", err)
			}
			obj[p.Name] = v
			continue
		}
		if !hasVal {
			if i+1 >= len(args) {
				return nil, usage("flag --%s needs a value", name)
			}
			i++
			val = args[i]
		}
		switch p.Kind {
		case api.ParamList:
			l, _ := obj[p.Name].([]string)
			obj[p.Name] = append(l, val)
			continue
		}
		if _, ok := obj[p.Name]; ok {
			return nil, usage("flag --%s given more than once", name)
		}
		if p.Kind == api.ParamInt {
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, usage("flag --%s takes an integer, not %q", name, val)
			}
			obj[p.Name] = n
			continue
		}
		obj[p.Name] = val
	}

	for i, p := range positionals {
		switch {
		case p.Rest && i < len(values):
			obj[p.Name] = values[i:]
			values = values[:i]
		case i < len(values):
			obj[p.Name] = values[i]
		}
	}
	if len(values) > len(positionals) {
		return nil, usage("too many arguments: %q", values[len(positionals):])
	}
	for _, p := range ps {
		if _, ok := obj[p.Name]; ok || !p.Required || p.Kind == api.ParamBool {
			continue
		}
		if p.Arg >= 0 {
			return nil, usage("missing <%s>", p.Name)
		}
		return nil, usage("missing --%s", p.Flag)
	}

	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return api.DecodeRequest(o, data)
}
