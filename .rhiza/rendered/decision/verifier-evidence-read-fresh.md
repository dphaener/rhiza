# verifier-evidence-read-fresh

decision · `01M7948KX29FV2FF2J93GANPJE`

## reasoning

The question's evidence shows leaf b's verifier did write Evidence and its verifies link (the fake claude exits 1 otherwise), yet the run judged 'neither Evidence nor findings' while leaf a's verifier wrote concurrently. In run_model_cache.go, readModel serves kept answers until catchUpModel refreshes, and refreshDue only refreshes after pollDebounce has passed or inside a writeSettle window that the runner opens for its own writes (wrote()), not for a session's writes; so a read right after a session ends can return an answer kept from before its writes, and another session's concurrent write can move the stamp in a way revalidate keeps a stale per-change-set answer. That makes the failure a product bug that could fail a real run's leaf (and burn builder retries), so quarantining alone would hide it. The change set's intent already says reproduce, then fix or quarantine; AGENTS.md says to reproduce timing flakes by injecting a delay, never by stress runs, and to quarantine with a bug filed when a fix is not possible. Fixing at the read (force a catch-up before judging an ended session) stays within the change set's intent.

## statement

Treat this as a real bug in the run driver, not a test flake: reproduce it locally by injecting a delay at the suspected race (e.g. hold leaf a's verifier write or the driver's post-session read in internal/api/run_model_cache.go) and running headless-run-agent-tabs a few times, then fix it so that the verdict on a session that has just ended is always read from a model caught up past that session's writes (the debounced catch-up in refreshDue and the answers kept across revalidate must not serve a pre-session answer). Add a regression test at the internal/api or internal/run level that fails without the fix. Quarantine the testscript (skip with a bug request filed via rhiza request --kind bug) only if the race cannot be reproduced or fixed within this change set.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-65c81086-1](../question/question-65c81086-1.md)
- derived-from [headless-tabs-flake-plan](../source/headless-tabs-flake-plan.md)

### In

- decided-by [verifier-verdict-reads-fresh-model](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step/verifier-verdict-reads-fresh-model.md)

## History

- `01M7942GRFNQ82YV0V8F4A9VNA` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172; reproduce it locally, repeatedly, then fix it or quarantine it
