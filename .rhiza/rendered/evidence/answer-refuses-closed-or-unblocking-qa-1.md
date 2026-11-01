# answer-refuses-closed-or-unblocking-qa-1

evidence · `01M6AZGXVDCAFH3669NYDVA2PW` · actor [claude-code-6df3bf06](../actor/claude-code-6df3bf06.md)

## result

PASS

## statement

QA at f270de1a: built go build -o $TMPDIR/rhiza-x ./cmd/rhiza and used it in a clean scratch git repo (RHIZA_* unset). With accepted change set A touching capabilities store and child (parent's child), and open questions store-q blocks store, par-q blocks parent, apart-q blocks apart, while another change set was current: question answer apart-q --change A exited 2 (usage: not an open question blocking this change set) and show answer-apart-q exits 4; par-q (blocks an ancestor of a touched node) and store-q were accepted and written to A; answering store-q again, with --json and with --dry-run, exited 2 with code usage and answer-store-q-2 does not exist; answering apart-q into the current Other change set exited 2; --record is refused. go test ./cmd/rhiza -run TestScripts/question-answer -count=1 ok; go vet and gofmt clean.

## Links

### Out

- verifies [answer-refuses-closed-or-unblocking](../capability/rhiza/roadmap-app-accept-answer/question-answer-op/answer-refuses-closed-or-unblocking.md)

## History

- `01M69VQG7GJARDSWZWKDJ4CWQ1` question.answer refuses a question that is not open or does not block the target change set, so an answer always lands in the change set the question blocks
