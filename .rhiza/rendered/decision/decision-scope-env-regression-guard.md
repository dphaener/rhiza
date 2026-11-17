# decision-scope-env-regression-guard

decision · `01M7M35AENJXKE5HXZ3DCVNJNP`

## reasoning

The question reports the bug does not reproduce on main: testscript starts from a fixed environment and since #605 the in-process rhiza gets only scriptEnv's kept names (cmd/rhiza/main_test.go scriptEnv/keepScriptEnv). The intent (scripts behave the same in a run session as in CI) already holds, so the remaining value is locking it in with a test that runs with the hostile variables set. The one host-env path, daemon.NewReexec (internal/daemon/reexec.go), only uses os.Environ() when Reexec.Do execs after the executable's size or mtime changes; a go test binary is not rebuilt in place during a script, and the self-update tests in internal/api already swap daemonReexec for a fake. It cannot make a script pass or fail differently, so changing the daemon's re-exec env is out of this bug's scope and would touch real daemon behaviour (it must keep the caller's env when it re-execs in production).

## statement

Scope this change set to a regression guard, not a fix: one leaf adds a test in cmd/rhiza that sets RHIZA_*, SHIP_* and REX_* variables in the test process's own environment and asserts that neither a testscript's environment nor the environment handed to the in-process rhiza (scriptEnv) contains them, so #605's isolation cannot regress silently. daemon.NewReexec is left unchanged here.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-d921f30c-1](../question/question-d921f30c-1.md)
- derived-from [plan-fix-testscript-env](../source/plan-fix-testscript-env.md)

### In

- decided-by [scripts-env-isolated-guard](../capability/rhiza/integrity/self-check/scripts-env-isolated-guard.md)

## History

- `01M7M2YWF4KEF5RV32ZS0C4XK7` Fix cmd/rhiza testscripts inherit the calling session's RHIZA_*, SHIP_* and REX_* variables, so a script run from inside a run session fails or passes differently than in CI
