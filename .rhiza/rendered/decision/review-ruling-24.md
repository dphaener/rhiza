# review-ruling-24

decision · `01M602SZESEBES1TT6BT4NX88B`

## session

0e5a9d92-bc68-4281-9e5b-01a17f6988d0

## statement

Review sessions get RHIZA_CHECK, RHIZA_CYCLE and RHIZA_ROLE (check, critic, fixer or verifier); check and critic prompts take {check} and {cycle}, and the fixer prompt takes {leaf} and {cycle}, with {leaf} empty for the change-set fixer.

## wording

the agent's choice in its plan, approved by the orchestrator acting as arbiter at plan approval

## Links

### In

- decided-by [review-checks-step](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-checks-step.md)
- decided-by [review-workflow-file](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-workflow-file.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
