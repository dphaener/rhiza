# Workflows

A workflow is a TOML state machine of named steps. Each step ends with an outcome, and its
`on.<outcome>` names the next step. Write your own instead of editing the binary's.

## Where a name resolves

`rhiza run start --workflow <name>` takes the first of: your personal
`$XDG_CONFIG_HOME/rhiza/workflows/<name>.toml`, the repo's `.rhiza/workflows/<name>.toml`, then a
built-in. A path with a `/` or ending in `.toml` is that file. A higher layer shadows the lower.

## File

`name`, `start` (the first step), then `[step.<name>]` tables. A step name is lowercase words
joined by `-`, at most 30 characters. An unknown key is an error.

## Step kinds

- `agent`: `prompt`, `agent`, `permission_mode`, `models`, `until`, `until_when` (`idle` or
  `done`), `timeout`, `interval`, `max_visits`. Outcomes `pass`, `fail`.
- `wait`: `until`, `timeout`, `interval`. Outcomes `pass`, `timeout`.
- `command`: `run`, `timeout`. Outcomes `pass`, `fail`.
- `leaves`: the change set's leaves, a `[step.<name>.builder]` then a
  `[step.<name>.verifier]` (each `prompt`, `agent`, `permission_mode`, `models`), `split`
  (`auto`, `right`, `below`). Outcomes `pass`, `fail`.
- `checks`: `[[step.<name>.check]]` (`name`, `prompt`), `critic`, `fixer`, `verifier`.
  Outcomes `pass`, `fail`.
- `signoff`: waits for the owner. Outcomes `pass`, `changes`.
- `arbiter`: decides a question in an autonomous run. Not a state: nothing goes to it. At most
  one, only to set its agent, permission mode and interval.
- `workflow`: `workflow = "<name>"` and `skip = [steps]` runs another workflow's steps,
  named `<workflow>/<step>`. It takes no prompt, run or other session key.

Prompts may use `{change}`, `{slug}`, `{title}`, `{run}`, `{step}`; leaf and check steps add
`{leaf}`, `{attempt}`, `{check}`, `{cycle}`.

## Outcomes and targets

`on.pass = "next-step"`. The reserved targets `done` and `failed` end the run. A missing
`on.fail` fails the run; a step may loop back, bounded by `max_visits`.

## The loop

```
rhiza workflow ls
rhiza workflow new mine --from ship-change
rhiza workflow describe mine
rhiza workflow check mine
rhiza run start --workflow mine --dry-run
```

`ls` shows each layer and what it shadows. `new` forks a workflow into
`.rhiza/workflows/mine.toml` (without `--from`, a skeleton) and never overwrites. `describe`
shows the steps as the run resolves them. `check` reports every problem with the step it is in.
`--dry-run` starts nothing. Edit the file, check, repeat.

## Outside input

`rhiza health` and `rhiza workflow check` warn when a workflow runs autonomously in
`bypassPermissions` and outside input can reach it: a schedule, or a request, whose text and
`--ref` go to agents with no permission checks. Every queued request counts. Rhiza cannot see
what a scheduled command running `rhiza request` does. `rhiza request --ref` prints a note
when its workflow is such a one. Nothing is refused.

Mitigations:

- A person gates intake: `rhiza.autonomy.plan` `owner`.
- A narrower `permission_mode` on agent steps.
- Owner sign-off: `rhiza.autonomy.signoff` `owner`.
