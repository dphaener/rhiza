# scripts-env-isolated-guard-qa-2

evidence · `01M7M4N37GS56FSSR0YGH165YP` · actor [claude-code-bcc50aa0](../actor/claude-code-bcc50aa0.md)

## result

PASS

## statement

QA rerun at 53f29ac4f (code tree clean; diff from merge-base touches only cmd/rhiza/script_env_test.go plus .rhiza records: no main_test.go, reexec.go, production, txtar or processScripts change, no t.Parallel). Read the test and its helpers scriptParams, keepScriptEnv, scriptEnv (main_test.go) and writeFixture. All runs under a capped systemd scope with go test -trimpath -p 2 -count=1, session RHIZA_*/SHIP_*/REX_* unset unless noted. Happy: -run '^TestScriptEnvIsolated$' PASS. Failure: keepScriptEnv temporarily also adding every os.Environ() name made it FAIL with 'scriptEnv: leaked from the test process: RHIZA_CHANGE (in RHIZA_CHANGE=), RHIZA_ACTOR, RHIZA_SESSION, RHIZA_RUN, SHIP_HEAD_BUILD, SHIP_HEAD_RHIZA, REX_SESSION, REX_BLOCK'; reverted with git checkout, tree clean. -run 'TestScripts/^session-run$' PASS. Edge cases: -count=2 PASS twice in one process; '^TestScriptEnvIsolated$|TestScripts/^session-run$' in one process PASS both; with the session's real RHIZA_*/REX_*/SHIP_* (16 vars incl. unlisted RHIZA_CHECK, REX_BIN_DIR) still set it PASS and the logged exec env holds none of them while RHIZA_TEST_NOW and RHIZA_TEST_SEED are present. Scratch build go build -trimpath -o scratchpad/rhiza-qa ./cmd/rhiza ok.

## Links

### Out

- verifies [scripts-env-isolated-guard](../capability/rhiza/integrity/self-check/scripts-env-isolated-guard.md)

## History

- `01M7M2YWF4KEF5RV32ZS0C4XK7` Fix cmd/rhiza testscripts inherit the calling session's RHIZA_*, SHIP_* and REX_* variables, so a script run from inside a run session fails or passes differently than in CI
