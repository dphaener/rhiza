# conductor-land-render-head-build-verified

evidence · `01M7C3WRB6C4BMQNPDREQ1FDNA` · actor [claude-code-9e71ab28](../actor/claude-code-9e71ab28.md)

## result

PASS

## statement

At eebc7598d re-ran the commit's Checks: go test -trimpath ./internal/api -run TestLandRenderHeadBuild (exit 0, all 5 subtests pass: named head binary execs 'render --json' and 'land --json' in ctx.Root with nothing rendered or landed in-process; SHIP_HEAD_RHIZA unset and running binary == head both stay in-process; missing binary with SHIP_HEAD_BUILD runs the build first; missing with empty SHIP_HEAD_BUILD fails render and land with 'not a built rhiza' and no exec or write), go test -trimpath ./internal/api -run 'TestShip|Lifecycle|TestRender|Land' (exit 0), rhiza ship generate (exit 0), rhiza render (exit 0). Read the diff: runViaHead in internal/api/ship_build.go resolves the head through shipHeadBuild and runs it via shipRunHead before any in-process work in runRender (render.go) and lifecycleVerb.run for land (lifecycle.go); headIsSelf compares by path and os.SameFile, so the child (Self == SHIP_HEAD_RHIZA) runs in-process with no recursion. Code changes touch only the four named files plus the builder's actor log/rendered record; no TODOs or stubs; tests use a fake Exec only.

## Links

### Out

- verifies [conductor-land-render-head-build](../capability/rhiza/orchestration/shipping/conductor-land-render-head-build.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
