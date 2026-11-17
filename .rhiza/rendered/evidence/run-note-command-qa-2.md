# run-note-command-qa-2

evidence · `01M7M1V7MB843Z9Y9G3E3PXYRX` · actor [claude-code-91a0ded1](../actor/claude-code-91a0ded1.md)

## result

PASS

## statement

QA cycle 2 at c404b4730, scratch-built rhiza: run note journals with role/leaf/session from env (role-only and leaf-only too), via --run, --json and api call run.note -; run status human '## Notes' and --json notes oldest first; empty, spaces-only, \n, \r, 201 chars, missing and extra args exit 2 with nothing journalled (200 x and 200 é accepted, tab accepted); no RHIZA_RUN names RHIZA_RUN; unknown run exit 4; --dry-run writes nothing. Scoped tests internal/run TestNote, internal/api TestRunNote, TestScripts/run-note and describe-golden pass.

## Links

### Out

- verifies [run-note-command](../capability/rhiza/workflows/workflow-runs/run-note-command.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
