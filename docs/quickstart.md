# Quickstart

This walks the whole model loop by hand, with no agents and no daemon. You'll start Rhiza in a new repo, write down a capability and a decision, and take that change through propose, accept and land. It takes about ten minutes. You need `rhiza` on your `$PATH` and `git`; see [Install](../README.md#install) if you don't have it yet.

Every command in the `bash` blocks below is run by a test (`go test ./cmd/rhiza -run TestScripts/quickstart`), so the page can't quietly fall out of date with the CLI.

## Start a repo

Rhiza lives inside a git repo. Make one, with a first commit so there's something to hang things on, then run `rhiza init`:

```bash
git init -q
git commit -q --allow-empty -m "Start"
rhiza init
```

That creates `.rhiza/`, where the log goes. Everything Rhiza knows is an entry in that log, and only the CLI writes to it. Don't edit it by hand.

## Open a change set

Nothing goes into the model directly. Every edit belongs to a change set, which is one intention: a sentence saying what you're about to do.

```bash
rhiza change new "Export reports as CSV"
```

The new change set is a draft, and it's now the current one, so the commands that follow write into it.

## Add a capability and a decision

A capability says what the product does. A decision says why it's done that way, and who decided. Add them one command at a time, starting with a capability:

```bash
rhiza add capability reports --field "statement=Users see their reports"
rhiza add capability export-csv --parent reports --field "statement=Users export any report as CSV"
```

The name after the kind is the handle, a short name you can use in place of an ID. `--parent` puts `export-csv` under `reports` in the capability tree. The decision comes next, and it keeps the owner's own words in `quote`. If someone made the call, write down what they said, not your summary of it.

```bash
rhiza add decision csv-streaming --field "statement=Stream CSV rows instead of building the file in memory" --field 'quote="Stream it, reports get huge" (owner, 2026-10-20)'
rhiza link export-csv decided-by csv-streaming
```

That `decided-by` link is what ties the capability to its decision. `rhiza set` changes a field on a node later if you get the wording wrong.

When you have a lot to write, or an agent is doing the writing, `rhiza apply` takes a file of writes and puts them all in at once. You won't need it here.

Add `--dry-run` to a write to see what it would do first.

## Propose, accept, land

A change set moves through a few states. You propose it, someone accepts it, and once the work is done you land it.

```bash
rhiza propose
rhiza accept
```

Accepting is the owner's call, and here that's you. With agents in the picture, this is where you decide whether their proposal goes ahead.

Between accept and land is where the code would get written. There's none in this walkthrough, so land it straight away:

```bash
rhiza land
```

## Read it back

Once it's landed, the nodes are part of the model. `show` prints one node with its fields and links:

```bash
rhiza show export-csv
```

`why` goes the other way: the decisions, change sets and sources behind a node.

```bash
rhiza why export-csv
```

`slice` is the one to know if you work with agents. It composes the context around a node (its parents, its decisions, the constraints that apply) in a single read, and it's the first thing an agent runs before touching code. `tree` shows the capabilities as a tree.

```bash
rhiza slice export-csv
rhiza tree
```

## Render and commit

The log is the source of truth, but it's not much fun to read. `rhiza render` writes the model out as markdown under `.rhiza/rendered/`, so you can browse it and see it change in diffs. `rhiza check` then confirms the log, the rendered files and the repo agree.

```bash
rhiza render
rhiza check
git add -A
git commit -q -m "Export reports as CSV"
```

Commit the log and the rendered files along with your code. From then on the reasons for each decision travel with the repo.

## Where next

- The [README](../README.md) says what Rhiza is for and how to install it.
- `rhiza guide` prints the short reference that agents read, and `rhiza --help` lists every command.
- The [model](concepts/model.md) and [change sets](concepts/change-sets.md) pages cover the concepts behind all this, and [working with coding agents](agents.md) covers the agent side.
