# Workflow primitives

Every workflow, the built-in `ship-change` included, is put together from eight pieces: agent, spawn, gate, push/pr, wait, merge, sign-off and hook. Each one is either a step kind or a `rhiza` command a step runs, and each ends with an outcome that picks the next step. This page says what each piece does, what it's made of and what happens when it goes red. The keys themselves are in the generated [workflow file reference](reference/workflow-file.md), so they're not repeated here.

If you haven't written a workflow yet, [writing a custom workflow](custom-workflows.md) builds one from scratch and is the better place to start.

## How a red hands back

Most of these pieces can fail, and the useful thing to do with a failure is usually to give it to an agent that can fix it. That's a hand-back. Instead of naming a step in `on.fail`, you name an agent step with a `handback:` prefix:

```toml
handback_cap = 3

[step.gate]
kind = "gate"
commands = ["rhiza check", "make test"]
on.pass = "signoff"
on.fail = "handback:conduct"
```

When the gate goes red, the run writes `handback.md` into the `conduct` step's session folder and sends that session the step's `handback` text. A live session gets it as a message; one that already ended is resumed in the same conversation. `handback.md` names the failing step and its outcome, then whatever the run knows about why: the failing commands and the files holding their full output, the output's tail, the failed CI runs or the conflicted files. The agent reads it through `{handback}` in its template or `$RHIZA_HANDBACK`. When the agent step passes again, the run goes on from there to wherever its `on.pass` points.

`handback_cap` bounds this. It counts hand-backs across the whole run (unset is 3), and the red after the last one parks the run stopped, with an inbox item naming the step and the cause. Each hand-back also gives every step its full `max_visits` again. If you route a red outcome back to an agent step by its plain name instead, nothing caps it but `max_visits`, no `handback.md` is written and nothing reaches the inbox. `rhiza workflow check` warns about that, since it's almost always a mistake.

## Agent

An `agent` step starts one agent session with a prompt and waits for it. It's the only piece that writes code.

`invoke` (or its alias `prompt`) is the text the session starts with. It's a template: `{change}`, `{title}`, `{run}` and the rest are filled in, and `{brief}`, `{slice}` and `{handback}` are the paths of the session's context files (`rhiza guide contract` lists them all). `handback` is the text a hand-back sends, and an agent step that's the target of a `handback:` route needs one. `agent`, `permission_mode` and `models` pick who runs it and how; unset, they take the run's.

Outcomes are `pass` and `fail`. The step passes when the agent says it's done with `rhiza session done` or exits 0, and fails when it exits any other way. With `until`, the step also passes once that command exits 0 while the session sits idle (or, with `until_when = "done"`, only once the agent said it's done), and the session is left open. That's how `ship-change`'s `conduct` step stays around for later hand-backs.

An agent step is where hand-backs land, so it's the one piece that doesn't hand back itself. A failed agent step goes wherever its `on.fail` says, and with no route it ends the run failed.

## Spawn

`rhiza run spawn` isn't a step kind. It's the command an agent inside a run uses to start a child session, which is how one conductor session builds a whole change set: it spawns a builder per leaf, a verifier for each, then the checks.

```bash
rhiza run spawn --role builder --leaf export-csv --brief brief.md
```

The child starts in the run's worktree through the run's agent adapter and permission mode, with a prompt telling it to read its brief and follow `rhiza guide <role>`. Its `brief.md` is the contract section followed by the file you passed. `--leaf` gives it a leaf to work on, `--model` picks its model and `--check` names the review check it works on (a `check-<name>` role gets that one without the flag). The spawn is journaled as a child of the calling step and doesn't count against `rhiza.daemon.maxRuns`. Only a usage hold or the machine load hold makes it wait, and that wait is journaled too.

The command blocks until the child exits. Its outcome is `pass` when the child exited 0 and `fail` otherwise, and `rhiza run spawn` exits 1 on a fail. A red spawn doesn't hand anything back by itself: the agent that spawned it decides what's next, usually a retry or a fixer. If that agent gives up and its own step fails, the step's `on.fail` takes over.

## Gate

A `gate` step runs a list of `commands` in order, in the run's worktree. They're commands the daemon owns, never an agent, so an agent can't talk its way past one.

Every command runs even after one fails, and the step passes only when all of them exit 0. Outcomes are `pass` and `fail`. Each failing command's full output is kept in a file in the run folder, and the step's detail names it. That's what makes a gate a good hand-back source: `handback.md` lists each failing command with its output file.

`ship-change`'s gate runs `rhiza run build-check` (every leaf has Evidence from a verifier who didn't build it), `rhiza run review-check` (no open blocker or major finding), `rhiza check` and `rhiza ship gate` (the repo's local gate), and hands a red back with `on.fail = "handback:conduct"`.

## Push/PR

Pushing and opening the pull request is `rhiza ship open-pr`, run by a `command` step:

```toml
[step.push]
kind = "command"
run = "rhiza ship open-pr"
on.pass = "wait-ci"
on.fail = "handback:conduct"
```

It pushes the change set's branch and opens its pull request, titled by the change set with the run's review summary as the body. If one is already open, it pushes and replaces the body. The repo's pre-push check (`rhiza.ship.prePushCheck`) runs inside it, before the push. A command step passes on exit 0 and fails otherwise, so a refused push or a red pre-push check is a `fail`, and `ship-change` hands that back.

## Wait

A `wait` step runs its `until` command every `interval` until it exits 0, and gives up after `timeout`. Outcomes are `pass` and `timeout`. With `max_interval`, each check whose output didn't change doubles the sleep up to that, and a changed output resets it, which keeps a long CI wait from hammering GitHub.

`ship-change` waits on CI with `rhiza ship check --settled`, then checks it with a `pr-green` command step:

```toml
[step.wait-ci]
kind = "wait"
until = "rhiza ship check --settled"
timeout = "2h"
interval = "30s"
max_interval = "2m"
max_visits = 6
on.pass = "pr-green"
on.timeout = "handback:conduct"

[step.pr-green]
kind = "command"
run = "rhiza ship check"
max_visits = 6
on.pass = "merge"
on.fail = "handback:conduct"
```

`--settled` passes once every check has finished, green or red, so a red CI ends the wait as soon as it's in rather than at the timeout. When GitHub reports the pull request conflicting, the wait merges `origin/main` into the branch: a clean merge is pushed and waited on, and a code conflict is left in progress and settles the wait. `pr-green` passes only when every check is green and records the head it tested, so a red CI, a conflict left in progress or a merged head the local gate failed on fails there and hands back. A timeout hands back with `on.timeout = "handback:conduct"`.

## Merge

The merge is `rhiza ship merge`, again run by a `command` step, after `pr-green` passed. A failed merge goes to a `rebase` step rather than handing back:

```toml
[step.merge]
kind = "command"
run = "rhiza ship merge"
max_visits = 6
on.pass = "main-ci"
on.fail = "rebase"

[step.rebase]
kind = "command"
run = "rhiza ship rebase"
stall = "70m"
max_visits = 6
on.pass = "wait-ci"
on.fail = "handback:conduct"
```

`rhiza ship merge` squash-merges the pull request, pinned to the head whose checks passed. If `origin/main` has moved past that head it merges nothing and fails, and `rebase` brings main in: `rhiza ship rebase` merges `origin/main` into the branch, settles rendered-only conflicts by rendering again, runs the local gate and the pre-push check on the merged head, pushes it and goes back to `wait-ci` for a CI round on the new head (sign-off isn't asked again). When main hasn't moved it merges nothing, and the wait, check and merge retry. It hands back only on a code conflict left in progress or a local gate or pre-push check failing on the merged head, so the conductor spawns the conflict or gate fixer and never runs `rhiza ship rebase` itself.

A command step can also stop the run instead of failing it: exit 75 and the run stops with the last `stop:` line of the output, waits for the owner, and picks up at the same step on `rhiza run resume`. `rhiza ship merge` does that for a draft pull request or a merge GitHub refused, since no agent can fix either.

After the merge, nothing on the branch can fix the change set, so `ship-change`'s later steps (`main-ci`, `close-out`, `ship`, `deploy-check`) have no `on.fail` route and no hand-back. A red there stops or fails the run.

## Sign-off

A `signoff` step runs nothing. It stops the run until the owner signs off. Outcomes are `pass` and `changes`, and `on.changes` is required.

`rhiza run approve` passes it. `rhiza run send-back` (with a note) or `rhiza run request-changes` (with a findings file) gives `changes`, and the owner's words become the findings. In `ship-change`, `on.changes = "handback:conduct"`, so the conductor gets them in `handback.md` like any other red. When `rhiza.autonomy.signoff` is `auto`, or the run is autonomous, the run carries on past the step without waiting.

## Hook

A hook is any command you want run around the agents: setting up a database, seeding a cache, posting a message. There are two places for one.

A `command` step runs `run` as a shell command at its place in the workflow. It passes on exit 0 and fails otherwise, and its `on.fail` can hand back like any other red, with the output's tail in `handback.md`.

`[[cleanup]]` runs once when the run ends, done, failed or abandoned, in the run's worktree with the run's environment, in file order. A failing cleanup is logged and never changes the run's result, so there's nothing to hand back:

```toml
[[cleanup]]
name = "stop-services"
run = "docker compose down"
```

Both get the workflow's `[env]` variables, `{run_tmp}` included.

## The other step kinds

Four more step kinds aren't primitives so much as bundles of them, kept for the older flows:

- `leaves` runs a builder and then a verifier session for each leaf, retrying a failed one. It's what `ship-change-classic` builds with; `ship-change` does the same through spawns.
- `checks` runs the review checks in parallel, then a critic and fixers over what they found.
- `workflow` includes another workflow's steps inside this one.
- `arbiter` isn't a state of the machine at all. In an autonomous run it answers a question an agent step raised.

The full list, `agent`, `command`, `wait`, `leaves`, `checks`, `signoff`, `workflow`, `gate` and `arbiter`, with every key, is in the [workflow file reference](reference/workflow-file.md). `rhiza workflow describe ship-change` shows how the built-in puts the pieces together.
