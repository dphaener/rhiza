# run-note-command-verified

evidence · `01M7KX2SBYQ94XS6GJ49VH7V8A` · actor [claude-code-dd51a582](../actor/claude-code-dd51a582.md)

## result

PASS

## statement

Verifier at d122c8ae9 re-ran every Checks: line, all exit 0: go test -trimpath ./internal/run -run TestNote (TestNoteFolds, TestNoteProblem), ./internal/api -run TestRunNote (InStatus, Refused, DryRun), ./cmd/rhiza TestScripts/run-note, describe-golden -update (golden unchanged) then describe-golden, rhiza ship generate (0 changed), rhiza render (only the verifier's own actor, reverted). Tests journal two notes found in run status notes oldest first with role and leaf; empty, newline and 201-char notes exit 2 with nothing journaled; no --run names RHIZA_RUN. A scratch-built binary confirmed exit 2 for a newline, an empty and a run-less note. Diff adds only the run.note op, the note journal entry and a Notes field/section; rhiza note add and existing status fields untouched; the script cmp's .rhiza/log unchanged; no TODOs.

## Links

### Out

- verifies [run-note-command](../capability/rhiza/workflows/workflow-runs/run-note-command.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
