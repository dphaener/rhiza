# gate-check-head-build-qa-1

evidence · `01M7EPCH3HTW8QJ2N8E4NP571J` · actor [claude-code-4965e067](../actor/claude-code-4965e067.md)

## result

PASS

## statement

Built ./cmd/rhiza to a scratch bin put first on PATH, in a fresh git repo with a fake head (logs argv and cwd, prints a given report). rhiza check with SHIP_HEAD_RHIZA=fake execed 'check --json' in the repo root (also from a subdir) and printed its ok report, exit 0; --base/--skip-reference/--branch-newer-format passed through; rhiza check --json and api call check gave identical bytes; a fail report with child exit 1 printed the fail finding, exit 1; garbage output, a missing head with SHIP_HEAD_BUILD empty, and a failing build each errored naming SHIP_HEAD_RHIZA/the head build, exit 1; a missing head with SHIP_HEAD_BUILD set was built then run. In-process with no exec: env unset, SHIP_HEAD_RHIZA naming the PATH binary, and the same build run by a path off PATH (even with a failing SHIP_HEAD_BUILD); a symlinked PATH rhiza still used the head. go test ./internal/api -run 'TestCheckHeadBuild|TestLandRenderHeadBuild|TestShipHeadBuildOnDemand|TestShipPrePushCheck' and TestScripts/ship-gate-head-check pass; make self-check and SHIP_HEAD_RHIZA=$TMPDIR/no-head-rhiza SHIP_HEAD_BUILD= make self-check both pass from a clean detached checkout of HEAD. Only finding: the missing/failed-head error tail mentions land/render/push (minor).

## Links

### Out

- verifies [gate-check-head-build](../capability/rhiza/orchestration/shipping/gate-check-head-build.md)

## History

- `01M7EHHAG3CTRQCJCWSG2W7BTZ` Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN)
