# lost-session-judged-verified

evidence · `01M6M9T5CA6ZMKHA23EXC232T8` · actor [claude-code-db66367e](../actor/claude-code-db66367e.md)

## result

PASS

## statement

At d1e34cea (a model-only commit; the code landed in 7fd9c685), ran go test -count=1 -run TestAgentLostSession -v ./internal/run: all 4 subtests PASS. Each one rejoins journalled session S-earlier (Entered+Session entries, not started by this process) whose AgentStatus is Lost with Outcome fail. In 'plain' (agentWF) the step is Failed with detail containing 'lost: SessionEnd'. In 'done and until holds' (agentDoneWF, until_when = done, SaidDone, until true) it is Succeeded with 'until held'. In 'done, until fails' it is Failed with 'until did not hold'. In 'until step without done mark' it is Failed without checking until. Every case asserts r.closed == [S-earlier], so the Runner was asked to close the pane. Read pollAgent in internal/run/engine.go: a Lost status passes only through exitedUntil when step.Until is set, WaitsForDone is true and SaidDone is true; otherwise its fail outcome stands. stepRunner.AgentStatus in internal/api/run.go maps a lost exit to Done, OutcomeFail, Lost and never to exit 0. go test ./internal/run ./internal/session pass and go vet is clean. ./internal/api passed on a rerun (387s); its first run failed one timing budget under CPU contention (539s), which is unrelated to this leaf.

## Links

### Out

- verifies [lost-session-judged](../capability/rhiza/workflows/workflow-runs/lost-session-judged.md)

## History

- `01M6M4YHKJAER13CM0YHZQFKCZ` Fix A run resumed after a reboot can wait forever on a session that already finished: text references' run 01M6K32AED9M0G432791PNBS8B was resumed by the daemon at 18:35 on 2026-11-04 after the VM reboot, logged 'continuing visit 1' and then nothing for 2+ hours, while its fixer session b0eafd65 shows Status: done at 18:27 (from hook) and its Rex block no longer exists. On resume, re-read each in-flight session's recorded status and judge sessions already done or exited (and treat a missing Rex block as exited) instead of polling a pane that is gone
