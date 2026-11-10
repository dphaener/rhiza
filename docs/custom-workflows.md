# Writing a custom workflow

The built-in workflows cover the usual path: plan, build, review, ship. Sooner or later you'll want a run to do something else. Maybe bugs shouldn't get the full feature treatment, or you want a command to run before an agent starts. A workflow is a small TOML file, so changing one is mostly copying and editing.

This page builds one from scratch: a `bugfix` workflow that reproduces a bug first, fixes it, runs the tests and waits for CI, and then points bug requests at it. If you haven't read how [workflows and runs](runs-and-daemon.md) work yet, skim the first half of that page.

## Start from a skeleton

```bash
rhiza workflow new bugfix
rhiza workflow check bugfix
```

`workflow new` writes `.rhiza/workflows/bugfix.toml` and never overwrites a file that's there. Add `--from ship-change` to fork a built-in instead of starting empty, which is the better move when you only want to tweak one. A skeleton already passes `check`, so you can run it after every edit and see what you broke.

## The workflow

Open the file and replace it with this:

```bash
cat > .rhiza/workflows/bugfix.toml <<'EOF'
name = "bugfix"
start = "reproduce"

[step.reproduce]
kind = "agent"
agent = "claude-code"
permission_mode = "acceptEdits"
prompt = "Change set {change} ({title}) is a bug. Write a failing test that reproduces it, commit it, then run `rhiza session done`. Do not fix anything yet."
on.pass = "fix"

[step.fix]
kind = "agent"
agent = "claude-code"
permission_mode = "acceptEdits"
prompt = "Make the failing test from the reproduce step pass for change set {change}, with the smallest change you can. Commit it, then run `rhiza session done`."
on.pass = "test"

[step.test]
kind = "command"
run = "make test"
max_visits = 3
on.pass = "ci"
on.fail = "fix"

[step.ci]
kind = "wait"
until = "gh pr checks"
timeout = "30m"
interval = "30s"
on.pass = "done"
on.timeout = "failed"
EOF
rhiza workflow check bugfix
```

The file is a state machine. `start` names the first step, and each `[step.<name>]` table is one step. When a step finishes it has an outcome, and the step's `on.<outcome>` key says where to go next. `done` and `failed` end the run.

Here's how this one reads:

- **reproduce** starts an agent with a prompt. `{change}` and `{title}` are filled in from the change set, and the agent ends with `rhiza session done`, which is what makes the step pass. Asking for a failing test first keeps the agent from "fixing" something it never saw break.
- **fix** is a second agent step, so a fresh agent session. It sees the committed test and makes it pass.
- **test** runs a plain command. Pass moves on, fail loops back to `fix`. A loop needs a bound, so `max_visits = 3` fails the run on the fourth visit instead of letting the two steps go back and forth all night. Swap `make test` for your repo's own scoped command.
- **ci** waits. It runs `gh pr checks` every 30 seconds until it exits cleanly, and gives up after 30 minutes. Passing there ends the run.

Every step kind and its keys, including the ones this example doesn't use (`leaves`, `checks`, `signoff`, `arbiter`, and `workflow` to include another workflow), are in the generated workflow-file reference, `docs/reference/workflow-file.md`, and in `rhiza guide workflows`. Prompts can use `{change}`, `{slug}`, `{title}`, `{run}` and `{step}`, plus a few more in the `leaves` and `checks` steps.

## Check it, then look at it

```text
rhiza workflow describe bugfix
rhiza workflow check bugfix
```

`describe` prints the workflow as a run would resolve it: where the file came from, what it shadows and every step's targets. `check` is stricter. It reports every problem at once with the step it's in, say a target that isn't a step, a misspelled outcome or an agent Rhiza doesn't know, and exits non-zero. It starts nothing, so it's safe in a script or a pre-commit hook.

To see what a run would do without queuing one:

```text
rhiza run start --workflow bugfix --dry-run
```

## Point bugs at it

A workflow only runs when something asks for it. You can name it for one run with `rhiza run start <change> --workflow bugfix`, but the more useful move is to map bug requests to it. `.rhiza/dispatch.toml` does that:

```bash
cat > .rhiza/dispatch.toml <<'EOF'
[ready]
bug = "bugfix"
EOF
```

Now a change set that came in as a bug (`rhiza request "Totals are wrong on the invoice page" --kind bug`) gets built by `bugfix` once you've accepted its plan. Features still get `ship-change`. The `[request]` table in the same file does the same for the planning run, if you want a different workflow on that side too. Both default to the built-ins, and a name that doesn't resolve is never swapped for a default: `rhiza daemon status` shows the run as skipped.

## Pin one change set to it

Sometimes you don't want a rule, just this one change set on this one workflow. Name it when you open the change set:

```text
rhiza change new "Totals are wrong on the invoice page" --workflow bugfix
```

Or pin one that already exists, and clear the pin later if you change your mind:

```text
rhiza change workflow bugfix
rhiza change workflow --clear
```

A pinned change set gets that workflow as its ready pick and skips `dispatch.toml` entirely. Clear the pin and the dispatch rules choose again.

## Where it lives

A name resolves to your personal `$XDG_CONFIG_HOME/rhiza/workflows/` first, then the repo's `.rhiza/workflows/`, then a built-in. So a repo file named `ship-change` replaces the built-in for everyone on the repo, and a personal one replaces both just for you. That's handy, but it's also easy to lose track of which file ran. `rhiza workflow describe <name>` says which one won.

If you build a workflow that runs without anyone at the keyboard, watch for an agent step in `bypassPermissions` that a request's text can reach. `rhiza workflow check` warns about it, and [what can go wrong](runs-and-daemon.md#what-can-go-wrong) has the details.

## Where next

- [Runs and the daemon](runs-and-daemon.md) covers requests, autonomy and schedules, which all take a workflow name.
- [Working with coding agents](agents.md) covers the loop each agent step follows.
