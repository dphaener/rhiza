# cli-no-process-globals-fix-verified

evidence · `01M7HQ29SZV7TZBHZNYRSP0XY5` · actor [claude-code-be43c5a6](../actor/claude-code-be43c5a6.md)

## result

PASS

## statement

Verified fix commit b6447abda in a detached checkout, RHIZA_/SHIP_/REX_ vars unset, capped scopes. Checks block re-run: go test ./cmd/rhiza -run '^Test(RunWorkingDir|NoProcessGlobals)$' exit 0; ./cmd/rhiza -skip 'TestScripts|Budget' exit 0 (28.6s); TestScripts/^(apply|apply-concurrent|e2e|import-retires-docs|import-unbuilt-as-drafts|leaf-cycle|tui-.*)$ exit 0 (13 scripts); ./internal/api ./internal/index -skip Budget exit 0 (the block's exit 1 was the pre-test-fix run it explains; now green); ./internal/api -run '^(TestInit|TestWorkflowFile|TestRunStartWorkflowRefusals)$' exit 0; go vet ./cmd/rhiza ./internal/api ./internal/index exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (only an unrelated actor from shared session records appeared). Acceptance: -run '^Test(NewCtx|CtxFromEnv|TUI|RunWorkingDir)' exit 0 and TestScripts/tui- exit 0. TestRunWorkingDir t.Chdir's to an outside folder holding in/writes.json, in/s.json, in/raw.txt with other content, runs apply and source add --raw with relative paths and run() given the repo: q-given exists, q-process is exit 4, founding shows 'Given call' not 'Process call', and a missing relative file is exit 2 'apply: cannot read in/none.json'. TestNoProcessGlobals now keys by import path and flags path/filepath.Abs under any name: fixtures show internal/api/apply.go:6 (aliased fp.Abs) and cmd/rhiza/ctx.go:8 (non-allowlisted root()) flagged while selfPath's Abs is allowlisted; grep finds the only remaining filepath.Abs in guarded packages at cmd/rhiza/ctx.go:451 (selfPath, called by main() only), with its reason in the allowlist. Diff read: apply, source add file/--raw, run findings/verdict/sign-off/inbox files, workflow file refs, init/root init/session/pre-push/daemon-projects resolve via Ctx.abs/absIn; run spawn --brief and project add already joined ctx.Wd; RHIZA_ROOT and child Exec/Start/Detach dirs use childDir(wd). Error messages still name the file as given and keep their exit codes (Usagef/NotFoundf unchanged); no flag or output change. Observation, not blocking: runUserStatusLine (internal/api/usage_live.go:104) and shellWithStdin (cmd/rhiza/ctx.go:306) start sh -c with no Dir, so they inherit the process's folder; neither takes a path from the command line.

## Links

### Out

- verifies [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
