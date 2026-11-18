# question-2debf9c0-1

question · `01M7RBGT79TXH5BTXQZM4NHY4T`

## session

2debf9c0-0390-4456-af6c-01a1f0aeb98a

## statement

Builder of leaf byo-workflow-daemon-e2e found that the BYO example myteam-flow (docs/examples/byo/workflow.toml) cannot push: `rhiza ship open-pr` refuses a run with no journaled review (internal/api/ship_openpr.go:86, reviewJournaled), and myteam-flow declares no review (no checks step, no conduct checks). Question 01M7RBEHJ5 blocks the leaf. Which fix should this change set make?

## Links

### In

- answers [open-pr-review-only-when-declared](../decision/open-pr-review-only-when-declared.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
