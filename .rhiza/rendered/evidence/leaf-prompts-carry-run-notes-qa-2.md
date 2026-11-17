# leaf-prompts-carry-run-notes-qa-2

evidence · `01M7M1V7XC2N4T129B4GKSVWXK` · actor [claude-code-91a0ded1](../actor/claude-code-91a0ded1.md)

## result

PASS

## statement

QA cycle 2 at c404b4730, same scratch leaves run: builder and verifier prompts after notes were journalled end with 'Notes from earlier sessions of this run:' and '- <note>' lines oldest first (retry prompts after the retry section); the first builder, with no notes, has no heading; rhiza run spawn children: builder and verifier brief.md carry both notes oldest first after the caller's brief (one with no trailing newline), check-plan, critic and fixer briefs carry none. TestLeafPrompt|TestRunSpawn, TestRenderSpawnBrief and TestScripts/run-spawn pass.

## Links

### Out

- verifies [leaf-prompts-carry-run-notes](../capability/rhiza/workflows/built-in-workflows/build-workflow/leaf-prompts-carry-run-notes.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
