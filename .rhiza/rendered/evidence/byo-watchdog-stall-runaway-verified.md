# byo-watchdog-stall-runaway-verified

evidence · `01M7RCZNABHNVFMFC52AHDZZX3` · actor [claude-code-ead76344](../actor/claude-code-ead76344.md)

## result

PASS

## statement

At f04be7546: go test -trimpath -p 2 ./internal/api -run 'TestStallLimit|TestWatchdog|TestRunaway' -count=1 exit 0 (capped scope); -v shows the new TestWatchdogPolicy/gate_step_gets_the_command_limit (StallLimit("gate",0,0)==CommandLimit, explicit 70m stall wins), TestWatchdogStalledBYO (myteam-flow parsed from docs/examples/byo/workflow.toml: work silent nudged at 21m then stopped and resumed with --retry at 41m; gate stopped and resumed at 31m past commandLimit, nothing done at 29m inside it) and TestWatchdogRunawayBYO (fake process table, 450 procs: work nudged with the reading, then run stopped and step work journalled with the runaway reason, never resumed, bug filed) all PASS. Scratch-built rhiza docs exit 0 with docs/reference/config.md unchanged (commandLimit described as covering command and gate steps); ship generate exit 0, 0 changed; render exit 0. Diff limited to the leaf's files plus its log/actor record; StallLimit switches on kind only, no workflow or step name special-cased, built-in stall= values untouched, no real process table read.

## Links

### Out

- verifies [byo-watchdog-stall-runaway](../capability/rhiza/workflows/workflow-runs/byo-watchdog-stall-runaway.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
