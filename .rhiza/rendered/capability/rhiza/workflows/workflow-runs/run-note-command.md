# run-note-command

capability · `01M7KWHEF0XF0PZQ7KB7CEMBXJ` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/run -run TestNote and go test -trimpath ./internal/api -run TestRunNote pass with tests that journal two notes and find them in run status's notes oldest first with their role and leaf; (failure case) an empty note, a note holding a newline and a 201-character note are each refused with exit 2 and nothing journaled, and a note outside a run with no --run names RHIZA_RUN; go test -trimpath ./cmd/rhiza -run TestScripts/run-note passes; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after -update.

## files

internal/run/journal.go, internal/run/journal_test.go, internal/run/status.go, internal/api/run_note.go (new), internal/api/run_note_test.go (new), internal/api/run.go, cmd/rhiza/testdata/script/run-note.txtar (new), cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Write to the model or .rhiza/log; collide with or change rhiza note add; change any existing run status field

## statement

`rhiza run note <text>` (RHIZA_RUN's run, or --run) journals a one-line learning on the run with its session, role and leaf from RHIZA_SESSION, RHIZA_ROLE and RHIZA_LEAF when set; rhiza run status (human and --json) lists the run's notes oldest first; empty, multi-line or over-200-character text is refused with a usage error.

## Links

### Out

- decided-by [run-notes-journaled](../../../../decision/run-notes-journaled.md)

### In

- depends-on [leaf-prompts-carry-run-notes](../built-in-workflows/build-workflow/leaf-prompts-carry-run-notes.md)
- depends-on [build-topics-notes-and-drift](../built-in-workflows/build-workflow/build-guide/build-topics-notes-and-drift.md)
- verifies [run-note-command-verified](../../../../evidence/run-note-command-verified.md)
- verifies [run-note-command-qa-1](../../../../evidence/run-note-command-qa-1.md)
- verifies [run-note-command-qa-2](../../../../evidence/run-note-command-qa-2.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
