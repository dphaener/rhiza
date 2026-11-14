# every-agent-session-plan-interview

source · `01M7AHRDMP54HEN975Q46FSARW`

## quotes

```json
[]
```

## raw_sha256

2f59cb39b3046056b58f6385c57697739deffc70ad477c8755c89321940f3423

## session

e520f785-7bad-4af3-ae27-01a1d514808d

## summary

Autonomous plan run for change set 'Every agent session a run' (milestone 36), with no owner present. The planner read the change set's three proposed leaves (run-brief-files, agent-invoke-template, guide-contract), the sibling change set 'The daemon offers spawn, gate and hand-back primitives' (whose run-handback leaf also said it writes handback.md), and the code: agent sessions start in internal/api (run.go StartAgent, run_leaves.go, run_checks.go, session_run.go launchEnv); prompts are templates rendered by internal/workflow; there is no per-run directory beside the journal, only the run's temp folder, which is removed when the run ends. Four questions went to the arbiter as one question: who owns handback.md, brief.md's format (left open by the milestone interview), where the files live, and where the invoke key applies. The arbiter decided (brief-files-and-invoke-placement): this change set owns handback.md's path, format and writer, and run-handback only triggers it; brief.md is deterministic Markdown with fixed sections (change set, leaves, decisions, constraints, check commands, a short contract summary ending with 'Full contract: rhiza guide contract'); the files live per session under <commondir>/rhiza/runs/<run id>/<step>[/<role or check>-<leaf>-<attempt>]/ and are kept after the run; invoke and handback are accepted wherever prompt is, with prompt as an alias and both on one table refused. The planner then split the brief writer from the daemon's wiring into every session start, so each leaf is one concern.

## title

Every agent session a run: autonomous plan interview

## Links

### In

- derived-from [brief-files-and-invoke-placement](../decision/brief-files-and-invoke-placement.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
