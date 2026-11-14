# handback-route-verified

evidence · `01M7BADRPZ6SHXD8KQBMN9E7FV` · actor [claude-code-a93310fa](../actor/claude-code-a93310fa.md)

## result

PASS

## statement

At 401f4b156 (attempt 2, leaf commits 62744d121 and 401f4b156): go test -trimpath ./internal/workflow -run TestHandbackTarget, ./internal/run -run TestHandbackRoute (3 subtests: two reds then pass re-enter build each time; a fourth red parks the run stopped with 'handback cap reached' naming check, build and the cause; handback_cap=1 sets the limit) and ./internal/api -run TestInboxHandbackParked all exit 0. In a clean detached worktree of HEAD with env -i HOME PATH TMPDIR GOTMPDIR: make self-check exit 0 (attempt 1's finding f1 is fixed), rhiza render 0 written/0 removed, rhiza ship generate 0 changed. Diff read: handback:<step> parsed in Next/HandbackTarget, handback_cap defaults to 3 and is refused below 1 or non-integer in Parse and Validate, checkTargets refuses unknown and non-agent targets naming the step, the journal's handback entry counts hand-backs and resets visit caps, the inbox handback kind names the step and cause; supporting edits to include.go, spec.go, docs, web_inbox.go and daemon.go are within scope; no TODOs or stubs; nothing is delivered to the session.

## Links

### Out

- verifies [handback-route](../capability/rhiza/workflows/workflow-runs/handback-route.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
