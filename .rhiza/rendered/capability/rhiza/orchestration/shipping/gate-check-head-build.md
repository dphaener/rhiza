# gate-check-head-build

capability · `01M7EHQ2K4VECK48QNRCE286P5` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./internal/api -run TestCheckHeadBuild shows, with SHIP_HEAD_RHIZA naming a fake head that records its argv: rhiza check invoking it as check --json with --base, --skip-reference and --branch-newer-format passed through when set, and returning the fake's report, an ok one and one with a fail finding that the fake prints while exiting 1; with SHIP_HEAD_RHIZA unset, or naming the running binary, check running in-process with no exec; with the head missing and SHIP_HEAD_BUILD set, the build running first; with it missing and SHIP_HEAD_BUILD empty, or the fake printing no report, check failing with an error naming the head build; with SHIP_HEAD_RHIZA naming the fake and the running binary not the file PATH's rhiza resolves to (a self-check build in a temp folder), check running in-process with no exec and no build, even with SHIP_HEAD_BUILD set; and runCheck called directly (as the pre-push check and run code-check call it) running in-process with SHIP_HEAD_RHIZA set. go test -trimpath ./cmd/rhiza -run 'TestScripts/ship-gate-head-check$' passes a new testscript: a repo whose docs/reference was written by rhiza docs and then edited to document a CLI flag this build does not have (a branch that adds the flag), where plain rhiza check fails on the generated-reference comparison; a workflow run whose gate step runs rhiza check, with SHIP_HEAD_RHIZA naming a fake head build that logs its argv and prints an ok check report, passes the gate and the log shows check --json; the same run with SHIP_HEAD_RHIZA unset fails the gate on command 1 rhiza check. go test -trimpath ./internal/api -run 'TestLandRenderHeadBuild|TestShipHeadBuildOnDemand|TestShipPrePushCheck' still passes, make self-check passes, and SHIP_HEAD_RHIZA=$TMPDIR/no-head-rhiza SHIP_HEAD_BUILD= make self-check passes too (its fresh build checks in-process instead of failing on the missing head). Touches internal/api/check.go, internal/api/ship_build.go, internal/api/ship_build_test.go, cmd/rhiza/testdata/script/ship-gate-head-check.txtar and .rhiza/workflows/ship-rhiza.toml.

## statement

In a run whose environment names SHIP_HEAD_RHIZA, rhiza check (the CLI and api call check) runs the PR head's build instead of checking in-process: it gets the binary from shipHeadBuild (built on demand under SHIP_HEAD_BUILD) and execs its check --json in the same directory with the request's --base, --skip-reference and --branch-newer-format, only when the running binary is the rhiza that PATH resolves (the same file as exec.LookPath("rhiza"), as the gate's plain rhiza check is) and is not that build, and reports the child's check report (read from its stdout even when it exits non-zero on failing findings) with check's usual output and exit; so the conductor default's gate step (plain rhiza check) judges docs/reference and the docs-command rule by the branch's own registry, and a change set that adds a CLI command, flag or setting no longer fails the gate and hands back to conduct in a loop. With SHIP_HEAD_RHIZA unset, naming the running binary, or a running binary that is not the PATH rhiza, check runs in-process as today: so make self-check's fresh build in its temp folder ($dir/rhiza check, run by builders, fixers and tools/gate under the workflow [env]) still checks the repo with this commit's build, with no exec and no head build; a missing head with SHIP_HEAD_BUILD empty, a failed build, or output that is no check report fails with an error naming the head build. The pre-push check, run code-check and health keep calling runCheck in-process. .rhiza/workflows/ship-rhiza.toml's comment says the head build serves check as well as land and render.

## Links

### Out

- decided-by [gate-check-via-head](../../../../decision/gate-check-via-head.md)
- decided-by [gate-check-head-fallback](../../../../decision/gate-check-head-fallback.md)

### In

- verifies [gate-check-head-build-verified](../../../../evidence/gate-check-head-build-verified.md)
- verifies [gate-check-head-build-qa-1](../../../../evidence/gate-check-head-build-qa-1.md)

## History

- `01M7EHHAG3CTRQCJCWSG2W7BTZ` Fix the gate step: the conductor default's gate runs rhiza check with the installed rhiza, not the PR head's build (SHIP_HEAD_RHIZA), so a change set that adds a CLI command, flag or setting fails the gate on docs/reference files its own build passes and hands back to conduct in a loop (runs 01M7DVV7GNJAWG29QZC5F2JWC3 and 01M7E56MGTY4M60XZNEG2A6537; request 01M7EAYSYK7Z32QRP0FB4XPKJN)
