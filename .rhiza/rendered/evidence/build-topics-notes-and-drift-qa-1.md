# build-topics-notes-and-drift-qa-1

evidence · `01M7KZEEJ6ACQSGHXVKDHXA86F` · actor [claude-code-805f61bc](../actor/claude-code-805f61bc.md)

## result

PASS

## statement

QA at cf68954a8, scratch-built rhiza guide: builder and verifier topics carry `rhiza run note "<learning>"` before session done; builder carries the update-landed-node-with-rhiza-set-or-raise-a-question rule; build-retry says findings are quoted verbatim in the prompt; conduct says the retry brief quotes findings verbatim and inspects the working tree; verifier keeps 'Never edit code'; no check-*, critic, plan-critic or fixer topic names `rhiza run note`. rhiza show of builder-topic, verifier-topic, build-retry-topic, conduct-skill matches. go test ./internal/guide/... and internal/api -run TestGuideInstall pass.

## Links

### Out

- verifies [build-topics-notes-and-drift](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-topics-notes-and-drift.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
