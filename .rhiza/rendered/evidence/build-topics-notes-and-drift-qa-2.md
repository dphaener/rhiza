# build-topics-notes-and-drift-qa-2

evidence · `01M7M1V81VEF98FGFMMFG0G8FK` · actor [claude-code-91a0ded1](../actor/claude-code-91a0ded1.md)

## result

PASS

## statement

QA cycle 2 at c404b4730, scratch-built rhiza guide: builder and verifier carry `rhiza run note "<learning>"` before session done; builder carries update-the-landed-node-with-rhiza-set-or-raise-a-question; build-retry says findings are quoted verbatim in the prompt; conduct says the retry brief quotes findings verbatim and inspects the working tree; verifier keeps 'Never edit code'; check-*, critic, plan-critic and fixer topics tell no session to journal a note (checks name run notes only in their unchanged never-read rule). internal/guide tests and internal/api TestGuideInstall pass.

## Links

### Out

- verifies [build-topics-notes-and-drift](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-topics-notes-and-drift.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
