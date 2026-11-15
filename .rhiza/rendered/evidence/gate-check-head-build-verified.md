# gate-check-head-build-verified

evidence · `01M7ENMBSPQ54J23X07PZ1NBCV` · actor [claude-code-da592419](../actor/claude-code-da592419.md)

## result

PASS

## statement

Verified commit 21f24b894 in a detached worktree. go test -p 2 -trimpath ./internal/api -run 'TestCheckHeadBuild|TestLandRenderHeadBuild|TestShipHeadBuildOnDemand|TestShipPrePushCheck': exit 0. TestCheckHeadBuild covers each acceptance point: a named head runs as check --json with --base, --skip-reference and --branch-newer-format passed through; an ok report, and a fail report with exit status 1, are reported with check's usual exit; unset or self-named head runs in-process with no exec; a missing head with SHIP_HEAD_BUILD set is built first; a missing head with no build command, or output that is no report, errors and names the head build; a running binary that is not the PATH rhiza runs in-process with no build; runCheck called directly runs in-process. go test -p 2 -trimpath ./internal/api (full package, 750 tests): exit 0 in 416s under machine load; two earlier runs hit the 240s cap with no failures, TestReadBudget alone taking 105s. go test -p 2 -trimpath ./cmd/rhiza -run 'TestScripts/(ship-land-branch-build|run-ship-land-conflict|docs-reference-stale|ship-local-gate|ship-gate-head-check)$': exit 0; the new ship-gate-head-check script shows plain rhiza check failing on an edited docs/reference/cli.md, the head-gate run passing with head.log holding check --json, and the plain-gate run failing on command 1 rhiza check. make self-check: exit 0. SHIP_HEAD_RHIZA=$TMPDIR/no-head-rhiza SHIP_HEAD_BUILD= make self-check: exit 0, no head built. rhiza ship generate and rhiza render: exit 0, no changes. Diff read against the Decisions gate-check-via-head and gate-check-head-fallback: the check op's entry (runCheckOp) redirects only outside dry-run; the in-process runCheck callers are unchanged; the child runs with headChildEnv. ship-rhiza.toml's comment now names check. Beyond the listed files: ship_prepush.go passes headChildEnv to the pre-push head check (a recursion guard), and ship_build.go adds nil ctx.Git guards; both are small and in scope. No stubs or TODOs.

## Links

### Out

- verifies [gate-check-head-build](../capability/rhiza/orchestration/shipping/gate-check-head-build.md)

## History

- `01M7EHHAG3CTRQCJCWSG2W7BTZ` Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN)
