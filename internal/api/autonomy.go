package api

// This file reads the autonomy settings that govern the stages: each is
// auto or owner, in the repo's git config under rhiza.autonomy.<stage>,
// re-read every time it is used.

import (
	"strings"
)

// The values of an autonomy setting.
const (
	AutonomyAuto  = "auto"
	AutonomyOwner = "owner"
)

// autonomyStage reads git config rhiza.autonomy.<stage>: auto or owner, def
// when unset. Any other value is a usage error naming the key.
func autonomyStage(ctx *Ctx, stage, def string) (string, error) {
	if ctx.Git == nil {
		return def, nil
	}
	key := "rhiza.autonomy." + stage
	out, err := ctx.Git("config", "--get", key)
	if err != nil {
		return def, nil // unset
	}
	switch v := strings.TrimSpace(out); v {
	case "":
		return def, nil
	case AutonomyAuto, AutonomyOwner:
		return v, nil
	default:
		return "", Usagef("git config %s is %q: it must be auto or owner", key, v)
	}
}
