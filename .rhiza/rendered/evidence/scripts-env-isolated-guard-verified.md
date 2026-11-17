# scripts-env-isolated-guard-verified

evidence · `01M7M3WC4CBSG7B8A22NCQB9A8` · actor [claude-code-2507c73f](../actor/claude-code-2507c73f.md)

## result

PASS

## statement

In a detached checkout at 53f29ac4f: go test -trimpath ./cmd/rhiza -run TestScriptEnvIsolated -count=1 exit 0; go test -trimpath ./cmd/rhiza -run 'TestScripts/^session-run$' -count=1 exit 0; go vet ./cmd/rhiza exit 0; rhiza ship generate and rhiza render exit 0 with a clean tree. Failure case: keepScriptEnv changed to also add every os.Environ() name made TestScriptEnvIsolated exit 1 with 'scriptEnv: leaked from the test process: ... RHIZA_CHANGE (in RHIZA_CHANGE=) ...', then reverted. Diff: code change is only cmd/rhiza/script_env_test.go (new) plus the builder's actor record under .rhiza; the test uses t.Setenv for all eight RHIZA_/SHIP_/REX_ names, a t.TempDir() fixture via writeFixture, checks exec sh -c env and scriptEnv for names and marker values and requires RHIZA_TEST_NOW/RHIZA_TEST_SEED; no main_test.go, reexec.go, production, txtar or processScripts change, no t.Parallel, no TODOs.

## Links

### Out

- verifies [scripts-env-isolated-guard](../capability/rhiza/integrity/self-check/scripts-env-isolated-guard.md)

## History

- `01M7M2YWF4KEF5RV32ZS0C4XK7` Fix cmd/rhiza testscripts inherit the calling session's RHIZA_*, SHIP_* and REX_* variables, so a script run from inside a run session fails or passes differently than in CI
