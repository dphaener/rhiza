# ship-head-build-on-demand-verified

evidence · `01M7C24J3ASSSD81BV9KQ0NXGX` · actor [claude-code-9f49e4c4](../actor/claude-code-9f49e4c4.md)

## result

PASS

## statement

At d97b0e7f1: go test -trimpath ./internal/api -run TestShip exit 0; go test -trimpath ./internal/api -run TestShipHeadBuildOnDemand -v passes all six cases (missing binary built once then rendered with and stamped with HEAD; built-at an earlier commit rebuilt; built at current HEAD not rebuilt; empty SHIP_HEAD_BUILD fails 'not a built rhiza' with no build; failing build carries its output and 'nothing was rendered or pushed'; ship.land with a failing build makes no commit or push); go test -trimpath ./cmd/rhiza -run TestScripts/ship-land-branch-build exit 0 with the new case where a missing SHIP_HEAD_RHIZA is built by ship land, logs land and render to head.log and lands. Diff read: shipHeadBuild (used by ship land and every ship render) builds via shipRebuildHead (sh -c in ctx.Root) when the binary is missing or its .built-at record differs from HEAD; shipRebuildHead stamps .built-at and adds the build output to its error. runShipCloseOut never renders or uses the head binary, so the land path, which shares shipHeadBuild, covers the acceptance's close-out case. Only the three named files plus the actor record changed; no stubs or TODOs.

## Links

### Out

- verifies [ship-head-build-on-demand](../capability/rhiza/orchestration/shipping/ship-head-build-on-demand.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
