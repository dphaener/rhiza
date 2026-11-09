# Runs and the daemon

Everything up to here has been you at a terminal. This page is about letting Rhiza do the typing: starting coding agents on change sets, keeping several going at once, and leaving it running while you're away. Nothing on this page is needed for the model itself. It's the part that needs Rex and Claude Code.

Read [Working with coding agents](agents.md) first. A run is that loop, driven by Rhiza instead of by you.

## What you need

Runs start agents in terminals, and Rhiza uses Rex for those. The agents are Claude Code. So before a run can do anything:

- `rex` is installed and its server is running, and you start Rhiza from inside it.
- `claude` is installed and logged in.
- Rhiza's skills are in Claude Code: `rhiza guide install --agent claude`.
- Your repository is a project under a Rhiza root. Runs, requests and the daemon work only inside one, and otherwise stop with `is not inside a Rhiza root`. A root is a folder holding your checkouts, with a `rhiza.toml`:

  ```text
  rhiza root init ~/code
  cd ~/code
  rhiza project add ~/code/myapp
  ```

  `project add` finds the root from the directory you run it in, so run it inside the root. It takes a checkout directly under the root, or a git URL to clone there.

Then connect Rhiza to Rex:

```text
rhiza rex install
```

That writes Rhiza's app actions and a key mode into your Rex config and reloads Rex. Rerun it after you upgrade Rhiza. If a `rhiza.lua` is there that Rhiza didn't write, it stops, and `--force` overwrites it.

The CLI doesn't change when you skip all this. `rhiza apply`, `land`, `check` and the rest still work on plain files.

## Workflows and runs

A workflow is a recipe: a TOML file of named steps, where each step's outcome picks the next one. A step can start an agent, run a command, wait for something, ask the owner to sign off, or include another workflow. A run is one go at a workflow, on one change set, with a journal of every step it took.

Rhiza ships with a few:

```text
rhiza workflow ls
```

- `plan` turns a request into a draft change set with leaves, and has it reviewed.
- `build` builds each leaf with a builder and has a verifier check it.
- `review` runs the review checks on what was built and fixes what they find.
- `ship-change` is `build` and then `review`, then opens the pull request, waits for CI, merges, and cleans up.

To start one on an accepted change set:

```text
rhiza run start <change> --workflow ship-change --permission-mode acceptEdits
```

That queues the run and returns. `--permission-mode` is the Claude Code permission mode the run's agent steps get when the workflow doesn't set one. The built-in workflows' builders set none, so without it `run start` stops and asks for one. A mode that bypasses permission checks, such as `bypassPermissions`, is refused unless the machine's `machine.toml` sets `trust = "dedicated"`. The daemon, below, picks it up. Leave the workflow off and you get `ship-change`. `--foreground` drives the run in your own terminal instead, which is handy once, to see how it behaves.

Each agent step runs in its own Rex tab, inside a git worktree of its own, so runs on different change sets never touch each other or your checkout. While one is going:

```text
rhiza run ls
rhiza run status <run>
rhiza run tail <run>
```

A run that stops, because you restarted the daemon or the machine went down, picks up where its journal left off with `rhiza run resume`. One you've given up on is `rhiza run abandon`.

### Sign-off

`review` and `ship-change` can wait for you at a sign-off step. Read the summary and decide:

```text
rhiza run approve <run>
rhiza run request-changes <findings.json> <run>
rhiza run send-back <run> --note "..."
```

Approving lets a `ship-change` run go on to merge. The other two send it back to a fixer. `request-changes` takes a JSON file of your findings (`-` reads stdin); `rhiza describe run.request-changes` shows its shape.

### Questions

An agent that isn't sure parks a question in the model. The run goes on with what it can, and the question waits for you in the inbox with everything else that needs the owner: sign-offs, failed runs and runs the daemon gave up on.

```text
rhiza inbox
rhiza question answer <question> "<your words>"
```

## The daemon

The daemon is what makes this continuous. It's the same `rhiza` binary. It polls the repo, starts runs for whatever is ready, resumes the ones that were interrupted, and watches for stalls. There's one per repository, and the CLI never needs it. It does only the continuous work.

```text
rhiza daemon start
rhiza daemon status
rhiza daemon stop
```

`start` opens it in a Rex tab. To have it run on its own, as a systemd user service that survives you logging out, use `rhiza daemon install` once, and `start` and `stop` will go through the unit from then on.

`rhiza daemon status` shows whether it runs, its settings, the live runs and what it last started. If you want a one-screen answer to "is anything wrong?", `rhiza health` reports the daemon, free disk, the usage hold, CI stability, the inbox and `rhiza check`, each as ok, warn or bad.

### Settings

`rhiza daemon settings` lists every setting with its value and where it came from. `rhiza daemon set <key> <value>` changes one, and a running daemon reads it at its next poll. The ones you'll reach for first:

- `rhiza.daemon.maxRuns` is how many runs go at once. The daemon doesn't know how loaded your machine is, so keep it modest.
- `rhiza.daemon.paused` stops it starting anything new.
- `rhiza.daemon.minFree` holds everything when the disk gets low.
- `rhiza.daemon.usage.session` and `rhiza.daemon.usage.weekly` are token limits for your plan, with `rhiza.daemon.usage.reserve` the share to keep back. The daemon holds new runs when it gets near them.
- `rhiza.daemon.permissionMode` is the permission mode agent steps get when they don't set one, in the runs the daemon starts itself, such as `plan` runs for requests. A run you queue with `rhiza run start` takes its own `--permission-mode` instead.

### Requests

Most work doesn't start as a change set. It starts as a sentence:

```text
rhiza request "Export reports as CSV"
rhiza request "Totals are wrong on the invoice page" --kind bug
```

That queues a plan request, and the daemon starts a `plan` run for it, which comes back as a draft change set for you to review and accept. `--kind` is `feature`, `bug` or `milestone`. `--ref` carries an issue URL along, which Rhiza stores and never reads. `rhiza request ls` shows the queue.

Which workflow handles what lives in `.rhiza/dispatch.toml`. Its `[request]` table names the workflow for each kind of plan run, and `[ready]` names the one that builds a ready change set. Anything unset uses `plan` and `ship-change`. A name that doesn't resolve is never quietly swapped for a default: `rhiza daemon status` shows it as skipped.

## Autonomy

Autonomy decides which gates wait for you. There's a setting per stage, each either `auto` or `owner`. Out of the box only `plan` is `owner`. `accept`, `signoff` and `ship` are `auto`, so a change set that passed its review is accepted, a run carries on past sign-off, and `ship-change` merges without you. Set them yourself if you want to be asked:

```text
rhiza daemon set rhiza.autonomy.plan owner
rhiza daemon set rhiza.autonomy.accept auto
rhiza daemon set rhiza.autonomy.ship auto
rhiza daemon set rhiza.autonomy.signoff owner
```

- `plan` is whether a finished plan goes straight to review and acceptance, or waits for you.
- `accept` is whether Rhiza accepts a change set that passed its review.
- `signoff` is whether a run stops at sign-off or carries on.
- `ship` is whether the merge happens without you.

An autonomous run is one started with `--autonomous`. It skips the owner's gates and questions and records each one it skipped, so you can see afterward what was decided for you. It never changes agent permissions. Those are the permission mode's job.

If you'd rather stay in the loop, set `owner` on `plan` and `signoff`. You still get the unattended building and the fixing. You just decide what goes in and what goes out.

### Per workflow

A workflow file can declare its own autonomy, one of two ways. A top-level `autonomous = true` (or `false`) sets all four stages to `auto` (or `owner`). An `[autonomy]` table sets only the stages it names:

```toml
[autonomy]
signoff = "owner"
ship = "owner"
```

Declare one or the other, not both. A stage the workflow declares overrides `rhiza.autonomy.<stage>` for runs of that workflow, and a stage it leaves out follows the setting. `rhiza run start --dry-run` prints each stage's mode and where it came from.

## Schedules

Some work happens on a clock: a weekly dependency check, a nightly sweep for stale questions. Put it in `.rhiza/dispatch.toml`:

```toml
[[schedules]]
name = "weekly-deps"
cron = "0 9 * * 1"
command = "rhiza request 'Update dependencies' --kind bug"
catch_up = "skip"
overlap = "skip"
```

Each schedule has a `name` and either `every` (a Go duration, at least `1m`) or `cron` (five fields, in local time, or `CRON_TZ=<zone>` in front). It runs either a `workflow` or a `command`. `catch_up` says whether a slot missed while the daemon was down runs once (`once`) or is dropped (`skip`). `overlap` says whether a slot that comes due while the last run is still going is dropped (`skip`) or waits its turn (`queue`).

A scheduled run has no change set. It works in your main checkout, leaves it as it found it, and touches the model only through things like `rhiza request`. That's why the example above files a request and doesn't edit anything itself. `rhiza health` lists the schedules, and the daemon log has their runs.

## What can go wrong

Runs build code from text, and some of that text comes from people other than you. A request's words, or the `--ref` it carries, end up in an agent's prompt. If that agent runs in `bypassPermissions`, which the daemon can use, and nothing stops the request, then whoever can file a request can steer an agent that has your shell. Rhiza doesn't refuse this. It warns:

- `rhiza health` and `rhiza workflow check` flag a request's text or `--ref` reaching an autonomous workflow in `bypassPermissions`.
- Rhiza can't see what a scheduled `rhiza request` sends, so a schedule gets no such warning.
- `rhiza.autonomy.plan` set to `owner` doesn't hold requests back. They're planned at once, under `rhiza.daemon.permissionMode`. What you control is that you must accept the plan.

Use bypassed permissions only on a machine that exists for this. To cut the exposure:

- Gate who can file requests, outside Rhiza.
- Narrow `rhiza.daemon.permissionMode`, or the `permission_mode` of the plan workflow's steps.
- Set `rhiza.autonomy.plan` and `rhiza.autonomy.signoff` to `owner`.

## Writing your own workflow

When the built-ins aren't what you want, copy one and change it:

```text
rhiza workflow new mine --from ship-change
rhiza workflow describe mine
rhiza workflow check mine
rhiza run start --workflow mine --permission-mode acceptEdits --dry-run
```

`workflow new` writes `.rhiza/workflows/mine.toml`. `describe` shows it as a run would resolve it, and `check` reports every problem and the step it's in, and starts nothing. A name resolves to your personal `$XDG_CONFIG_HOME/rhiza/workflows/`, then the repo's `.rhiza/workflows/`, then a built-in. The step kinds and their keys are in `rhiza guide workflows`.

## Where next

- [Working with coding agents](agents.md) covers the loop each agent step follows.
- [Change sets](concepts/change-sets.md) explains the states a run moves a change set through.
- `rhiza daemon settings` and `rhiza describe` list every setting and operation.
