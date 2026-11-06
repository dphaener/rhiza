package daemon

// ProjectPick is one project's turn in a round-robin pick: its name, the
// change sets it could start now in its own roadmap order (what Pick would
// return for it with no slot limit), how many of its runs are already
// running and its optional per-project cap on running runs (0 for none).
type ProjectPick struct {
	Project string
	Ready   []string
	Running int
	Max     int
}

// Picked is a change set the root daemon starts, and the project it is in.
type Picked struct {
	Project string
	Change  string
}

// RoundRobin fills the machine-wide maxRuns slots across projects. running
// is how many runs are already running on the machine, all projects
// together. Each turn goes to the next project in the order given that has
// a ready change set and is under its own Max, and takes that project's
// next ready change set; a project with nothing to take is passed over.
// It stops when the slots are full or no project can take a turn.
func RoundRobin(projects []ProjectPick, running, maxRuns int) []Picked {
	next := make([]int, len(projects))
	taken := make([]int, len(projects))
	var out []Picked
	for running < maxRuns {
		progress := false
		for i, p := range projects {
			if running >= maxRuns {
				break
			}
			if next[i] >= len(p.Ready) {
				continue
			}
			if p.Max > 0 && p.Running+taken[i] >= p.Max {
				continue
			}
			out = append(out, Picked{Project: p.Project, Change: p.Ready[next[i]]})
			next[i]++
			taken[i]++
			running++
			progress = true
		}
		if !progress {
			break
		}
	}
	return out
}
