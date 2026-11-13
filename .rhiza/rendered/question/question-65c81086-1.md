# question-65c81086-1

question · `01M7947261QQBNQHAXG9GKNDC0`

## session

65c81086-02c4-4705-8890-01a1d241b16b

## statement

Flaky TestScripts/headless-run-agent-tabs (PR #513, attempt 1 of run 43333843332): leaf b's verifier exited after `rhiza add evidence` + `rhiza link ... verifies` (the fake claude exits 1 if either fails), yet the run judged it 'verifier ended with neither Evidence nor findings'; builder retries then failed (b.txt already committed). Leaf a's verifier ran concurrently and passed. So the run can miss a verifier's Evidence written while another session writes — likely a stale answer from the run's cached model/index (internal/api/run_model_cache.go readModel/catchUpModel) or a read racing the write. How should this change set proceed?

## Links

### In

- answers [verifier-evidence-read-fresh](../decision/verifier-evidence-read-fresh.md)

## History

- `01M7942GRFNQ82YV0V8F4A9VNA` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172; reproduce it locally, repeatedly, then fix it or quarantine it
