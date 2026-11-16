# cli-no-process-globals-verified

evidence · `01M7HDWR50BDXQR8S75T1ZJCNJ` · actor [claude-code-024a032f](../actor/claude-code-024a032f.md)

## result

PASS

## statement

Detached checkout at 81a527fed, run under a capped systemd scope with RHIZA_/SHIP_/REX_ unset. Re-ran the Checks block: go test -p 2 -trimpath -count=1 ./cmd/rhiza -run '^Test(NewCtx|CtxFromEnv|TUI|RunWorkingDir)' exit 0 (TestRunWorkingDir, TestCtxFromEnv, TestNewCtxHome, TestTUIAppHeaders, TestTUISettingsScopes pass); -run 'TestScripts/tui-' exit 0 (all 7 tui- scripts pass); go vet ./cmd/rhiza exit 0; -skip 'TestScripts|Budget' exit 0; rhiza ship generate and rhiza render exit 0 with a clean tree. TestRunWorkingDir chdirs outside any root, runs init/status/show against a scratch repo through run()'s wd, checks errors land in the given stderr buffer, and checks a missing wd exits 1 with 'working directory'. Diff read: no os.Getwd, os.Args[0], os.Chdir or os.Executable left on run()'s path (only in main() and selfPath, which main calls); ctx.Stderr, Terminal and Rex use run()'s stderr; the TUI apps keep their folder in an atomic pointer and pass the switcher a Chdir that moves it; openModal takes exe and cwd. Env handling, flags and output are unchanged. The only os.Stderr left is startProgram's terminal handover, which the statement does not cover. A built binary run from a deleted cwd exits 1 with 'working directory: unknown'.

## Links

### Out

- verifies [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
