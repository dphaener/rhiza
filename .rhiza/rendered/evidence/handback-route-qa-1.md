# handback-route-qa-1

evidence · `01M7BEHMX3S02HG4SN1KS62GG7` · actor [claude-code-a9fe07bb](../actor/claude-code-a9fe07bb.md)

## result

PASS

## statement

QA cycle 1, scratch build, fake rex/claude shims: workflow agent build -> command check with on.fail = handback:build and no handback_cap; check red every time: the agent session was resumed 3 times (same session ID), the journal held 3 handback entries, and the 4th red parked the run stopped; rhiza inbox showed a handback item naming step check and the cause ('went fail after 3 hand-backs to build: exit 1 / tests red 4') with next rhiza run resume. Resume re-parked on red, then succeeded on green. handback_cap = 1 parked after one hand-back. A hand-back inside an included workflow (ihb/check -> ihb/build) works. workflow check refused handback_cap 0, -2, "3", 2.5, handback:nope, handback:check (a command step) and handback: with the step named.

## Links

### Out

- verifies [handback-route](../capability/rhiza/workflows/workflow-runs/handback-route.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
