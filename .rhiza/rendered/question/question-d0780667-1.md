# question-d0780667-1

question · `01M6W4S1C1J0WRG7PHPBRQZS98`

## session

d0780667-0b40-44e7-a17d-01a1b849f0e6

## statement

Which intent should Scheduled triggers (01M6W4K) serve? It serves nothing yet. The repo has no .rhiza/policy.toml, so backlog work is built after milestone work by default. Workflow dispatch (01M6VY7G), which this change set depends on, serves milestone-25 (01M6HKVRE1B7, the backlog milestone). / Where does a repo declare its [[schedules]]? The intent says 'a repo's rhiza.toml', but rhiza.toml is the root's machine-wide config (internal/rootcfg). It isn't committed in any repo and is shared by several projects. Workflow dispatch moved its rules to a committed .rhiza/dispatch.toml for the same reason, and calls a schedule another trigger kind. / Where does a scheduled run work, given every run today is on a change set and its worktree (run start <cs>)? A schedule such as Nebari's 15-minute issue intake has no change set, and opening one per firing would flood the model with draft change sets. / How do scheduled runs fit the daemon's existing limits, and what can a schedule override? The intent's defaults: after downtime a schedule runs once, and a slot is skipped while the previous run is still going.

## Links

### In

- answers [scheduled-triggers-shape](../decision/scheduled-triggers-shape.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
