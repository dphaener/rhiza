# agent-step-until

decision · `01M5WJJGTMG3YS0TYDVTNHMXX4`

## statement

An agent step may have an until command, as a wait step does. The step passes as soon as the session is idle at its prompt (not blocked, not working) and until exits 0; it still passes on exit 0 and fails on a non-zero exit. The session is then left open: the tab stays, and the owner can still read it. until is checked at the step's interval (default 30s) only while the session is idle; a blocked session is never passed over. In ship-change, build gets until = "rhiza ship check --pr-open", meaning a PR is open from cs/{slug}.

## wording

the orchestrator's ruling, under the owner's authorization to amend plans

## Links

### In

- decided-by shipit-agent-step (retired)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
