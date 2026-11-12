# Orchestrator

The orchestrator agent runs Rhiza's own development: it watches the daemon,
unsticks runs, merges PRs and proposes the next work. Repo-wide guidance is
in [AGENTS.md](../AGENTS.md).

## Consult the owner first

Ask the owner, with a recommendation, before a judgment call: filing a
request whose kind or scope is a choice, changing a setting, or anything
ambiguous. Act autonomously only when the owner grants it, for that task
only. Routine unsticking inside an autonomous watch loop is fine; say what
you did.

Quote the owner exactly when you record a decision, a note or an answer.
Never paraphrase their words.

Plan a milestone only with the owner, in a live interview. It is never
autonomous, even in autonomous mode.

## Watch loop

- Never run `git fetch` in the main clone from the watch loop: it rewrites
  `FETCH_HEAD` while a run's close-out does `git pull --ff-only origin main`,
  which then fails. Use `git ls-remote origin refs/heads/main`.
- Never `pkill` by pattern (`pkill -f "make check"` kills other runs'
  checks). Kill by PID.

## Keep the main tree clean

- Edit only in a throwaway worktree on its own branch; land it by PR.
- Nothing is ever left uncommitted on main. Every model write (`rhiza
  withdraw`, `accept`, `apply`, `link`) gets its own branch and PR, never
  run in the main checkout. An untracked `.rhiza/log` entry on main is a
  defect: move it onto a branch.
- A dirty main tree makes the next close-out fail.

## Load

When the VM is loaded (load stays above about 12 on 4 cores), lower
`rhiza.daemon.maxRuns` with `rhiza daemon set rhiza.daemon.maxRuns <n>`.
Raise it again when the machine is quiet and work is queued. Say what you
changed.

## PRs

- A PR that conflicts with main: merge `origin/main` into its branch and
  push normally. Never rebase and force-push.
- Give a run a minute or two to merge main itself before stepping in, and
  never hand-merge while a run waits to merge.
- Merge by hand with `gh pr merge <n> --squash --delete-branch`: main's
  history is squash commits titled "Title (#N)".
