# Workflow keys

Every key a workflow file takes, the default when it is unset, and what
`rhiza workflow new` writes. `rhiza guide workflows` covers resolution, dispatch and the loop.

## File keys

- `name`: unset takes the file's base name.
- `start`: the step a run enters; unset takes the first step in the file.
- `handback_cap`: how often a run may hand back (`on.<outcome> = "handback:<agent step>"`)
  before the next red parks it stopped with an inbox item; at least 1, unset 3. A hand-back
  gives every step its full `max_visits` again.
- `autonomous`: `true` or `false` sets every autonomy stage at once; not with `[autonomy]`.
- `autonomy`: a table giving stages `plan`, `accept`, `ship`, `signoff` `auto` or `owner`;
  a stage left out follows `rhiza.autonomy.<stage>`.
- `cleanup`: `[[cleanup]]`, repeated: `run` (required) and `name` (unset `cleanup-<n>`), run
  once when the run ends done, failed or abandoned; a failure never changes its result.
- `env`: `[env]`, `KEY = "value"` for every process the steps start; `{run_tmp}` is the
  run's private temp folder; `RHIZA_*` names are refused.
- `step.<name>`: a step; lowercase words joined by `-`, at most 30 characters, never `done`
  or `failed`.

## Shared step keys

- `kind` (required).
- `on.<outcome>`: a step, `done`, `failed` or `handback:<agent step>`.
- `max_visits`: unset or 0 means once.
- `stall`: such as `"20m"`; unset takes the watchdog's default for the kind.
- `retry_from`: the step a failed run restarts at when it fails here, by any retry; another
  step (not itself or the arbiter), `<workflow step>/<step>` for an included one, and prefixed
  like `on.*` targets when its workflow is included; unset restarts at the failed step. Not on
  a `workflow` or `arbiter` step.

## Step kinds

Each with its keys, sub-tables and outcomes.

- `agent`: `invoke` (required; `prompt` is an alias), `handback`, `agent` and
  `permission_mode` (unset: the run's), `models` (unset: the agent's default; the last stays
  for every later attempt), `until`, `until_when` (`idle`, the default, or `done`),
  `interval` (10s, or 30s with `until`). Outcomes `pass`, `fail`. An agent step declaring
  `[[step.<name>.check]]` is a conductor step: it spawns its own builders, verifiers and checks.
- `command`: `run` (required). Outcomes `pass`, `fail`.
- `wait`: `until` and `timeout` (required), `interval` (unset 10s), `max_interval`. Outcomes
  `pass`, `timeout`.
- `leaves`: `split` (`auto`, the default, `right` or `below`), `interval`;
  `[step.<name>.builder]`, `[step.<name>.verifier]`. Outcomes `pass`, `fail`.
- `checks`: `split`, `interval`; `[[step.<name>.check]]` (`name` plus the role keys),
  `[step.<name>.critic]`, `[step.<name>.fixer]`, `[step.<name>.verifier]`. Outcomes `pass`,
  `fail`.
- `signoff`: `on.changes` (required); waits for the owner. Outcomes `pass`, `changes`.
- `workflow`: `workflow` (required), `skip`. Outcomes `pass`, `fail`.
- `gate`: `commands` (required), run in order by the daemon, never an agent. Outcomes
  `pass`, `fail`.
- `arbiter`: `agent`, `permission_mode`, `interval`; at most one, its prompt fixed; without
  one a built-in arbiter is used. Outcomes `pass`, `fail`.

## Role keys

A role or check sub-table: `invoke` (required; `prompt` is an alias), `handback`, `agent`,
`permission_mode`, `models`, and for a verifier `low_risk_models` (unset: `models` for every
leaf) and `risky_paths`.

## rhiza workflow new

```
rhiza workflow new <name> [--from <workflow>]
```

With `--from` (a name, resolved personal, repo, then built-in, or a `.toml` path): a copy of
the resolved workflow with its `name` set. Without it, a skeleton: `start = "build"`, an
agent step `build`, and a command step `test` whose `run` is a placeholder that fails on
purpose, with `max_visits = 3` and `on.fail = "handback:build"`. It writes
`.rhiza/workflows/<name>.toml`, refuses an existing file and validates the result as
`rhiza workflow check` does.
