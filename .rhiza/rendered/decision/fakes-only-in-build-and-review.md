# fakes-only-in-build-and-review

decision · `01M6505FZ5YGMF00V28718TPB3`

## statement

The build guide and the review guide's non-qa checks tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh; the review guide's qa check runs the built binary hands-on, since fakes alone do not count for it

## wording

Build and review guides say to test runs with fakes only, never real Rex, agents or gh, except the review qa check, which runs the built binary hands-on: the change set's intent, the agent's wording (assumed, would have escalated)

## Links

### Out

- derived-from [agents-never-start-runs-plan](../source/agents-never-start-runs-plan.md)

### In

- decided-by [build-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-guide-fakes-only.md)
- decided-by [review-guide-fakes-only](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/review-guide-fakes-only.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
