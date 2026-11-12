# render-runs-generate-qa-1

evidence · `01M76M976BDH5W1NMKW97WW5Z6` · actor [claude-code-d94dc3e8](../actor/claude-code-d94dc3e8.md)

## result

PASS

## statement

This repo's rhiza.ship.generate reads 'make docs' (git config). make docs at the head left docs/reference unchanged ('3 unchanged'); make self-check on a clean worktree of HEAD printed ok. Built rhiza ship generate in a scratch repo: unset skipped with exit 0; 'echo new > gen.txt' passed and staged only gen.txt, leaving a file already dirty before unstaged; 'echo boom; exit 3' printed its output and exited 1. TestScripts/(render-generate|describe-golden) ok.

## Links

### Out

- verifies [render-runs-generate](../capability/rhiza/orchestration/shipping/wait-green-pr/render-runs-generate.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
