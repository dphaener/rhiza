# Workflows

A workflow is a TOML state machine of named steps. Each step ends with an outcome; its
`on.<outcome>` names the next step. Write your own instead of editing the binary's.

## Where a name resolves

`--workflow <name>` takes the first of `$XDG_CONFIG_HOME/rhiza/workflows/<name>.toml`,
`.rhiza/workflows/<name>.toml`, a built-in. A path with `/` or ending `.toml` is that file.

## File

`name`, `start` (the first step), `[step.<name>]` tables. Step names: lowercase words joined
by `-`, at most 30 characters. An unknown key is an error.

## Step kinds

- `agent`: `prompt`, `agent`, `permission_mode`, `models`, `until`, `until_when` (`idle` or
  `done`), `timeout`, `interval`, `max_visits`. Outcomes `pass`, `fail`.
- `wait`: `until`, `timeout`, `interval`. Outcomes `pass`, `timeout`.
- `command`: `run`, `timeout`. Outcomes `pass`, `fail`.
- `leaves`: `[step.<name>.builder]`, `[step.<name>.verifier]` (each `prompt`, `agent`,
  `permission_mode`, `models`), `split`. Outcomes `pass`, `fail`.
- `checks`: `[[step.<name>.check]]` (`name`, `prompt`), `critic`, `fixer`, `verifier`.
  Outcomes `pass`, `fail`.
- `signoff`: waits for the owner. Outcomes `pass`, `changes`.
- `arbiter`: decides questions in an autonomous run; no state, at most one, sets its agent.
- `workflow`: `workflow = "<name>"`, `skip = [steps]`: another workflow's steps, named
  `<workflow>/<step>`.

Prompts take `{change}`, `{slug}`, `{title}`, `{run}`, `{step}`, and in leaf and check steps
`{leaf}`, `{attempt}`, `{check}`, `{cycle}`.

## Outcomes and targets

`on.pass = "next-step"`. `done` and `failed` end the run. A missing `on.fail` fails it; a
loop back is bounded by `max_visits`.

## The loop

```
rhiza workflow ls
rhiza workflow new mine --from ship-change
rhiza workflow describe mine
rhiza workflow check mine
rhiza run start --workflow mine --dry-run
```

`new` never overwrites; `describe` shows resolved steps; `check` names each problem's step.

## Dispatch

`.rhiza/dispatch.toml` maps a kind to a workflow name in `[request]` (a request's plan run:
`feature`, `bug`, `milestone`) and `[ready]` (an accepted change set's run: `feature`, `bug`).
Unset: `plan` and `ship-change`. `rhiza request --workflow <name>` is recorded on the change
set the request opens. A ready run takes `run start --workflow`, else the request's
workflow, else `[ready]`, else `ship-change`. A name that does not resolve never falls back:
a `workflow-unresolved` skip in `rhiza daemon status` and the inbox names it.

## Outside input

`rhiza health` and `rhiza workflow check` warn when a request's text and `--ref` reach an
autonomous workflow in `bypassPermissions`. Rhiza cannot see what a scheduled
`rhiza request` sends. Nothing is refused.
`rhiza.autonomy.plan` `owner` does not gate intake: requests are planned at once under `rhiza.daemon.permissionMode`; the owner must
only accept the plan. Mitigations: gate intake outside Rhiza; narrow
`rhiza.daemon.permissionMode` or `permission_mode` on the plan workflow's steps; set
`rhiza.autonomy.plan` and `rhiza.autonomy.signoff` to `owner`.
