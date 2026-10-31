# shipped-skip-reason-verified

evidence · `01M69PQA2X9K6B3KHAND2VXW3E` · actor [claude-code-67542e90](../actor/claude-code-67542e90.md)

## result

PASS

## statement

Verified at 8d5d49ec. internal/daemon/pick_test.go TestSkips has shipped cases (succeeded ship run after the accept -> SkipShipped in both the one-slot and room tables); go test ./internal/daemon -run 'Skips|Pick' passes. daemon-status.txtar adds a Shipped change set with a hand-written succeeded ship-change journal (.git/rhiza/runs/SHIPPED0001.jsonl) under the fake bin/rex, and asserts 'Skipped: ... Shipped: its latest run shipped it' and '"reason": "shipped"'; go test ./cmd/rhiza -run TestScripts/(daemon-status|describe-golden)$ passes. DaemonSkip.Reason desc now lists shipped; re-running describe-golden with -update leaves the tree clean, so the golden is regenerated. make check fmt/vet/test pass; self-check passes with the run env (RHIZA_ROOT etc.) unset; with it set it only flags this run's live actor render files.

## Links

### Out

- verifies [shipped-skip-reason](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skip-reasons/shipped-skip-reason.md)

## History

- `01M69DFHXFG7197RGNYYPWS4A5` The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened
