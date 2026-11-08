# Workflows

A workflow is a TOML state machine of named steps; an outcome picks the `on.<outcome>` step.

## Files

`--workflow <name>` takes the first of `$XDG_CONFIG_HOME/rhiza/workflows/<name>.toml`,
`.rhiza/workflows/<name>.toml`, a built-in. A path (`/` or `.toml`) is a file.

A file has `name`, `start`, `[step.<name>]` tables; step names are lowercase words joined
by `-`, at most 30 characters.

## Step kinds

- `agent`: `prompt`, `agent`, `permission_mode`, `models`, `until`, `until_when` (`idle` or
  `done`), `timeout`, `interval`, `max_visits`. Outcomes `pass`, `fail`.
- `wait`: `until`, `timeout`, `interval`. Outcomes `pass`, `timeout`.
- `command`: `run`, `timeout`. Outcomes `pass`, `fail`.
- `leaves`: `[step.<name>.builder]`, `[step.<name>.verifier]` (each `prompt`, `agent`,
  `permission_mode`, `models`), `split`. Outcomes `pass`, `fail`.
- `checks`: `[[step.<name>.check]]` (`name`, `prompt`), `critic`, `fixer`, `verifier`. Outcomes
  `pass`, `fail`.
- `signoff`: waits for the owner. Outcomes `pass`, `changes`.
- `arbiter`: decides questions in an autonomous run; at most one.
- `workflow`: `workflow = "<name>"`, `skip = [steps]`: its steps, named `<workflow>/<step>`.

Prompts take `{change}`, `{slug}`, `{title}`, `{run}`, `{step}`; leaf and check steps also
`{leaf}`, `{attempt}`, `{check}`, `{cycle}`.

## Targets and the loop

`on.pass = "next-step"`. `done` and `failed` end the run. A missing `on.fail` fails it;
loops are bounded by `max_visits`.

```
rhiza workflow ls
rhiza workflow new mine --from ship-change
rhiza workflow describe mine
rhiza workflow check mine
rhiza run start --workflow mine --dry-run
```

## Dispatch

`.rhiza/dispatch.toml` maps a kind to a workflow name in `[request]` (plan runs: `feature`,
`bug`, `milestone`) and `[ready]` (ready change sets: `feature`, `bug`).
Unset: `plan` and `ship-change`. `rhiza request --workflow <name>` is recorded on the change
set it opens. A ready run takes `run start --workflow`, else the request's, else `[ready]`,
else `ship-change`. An unresolved name never falls back: `rhiza daemon status`
shows a `workflow-unresolved` skip.

## Schedules

`[[schedules]]` in `.rhiza/dispatch.toml`: `name`; `every` (min `1m`) or `cron` (five
fields, local time or a leading `CRON_TZ=<zone>`); `workflow` or `command`; `catch_up`
(`once` or `skip` a slot missed offline); `overlap` (`skip` or `queue` behind a live run).

A scheduled run has no change set: it runs on the main checkout, leaves it as it found it,
and acts on the model only via paths such as `rhiza request`. `rhiza health` shows
schedules; the daemon log, runs.

## Outside input

`rhiza health` and `rhiza workflow check` warn when a request's text or `--ref` reach an
autonomous workflow in `bypassPermissions`. Rhiza cannot see what a scheduled `rhiza request`
sends. Nothing is refused. `rhiza.autonomy.plan` `owner` does not gate intake: requests are
planned at once under `rhiza.daemon.permissionMode`; the owner must only accept the plan.
Mitigations: gate intake outside Rhiza; narrow `rhiza.daemon.permissionMode` or
`permission_mode` on the plan workflow's steps; set `rhiza.autonomy.plan` and
`rhiza.autonomy.signoff` to `owner`.
