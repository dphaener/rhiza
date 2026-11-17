# question-d921f30c-1

question · `01M7M34556E8DAX0FWZDJQYE5M`

## session

d921f30c-ba7e-42c6-89ad-01a1e82fbf90

## statement

Bug 'cmd/rhiza testscripts inherit RHIZA_*, SHIP_* and REX_*' does not reproduce on main: a probe script run from this run session (RHIZA_CHANGE, RHIZA_ACTOR, REX_SESSION... set) saw only RHIZA_TEST_* and GIT_CONFIG_*, because testscript starts from a fixed env and the in-process rhiza gets only scriptEnv's kept names (#605). The one host-env path left in the script path is daemon.NewReexec (internal/daemon/reexec.go) taking os.Environ()/os.Args, used by an in-process `rhiza daemon run`. How should this change set be scoped?

## Links

### In

- answers [decision-scope-env-regression-guard](../decision/decision-scope-env-regression-guard.md)

## History

- `01M7M2YWF4KEF5RV32ZS0C4XK7` Fix cmd/rhiza testscripts inherit the calling session's RHIZA_*, SHIP_* and REX_* variables, so a script run from inside a run session fails or passes differently than in CI
