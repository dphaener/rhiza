# conductor-land-render-head-build

capability · `01M7APRTE58PCDAC5VEJWK8Z59` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./internal/api -run TestLandRenderHeadBuild shows, with SHIP_HEAD_RHIZA naming a fake head binary that records its argv, rhiza render and rhiza land invoking it with their arguments and writing nothing in-process; with SHIP_HEAD_RHIZA unset, or the running binary being that build, both running in-process; with the binary missing and SHIP_HEAD_BUILD set, the build running first; with it missing and SHIP_HEAD_BUILD empty, both failing with nothing landed or rendered. Touches internal/api/ship_build.go, internal/api/render.go, internal/api/lifecycle.go and internal/api/ship_build_test.go.

## statement

In a run whose environment names SHIP_HEAD_RHIZA (a workflow's [env] reaches its agent steps), rhiza land and rhiza render run through ship's head build instead of in-process: they get the binary from shipHeadBuild (built on demand under SHIP_HEAD_BUILD, by ship-head-build-on-demand) and exec it with the same arguments and working directory, unless the running binary is that build, so the conductor's land and render on the branch use the branch's own rhiza, not an older installed one. With SHIP_HEAD_RHIZA unset they run in-process as today; a missing binary with SHIP_HEAD_BUILD empty, or a failed build, fails with nothing landed or rendered.

## Links

### Out

- depends-on [ship-head-build-on-demand](ship-head-build-on-demand.md)
- decided-by [builtin-ship-change-open-points](../../../../decision/builtin-ship-change-open-points.md)

### In

- depends-on [rhiza-repo-conductor-default](../../workflows/workflow-overrides/rhiza-repo-conductor-default.md)
- verifies [conductor-land-render-head-build-verified](../../../../evidence/conductor-land-render-head-build-verified.md)
- verifies [conductor-land-render-head-build-qa-2](../../../../evidence/conductor-land-render-head-build-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
