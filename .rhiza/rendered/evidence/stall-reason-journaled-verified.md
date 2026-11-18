# stall-reason-journaled-verified

evidence · `01M7QCSBM89H3KJ0CJHNZVENCY` · actor [claude-code-2e7832c0](../actor/claude-code-2e7832c0.md)

## result

PASS

## statement

Verified commit 0a5a81c50 at HEAD. Re-ran its Checks (capped, -p 2): go test ./internal/api -run TestWatchdogStalled (TestWatchdogStalled, TestWatchdogStalledSpawnedChild, new TestWatchdogStalledReason all PASS), -run 'TestWatchdog|TestDaemonNudge|TestInbox' ok, go vet ./internal/api ok, rhiza ship generate 0 changed, rhiza render no diff beyond verifier-local actor files (restored). Skipped the full ./internal/api run per brief. Diff: stallReason names each stall session with role and spawn parent (or 'the step's session', no parent) and last report time; used in the charged attempt action, the budget-spent hold Detail (and so watchdogHoldItems' inbox item) and a new warn dep that appends the existing run.TypeWarning entry (wired in daemonWatchdog) only on a stop+retry, so nudges journal nothing; bug evidence only gains the reason. No new entry type or run state, run not failed, budget and limits unchanged. Test covers the fixer warning/attempt/hold/item/budget-hold, nudge and moved-since journaling nothing, and the step's own session with no parent. Only daemon_watchdog.go, its test and the builder's actor record changed; no TODOs.

## Links

### Out

- verifies [stall-reason-journaled](../capability/rhiza/orchestration/daemon/daemon-watchdog/stall-reason-journaled.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch
