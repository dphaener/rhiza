# risk-warnings-plan

source · `01M6W4YYQA62MHZR7NMMC0MHG7`

## quotes

```json
[]
```

## raw_sha256

06c1fb50c1cecd34b8ede81745f75e8b3181a8cac9abe18b979d0fd4c0384249

## session

0381a0e6-2443-4fec-a3aa-01a1b849f614

## summary

Autonomous plan run for change set Risk warnings (01M6W4K4); no owner present. The planner read the intent (warn, never block, when an autonomous workflow with bypassPermissions can be triggered by outside input), health's autonomy part, workflow.check, request --ref, autonomyStage (rhiza.autonomy.*), rhiza.daemon.permissionMode, adapter Bypass, and the accepted Workflow dispatch change set. It raised four questions in one: which intent to serve, what counts as exposed, where health shows it, and whether to wait for Workflow dispatch. The arbiter answered with Decision risk-warnings-shape: serve the backlog milestone; every queued request is outside input (a schedule is a cron job running rhiza request, invisible to Rhiza); flag the request workflow (plan) when plan autonomy is auto and the ready workflow (ship-change) when plan and ship autonomy are both auto, when any agent, check or arbiter step's effective mode (its permission_mode, else rhiza.daemon.permissionMode) is a bypass mode per its adapter; workflow check warns without changing valid or the exit code; health puts it in the autonomy part as warn with one detail per workflow and Next rhiza guide workflows; do not wait for dispatch, compute the reached workflows in one function that whichever change set lands second makes read .rhiza/dispatch.toml. The planner then wrote five leaves.

## title

Risk warnings plan, autonomous run

## Links

### In

- derived-from [risk-warnings-shape](../decision/risk-warnings-shape.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
