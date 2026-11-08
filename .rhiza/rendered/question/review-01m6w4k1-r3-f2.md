# review-01m6w4k1-r3-f2

question · `01M6W61C9M95QBDME0DENYA5VE`

## statement

Plan review round 3 left finding f2 (contradiction on workflows-guide-outside-input): It lists 'a person gates intake (rhiza.autonomy.plan owner)' as a mitigation. In the code, plan owner only drops --autonomous from the plan run the daemon starts for each queued request (internal/api/request.go:336, 387-402). The planner agent still starts right away with the request text and --ref, under rhiza.daemon.permissionMode. Plan owner gates the gate (the owner accepts the plan), not intake, so the guide would tell owners that a setting stops outside text from reaching bypass-mode agents when it does not. The critic confirmed it: Same code path: rhiza.autonomy.plan owner does not gate intake; the daemon still starts the plan run (request.go:336) with the request text and --ref under rhiza.daemon.permissionMode, and only the gate waits for the owner's acceptance (run_review.go runRunGate). Listing 'a person gates intake (rhiza.autonomy.plan owner)' as a mitigation would tell owners a setting keeps outside text from bypass-mode agents when it does not; the guide should name what actually helps (a narrower permission_mode for the plan workflow's steps or the daemon, owner acceptance for what runs after the plan). How should it be settled?

## Links

### Out

- blocks [workflows-guide-outside-input](../capability/rhiza/workflows/workflow-authoring/workflows-guide-outside-input.md)

### In

- answers [outside-input-plan-owner-not-intake](../decision/outside-input-plan-owner-not-intake.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
