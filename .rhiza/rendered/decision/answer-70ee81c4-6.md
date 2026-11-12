# answer-70ee81c4-6

decision · `01M7644X7WPDGBY1SYYAMCBEZP`

## options

```json
[
  "Built-ins + project additions (Recommended)",
  "All project-configured",
  "Leave as is"
]
```

## question

The watchdog's transient-failure signatures include Go's go-build $WORK error, and ci-triage parses Go, RSpec, Minitest and pytest. What about these?

## quote

We may want to chat about this, but isn't just reading the exit code good enough? If we need to detect what kind of error, maybe a quick  call out to an LLM (headless) using a model like Haiku would be a better choice. Let's talk about it.

## quote_where

session 70ee81c4-37df-450f-86ae-01a1cbc527f8, answering "The watchdog's transient-failure signatures include Go's go-build $WORK error, and ci-triage parses Go, RSpec, Minitest and pytest. What about these?"

## session

70ee81c4-37df-450f-86ae-01a1cbc527f8

## statement

The watchdog's transient-failure signatures include Go's go-build $WORK error, and ci-triage parses Go, RSpec, Minitest and pytest. What about these? — We may want to chat about this, but isn't just reading the exit code good enough? If we need to detect what kind of error, maybe a quick  call out to an LLM (headless) using a model like Haiku would be a better choice. Let's talk about it.

## Links

### Out

- derived-from [workflow-cleanup-plan-2026-11-11](../source/workflow-cleanup-plan-2026-11-11.md)

## History

- `01M7648Q6T5DGZG4EKRS53XWND` Record the answers of session 70ee81c4-37df-450f-86ae-01a1cbc527f8, whose change set 01M75W99SF6Q1T4QJM5GKA4PM1 has landed
