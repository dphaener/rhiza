# scripts-in-process-e-q-verified-tip

evidence · `01M7HW20PBTY7RH7ED5TNXMPVY` · actor [claude-code-8afa6b71](../actor/claude-code-8afa6b71.md)

## result

PASS

## statement

Detached checkout at 50f3ea4e3 (leaf tip; leaf commits 878916791, b8da7ea13, 50f3ea4e3). Capped, RHIZA/SHIP/REX unset: go test -p 2 -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^[e-q]' ok 27.4s; -run 'TestProcessScripts|TestInProcess' ok (3 tests PASS); go vet ok; gofmt -l empty. grep -cE '^\s"[e-q][^"]*":\s+notConverted' cmd/rhiza/main_test.go = 0; 398 txtar files. git diff --word-diff=porcelain of the txtar files shows only removed exec: 809 (878916791), 0 (b8da7ea13), 20 (50f3ea4e3); only e-q txtar files and e-q processScripts entries changed. Failure case: model-commit-untracked (off the list) with line 11 back to exec rhiza fails with 'started rhiza as a process but is not on processScripts' (reverted). Reason probe: a scratch patch to rhizaProcess logged every rhiza process start with its parent's cmdline over the whole e-q range (1514 starts, stand-in rex/gh/claude first on PATH, none hit). Every kept e-q entry's reasons match its parents: all 15 'a workflow step runs rhiza' scripts start rhiza from a step shell (sh -c exec 2>&1 ...: 1 to 28 per script); fake bin/rex tab commands and bin/claude/bin/gate callbacks call rhiza where the entries say so; sh -c and the mkchanges.sh/mkcs.sh helpers start rhiza in the 'through sh -c' scripts; & lines exist in hook-commands, inbox-advisor, inbox-api-error, inbox-permission, inbox-stalled and feed-follow (exec sh -c 'exec rhiza feed --follow' &f1&); init has a stdin feed; quickstart starts rhiza only from the testscript binary (quickstart-doc); plan-workflow-ci-infra-wait's bin/rhiza execs the real one for steps. Conversion probe: plan-workflow-ci-infra-timeout (only reason: workflow step) taken off the list with every exec rhiza made bare fails at :29, its wait-pr step started rhiza as a process. Note: hook-commands keeps all 15 exec rhiza lines though only its daemon & line needs a process; the leaf allows kept scripts to keep exec lines.

## Links

### Out

- verifies [scripts-in-process-e-q](../capability/rhiza/integrity/self-check/scripts-in-process-e-q.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
