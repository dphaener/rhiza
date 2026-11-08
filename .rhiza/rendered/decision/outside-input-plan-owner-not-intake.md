# outside-input-plan-owner-not-intake

decision · `01M6W6524WEFTE8W9F4HM5W686`

## reasoning

I read internal/api/request.go: the daemon calls daemonStartPlan for each queued request regardless of rhiza.autonomy.plan, and planRunArgs passes the intent, --ref and --permission-mode=<daemon mode>, only omitting --autonomous when plan is owner. So the finding and critic are right: plan owner gates the plan's acceptance, not intake, and the guide as planned would tell owners a setting keeps outside text from bypass-mode agents when it does not. The change set's intent (warn-dont-safeguard: tell owners the true risk) is served by naming what actually reduces exposure: keeping intake human outside Rhiza, narrowing the permission mode of the agents that read the text, and owner acceptance/sign-off for later steps. Keeping rhiza.autonomy.plan in the section, correctly described, keeps the leaf's acceptance test meaningful without changing its scope.

## statement

The workflows guide's Outside input section must not present rhiza.autonomy.plan owner as gating intake. It says plainly that with plan owner the daemon still starts the plan run for every queued request at once, so the planner agent reads the request's text and --ref under rhiza.daemon.permissionMode; what plan owner adds is that the owner must accept the plan before anything after it runs. The mitigations the section lists are: a person gates intake outside Rhiza (only people queue requests, and no schedule runs rhiza request with outside text); a narrower permission mode for the agents that read outside text (rhiza.daemon.permissionMode, or permission_mode on the plan workflow's steps); owner acceptance of the plan (rhiza.autonomy.plan owner) and owner sign-off (rhiza.autonomy.signoff owner) for what runs after it. The leaf keeps naming rhiza.autonomy.plan, described this way; nothing is refused.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6w4k1-r3-f2](../question/review-01m6w4k1-r3-f2.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
