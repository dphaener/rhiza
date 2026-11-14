# landed-run-resumable-past-merge-qa-2

evidence · `01M7DTKQKQTRCBWG4ZEH6FS5E3` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: go test ./internal/run -run TestPastLanding, ./internal/daemon -run TestResumable and ./internal/api -run TestRunResume pass, and the inbox-landed-resumable script passes with fixtures that visited merge and ship-change/merge.

## Links

### Out

- verifies [landed-run-resumable-past-merge](../capability/rhiza/orchestration/shipping/landed-run-resumable-past-merge.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
