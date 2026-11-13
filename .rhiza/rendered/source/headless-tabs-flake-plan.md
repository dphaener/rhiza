# headless-tabs-flake-plan

source · `01M7949F09B7VWF6S1H3D2KMBX`

## quotes

```json
[]
```

## raw_sha256

f6cf0e15b04c70037a60072930e0ffe6f161bc42c90a590497781441c72c4d0c

## session

65c81086-02c4-4705-8890-01a1d241b16b

## summary

Autonomous plan run, no owner. The planner read CI run 43333843332 attempt 1 (PR #513, head 7a96ac2): TestScripts/headless-run-agent-tabs failed because leaf b's verifier was judged 'verifier ended with neither Evidence nor findings' although the fake claude only exits after rhiza add evidence and rhiza link verifies succeed; leaf a's verifier wrote concurrently and passed; b's builder retries then failed with exit 1 (its leaf file was already committed). The planner traced the verdict to stepRunner.Evidence, which reads through readModel's cached answers (internal/api/run_model_cache.go: catchUpModel, refreshDue, revalidate). The planner asked whether to fix the product race, fix only the test's timing, or quarantine; the arbiter ruled (Decision verifier-evidence-read-fresh) to treat it as a run-driver bug: reproduce with an injected delay, make the verdict on a just-ended session read a model caught up past its writes, add a regression test below the testscript, and quarantine with a bug filed only if it cannot be reproduced or fixed here. Not breaking, not severe.

## title

Plan of the headless-run-agent-tabs flake (autonomous)

## Links

### In

- derived-from [verifier-evidence-read-fresh](../decision/verifier-evidence-read-fresh.md)
- derived-from [verifier-verdict-reads-fresh-model](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step/verifier-verdict-reads-fresh-model.md)

## History

- `01M7942GRFNQ82YV0V8F4A9VNA` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172; reproduce it locally, repeatedly, then fix it or quarantine it
