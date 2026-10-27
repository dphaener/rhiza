# build-leaf-fails

capability · `01M5YX3G453SY3XT0RGZ5Z0C4A` · parent [build-workflow](../build-workflow.md)

## acceptance

An engine test over leaves A, B and C, C depending on A and A failing, finds a question blocking A, C skipped, B built and verified, and the run failed; after an answering Decision a resume builds A afresh with the Decision in its prompt and then C; in an autonomous run the arbiter is started on A's question and a second failure of A ends the run failed.

## statement

A leaf that fails becomes a question linked blocks to it and its dependents are skipped, while leaves that don't depend on it carry on, and the run ends failed once nothing more can proceed; a resume after the question is answered gives the leaf a fresh 3 attempts with the answering Decision in its builder's prompt, and in an autonomous run the arbiter answers at once and the leaf gets that one fresh round.

## Links

### Out

- depends-on [build-retries](build-retries.md)
- decided-by [build-workflow-shape](../../../../../decision/build-workflow-shape.md)
- decided-by [arbiter](../../../../../decision/arbiter.md)
- decided-by [answer-ab54d544-10](../../../../../decision/answer-ab54d544-10.md)
- decided-by [answer-ab54d544-11](../../../../../decision/answer-ab54d544-11.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-leaf-fails-verified](../../../../../evidence/build-leaf-fails-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
