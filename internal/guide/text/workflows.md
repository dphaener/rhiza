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

`new` forks into `.rhiza/workflows/mine.toml` (without `--from`, a skeleton), never
overwriting. `check` reports every problem with its step. `--dry-run` starts nothing.

## Dispatch

`.rhiza/dispatch.toml` maps a kind (`feature`, `bug`, `milestone`) to a bare workflow name in
`[request]` (a queued request's plan run) and `[ready]` (an accepted change set's run), e.g.
`[ready]` then `bug = "bugfix"`. Any other table, key or value is an error naming it. Unset
kinds, or no file: `plan` and `ship-change`. No kind is `feature`.

`rhiza request --workflow <name>` is recorded, with the kind, on the change set the request
opens,.

A ready change set's run takes, in order: `run start --workflow`, the workflow its request
named, the `[ready]` rule, `ship-change`. A request's plan run takes the `[request]` rule.

A name that does not resolve never falls back: no run starts, and a `workflow-unresolved`
skip in `rhiza daemon status` and the inbox names the workflow and the rule or flag behind it.
