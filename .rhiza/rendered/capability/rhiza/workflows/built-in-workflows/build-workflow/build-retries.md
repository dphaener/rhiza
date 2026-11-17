# build-retries

capability · `01M5YX3F9R6BSD6XKKTVQXYSDZ` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test whose verifier fails twice then passes finds three builder sessions on the ladder's models, each prompt naming its attempt, the second's spec carrying the first attempt's journalled findings and the third's the second attempt's failure, and the leaf passed; one whose verifier fails three times finds no fourth builder and the leaf failed. A prompt test finds a retry builder's prompt quoting the previous attempt's finding reasons byte for byte (a multi-line one and one over 200 characters), the inspect-the-working-tree line and the address-every-finding-as-written line, or with no findings that attempt's failure detail and no findings list.

## statement

A failed attempt goes to a fresh builder session on the next model of the ladder, its prompt naming the leaf and the attempt, telling it to inspect the working tree and the earlier attempts' commits before writing, quoting verbatim and in full each finding the previous attempt's verifier journalled on the leaf (or that attempt's failure detail when it journalled none), and telling it to address every finding exactly as written in one fix commit with the leaf's trailer; after 3 failed attempts the leaf fails.

## Links

### Out

- depends-on [build-model-ladder](build-model-ladder.md)
- depends-on [build-verifier-findings](build-verifier-findings.md)
- depends-on [build-leaf-commit](build-leaf-commit.md)
- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [answer-ab54d544-3](../../../../../decision/answer-ab54d544-3.md)
- decided-by [answer-ab54d544-5](../../../../../decision/answer-ab54d544-5.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-leaf-fails](build-leaf-fails.md)
- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-retries-verified](../../../../../evidence/build-retries-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
