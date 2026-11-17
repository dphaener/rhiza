# scripts-env-isolated-guard

capability · `01M7M3874NE4M6GA6XEVKK3RW2` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScriptEnvIsolated passes: the test sets, with t.Setenv, RHIZA_CHANGE, RHIZA_ACTOR, RHIZA_SESSION, RHIZA_RUN, SHIP_HEAD_BUILD, SHIP_HEAD_RHIZA, REX_SESSION and REX_BLOCK to marker values, runs testscript through scriptParams over a script written to t.TempDir() (not under testdata/script, so TestScripts never picks it up), and checks both that an exec'd command's environment (exec sh -c env) and the names and values scriptEnv returns hold none of those variables, while RHIZA_TEST_NOW and RHIZA_TEST_SEED are still there; go test -trimpath ./cmd/rhiza -run 'TestScripts/^session-run$' still passes. (failure case) With keepScriptEnv temporarily changed to keep every name in os.Environ(), TestScriptEnvIsolated fails and its message names a leaked variable (e.g. RHIZA_CHANGE); the change is reverted before commit

## files

cmd/rhiza/script_env_test.go (new)

## must_not

Change main_test.go's setup, keepScriptEnv or scriptEnv, internal/daemon/reexec.go or any production code; add a txtar under cmd/rhiza/testdata/script or an entry on processScripts; call t.Parallel (t.Setenv forbids it)

## statement

A cmd/rhiza test proves that RHIZA_*, SHIP_* and REX_* variables set in the test process's own environment reach neither a testscript's environment nor the environment the in-process rhiza is handed (scriptEnv), so a script behaves the same inside a run session as in CI

## Links

### Out

- decided-by [decision-scope-env-regression-guard](../../../../decision/decision-scope-env-regression-guard.md)

### In

- verifies [scripts-env-isolated-guard-verified](../../../../evidence/scripts-env-isolated-guard-verified.md)
- verifies [scripts-env-isolated-guard-qa-1](../../../../evidence/scripts-env-isolated-guard-qa-1.md)
- verifies [scripts-env-isolated-guard-qa-2](../../../../evidence/scripts-env-isolated-guard-qa-2.md)

## History

- `01M7M2YWF4KEF5RV32ZS0C4XK7` Fix cmd/rhiza testscripts inherit the calling session's RHIZA_*, SHIP_* and REX_* variables, so a script run from inside a run session fails or passes differently than in CI
