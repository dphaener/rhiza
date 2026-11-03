# review-01m6fdsp-r3-f1

question · `01M6FFM086X8RJX4Q6GVV6HFF9`

## statement

Plan review round 3 left finding f1 (contradiction on daemon-headless-run-drivers): m24-only-agent-tabs quotes the owner: 'whatever is watching the plan, that should also be swallowed up by the daemon. No need to see anything in a Rex tab that isn't doing agent work.' But the daemon also starts plan runs for requests: daemonStartPlan (internal/api/request.go:305-328) opens a Rex tab named plan-<request> running `rhiza run start --workflow plan --new <intent> --request <id>`, so that tab holds the plan run's driver. daemon-headless-run-drivers covers only `rhiza run start --foreground <change set>` and `rhiza run resume --foreground <run>`, and its acceptance checks only a change set's run. No leaf in 01M6FENEJE4B11EBYXDGP9T052 moves request plan runs out of Rex. daemon-repo-session (01M6FEMDQ1D2HWXT3GR4F3Z5H6) even names 'plan-request tabs' and 'run tabs' as tabs the daemon opens. After the milestone, every plan request would still open a driver tab. Widen daemon-headless-run-drivers (statement and acceptance) to start request plan runs headless (run start --foreground --workflow plan --new ... --request ...), and drop the plan-request and run tabs from daemon-repo-session's wording. The critic confirmed it: daemonStartPlan (internal/api/request.go:305-328) opens a Rex tab plan-<request> running rhiza run start --workflow plan --new <intent> --request <id>, so that tab holds a plan run's driver. In 01M6FENEJE4B11EBYXDGP9T052, daemon-headless-run-drivers names only 'rhiza run start --foreground <change set>' and 'rhiza run resume --foreground <run>', and its acceptance checks only a change set's run. No plan leaf moves request plan runs out of Rex. daemon-repo-session still lists 'plan-request tabs, run and agent tabs' among the tabs the daemon or its runs open. Request plan runs would keep a driver tab, against m24-only-agent-tabs ('whatever is watching the plan, that should also be swallowed up by the daemon'). How should it be settled?

## Links

### Out

- blocks [daemon-headless-run-drivers](../capability/rhiza/orchestration/daemon/daemon-headless-run-drivers.md)

### In

- answers [answer-review-01m6fdsp-r3-f1](../decision/answer-review-01m6fdsp-r3-f1.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
