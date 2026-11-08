package api

// This file reads the autonomy settings that govern the stages: each is
// auto or owner, in the repo's git config under rhiza.autonomy.<stage>,
// re-read every time it is used.

import (
	"strings"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

// The values of an autonomy setting.
const (
	AutonomyAuto  = "auto"
	AutonomyOwner = "owner"
)

// autonomyStage reads git config rhiza.autonomy.<stage>: auto or owner, the
// stage's default in settingDefs when unset, so the gate and the daemon act
// on the value rhiza daemon settings shows. Any other value is a usage error
// naming the key.
func autonomyStage(ctx *Ctx, stage string) (string, error) {
	key := "rhiza.autonomy." + stage
	d, ok := findSetting(key)
	if !ok {
		panic("autonomyStage: no setting " + key)
	}
	def := d.def
	if ctx.Git == nil {
		return def, nil
	}
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

// runAutonomy reads each stage's autonomy as a run of def starts with it:
// the workflow's declaration when it makes one, else the
// rhiza.autonomy.<stage> setting. The run's journal records it.
func runAutonomy(ctx *Ctx, def *workflow.Definition) ([]run.StageAutonomy, error) {
	out := make([]run.StageAutonomy, 0, len(workflow.AutonomyStages))
	for _, st := range workflow.AutonomyStages {
		if m, ok := def.Autonomy(st); ok {
			out = append(out, run.StageAutonomy{Stage: st, Mode: m, Source: run.SourceWorkflow, Workflow: def.Name})
			continue
		}
		m, err := autonomyStage(ctx, st)
		if err != nil {
			return nil, err
		}
		out = append(out, run.StageAutonomy{Stage: st, Mode: m, Source: run.SourceSetting})
	}
	return out, nil
}
