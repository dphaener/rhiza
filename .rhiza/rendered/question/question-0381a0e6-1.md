# question-0381a0e6-1

question · `01M6W4T069NY96AK9BGPRVHDEZ`

## session

0381a0e6-2443-4fec-a3aa-01a1b849f614

## statement

Risk warnings change set (01M6W4K4) serves no intent yet. Which should it serve? Its sibling Workflow dispatch (also from PR #353) serves milestone-25 (backlog: requests and polish ship on their own); milestone-29 (several projects under one root) is the other open one. / Rhiza has no schedule trigger of its own (a cron job would just run `rhiza request`), and autonomy is repo-wide (git config rhiza.autonomy.plan/ship), not per workflow. Which workflows should rhiza health and rhiza workflow check flag as 'outside text reaches agents with no permission checks'? / Where does the risk warning appear in rhiza health's report? / Workflow dispatch (01M6VY7G, accepted, not built) adds .rhiza/dispatch.toml rules that change which workflow a request or ready change set runs. Should Risk warnings wait for it?

## Links

### In

- answers [risk-warnings-shape](../decision/risk-warnings-shape.md)

## History

- `01M6W4K4DRTKBAFVPG4BC7TV51` Risk warnings, never blocks: rhiza health and rhiza workflow check warn when a workflow runs autonomously with bypassPermissions and can be triggered by outside input (a schedule, or requests carrying --ref), naming the workflow and saying outside text reaches agents with no permission checks; rhiza request --ref prints a one-line note when the workflow it will run is such a one. The workflows guide gains a section on outside input: what the warning means, that Rhiza cannot see what a scheduled command does, and the usual mitigations (a person gates intake, a narrower permission mode, owner signoff). Nothing is refused or fenced. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
