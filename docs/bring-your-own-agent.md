# Bring your own agent

Rhiza's runs start Rhiza's own skills: `/rhiza-builder` builds a leaf, `/rhiza-conduct` conducts a whole change set, and so on. Your team may already have skills of its own that do that work the way you like it. You can bring them in at three levels, each a bigger piece of the run:

1. A skill: your skill plays one role, say the builder, and Rhiza runs everything around it.
2. A conductor: your skill conducts the whole change set, and the daemon still gates what it hands over.
3. A workflow: your own steps from start to finish, put together from the [primitives](primitives.md).

Each level has a worked example under [`docs/examples/byo/`](examples/byo/). A test keeps all three passing `rhiza workflow check`, so you can copy one into `.rhiza/workflows/` and start from there. If you haven't written a workflow before, [writing a custom workflow](custom-workflows.md) is the place to start.

## Pick the agent

Every run, session and skill check launches the project's agent. `rhiza init` asks which one when you set the repo up, and you can name it instead:

```text
rhiza init --agent claude-code
```

The answer goes into the repo's git config as `rhiza.project.agent`. Change it later with `rhiza daemon set`, which refuses a name Rhiza doesn't know and lists the ones it does:

```text
rhiza daemon set rhiza.project.agent claude-code
```

An explicit `--agent`, or a step's own `agent` key, wins over the setting, and unset falls back to the built-in default. Claude Code is the only adapter today, so for now the agent you bring is your team's own Claude Code skills.

## Level 1: a skill

The smallest step is one role. A role is a session the run starts with a fixed contract: a builder ends with one commit carrying a `Leaf:` trailer, a verifier writes Evidence, a check records its findings. Your skill can play any of them, as long as it keeps the contract.

Check that first, before a real run depends on it. `rhiza skill check` runs your invoke as the role's session against a one-leaf fixture change set in a scratch repo, and reports each gate of the contract as pass or fail:

```text
rhiza skill check "/myteam:build Build leaf {leaf} of change set {change}, attempt {attempt}. Your brief is {brief}." --permission-mode acceptEdits
```

The role is `builder` unless you pass `--role verifier` or `--role check`. The permission mode is always explicit, and a mode that bypasses the agent's permission checks is refused. The scratch repo sits outside every git repository and is removed afterwards, so the current repo is never touched; `--keep` keeps it so you can look at what the session did. The command exits 1 when any gate fails.

Once it passes, put the invoke in a role table. Role tables are the `builder` and `verifier` of a `leaves` step and the `check`, `critic`, `fixer` and `verifier` of a `checks` step. `ship-change-classic` builds with those, so fork it:

```text
rhiza workflow new myteam-ship --from ship-change-classic
```

The classic includes the `build` workflow as one step. Replace that step with the `leaves` step written out, and name your skill in its builder table:

```toml
[step.build]
kind = "leaves"
on.pass = "render"

[step.build.builder]
invoke = "/myteam:build Build leaf {leaf} of change set {change}, attempt {attempt}. Your brief is {brief}."
permission_mode = "acceptEdits"

[step.build.verifier]
invoke = "/rhiza-verifier Verify leaf {leaf} of change set {change} ({title}) at commit $RHIZA_COMMIT: follow rhiza guide verifier."
```

The whole file is [skill.toml](examples/byo/skill.toml). It keeps Rhiza's verifier, so a session that didn't write the code still checks it, and it cuts the classic's path short to fit on a page; a real fork keeps the classic's conflict, gate and CI fix steps.

Why not `ship-change`? Its conductor starts every child with `rhiza run spawn`, and a spawned session is always told to follow `rhiza guide <role>`. There's no invoke to swap there, so one role of your own needs a workflow with role tables, or a conductor of your own.

## Level 2: a conductor

`ship-change` runs one conduct session that plans the leaves if there are none, builds and verifies each through spawns, runs the plan, code and QA checks, and lands the change set on its branch. The daemon owns the rest. To bring your own conductor, fork it and change the conduct step's prompt:

```text
rhiza workflow new myteam-conduct --from ship-change
```

```toml
[step.conduct]
kind = "agent"
invoke = "/myteam:conduct --prompt {brief} --slice {slice}"
handback = "/myteam:conduct --handback {handback} --prompt {brief} --slice {slice}"
until = "true"
until_when = "done"
on.pass = "gate"
```

Three placeholders carry the run's context, each the path of a file in the session's folder:

- `{brief}` is `brief.md`: the change set, its leaves, the decisions and constraints that apply, the check commands and a short summary of the contract every role follows.
- `{slice}` is `slice.json`, the model around the change set, the same bytes `rhiza slice --json` prints.
- `{handback}` is `handback.md`, what went red and why, written when a later step hands the run back. That's why the `handback` text names it.

The whole file is [conductor.toml](examples/byo/conductor.toml). Everything after `conduct` is the built-in's, which matters, because the daemon's gate still checks what your conductor hands over. It runs four commands, and your skill has to leave the branch so each one passes:

- `rhiza run build-check` passes when every leaf has Evidence from a verifier that wrote none of its commits. Your conductor still needs a second session to verify each leaf.
- `rhiza run review-check` passes when each check the step declares (`plan`, `code` and `qa` here) recorded a round with no open blocker or major finding. A declared check with no round fails it, so a check can't be skipped. Drop a `[[step.conduct.check]]` table if your team doesn't run that check.
- `rhiza check` passes on the model and the code.
- `rhiza ship gate` is the repo's local gate.

A red there, a send-back at sign-off, a refused push, a CI timeout, a red `pr-green` (a red CI, a conflict left by merging main, or a gate that failed on the merged head) or a red `rebase` all hand back to `conduct`, up to `handback_cap` times. A merge that fails doesn't hand back: it goes to `rebase`, which brings `origin/main` into the branch, gates and pushes the merged head and waits on CI again, handing back only for a code conflict, a gate that fails on the merged head or another error. Your skill reads `{handback}`, fixes what it names (never running `rhiza ship rebase` itself; the run's own steps do) and ends with `rhiza session done` again.

## Level 3: a workflow

The last level is your own workflow, not a fork. Every built-in is made of the same eight pieces, agent, spawn, gate, push/pr, wait, merge, sign-off and hook, and [workflow primitives](primitives.md) covers each of them. Pick the ones your team needs.

[workflow.toml](examples/byo/workflow.toml) is one way to do it. A command step starts the services, one agent step runs the team's `/myteam:implement` skill, which builds the change set and lands it, then a gate of the team's own commands checks it:

```toml
[step.gate]
kind = "gate"
commands = ["rhiza check", "make lint", "make test"]
on.pass = "signoff"
on.fail = "handback:work"
```

After the gate come the owner's sign-off, `rhiza ship open-pr`, a wait on CI, `rhiza ship merge` and `rhiza ship close-out`, with every red before the merge handing back to the `work` step. A `[[cleanup]]` hook stops the services however the run ends.

There's no verifier and no review checks in it, so there's no `rhiza run build-check` or `rhiza run review-check` in its gate either. Those belong to the runs that build with them. What your gate checks is up to you, but keep `rhiza check` in it, so the model and the code still agree when the run lands.

## Run it

All three run like any other workflow. Check one, then name it for a change set or map ready change sets to it in `.rhiza/dispatch.toml`:

```text
rhiza workflow check myteam-ship
rhiza change workflow myteam-ship
```

[Writing a custom workflow](custom-workflows.md#point-bugs-at-it) covers dispatch and pinning.
