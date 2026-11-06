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
// together. after is the project that took the last turn of the previous
// pick ("" for none): turns start at the project after it in the order
// given, wrapping round, so the rotation carries on from pick to pick
// instead of restarting at the first project. Each turn goes to the next
// project in that order that has a ready change set and is under its own
// Max, and takes that project's next ready change set; a project with
// nothing to take is passed over. It stops when the slots are full or no
// project can take a turn. LastTurn of the result is the next pick's after.
func RoundRobin(projects []ProjectPick, running, maxRuns int, after string) []Picked {
	n := len(projects)
	start := 0
	for i, p := range projects {
		if p.Project == after {
			start = i + 1
			break
		}
	}
	next := make([]int, n)
	taken := make([]int, n)
	var out []Picked
	for running < maxRuns {
		progress := false
		for j := 0; j < n; j++ {
			if running >= maxRuns {
				break
			}
			i := (start + j) % n
			p := projects[i]
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

// LastTurn is the project that took the last turn in picked, or after when
// nothing was picked: the after of the next RoundRobin.
func LastTurn(picked []Picked, after string) string {
	if len(picked) == 0 {
		return after
	}
	return picked[len(picked)-1].Project
}
