# Working with coding agents

Rhiza doesn't care which coding agent you use. An agent is anything that can read a prompt, run shell commands and make git commits, and everything it gets from Rhiza or hands back goes through the CLI. This page covers setting one up and what the loop looks like once it's running.

If you haven't done the [quickstart](quickstart.md) yet, do that first. It's the same loop, just done by hand.

## Install the guide

The agent has to be told Rhiza exists. `rhiza guide` prints a short reference written for it (the vocabulary, the core loop, the hard rules and a few recipes), and you can also get it a topic at a time:

```text
rhiza guide ls
rhiza guide build
```

Most agents do better with that as a skill they load themselves than with a command they have to remember. For Claude Code, `rhiza guide install` writes the skills for you:

```text
rhiza guide install --agent claude
```

That puts them in `~/.claude/skills/`, so every repo gets them. Add `--project` to install into this repo's `.claude/skills/` instead, which is the better choice if you want the skills to travel with the repo. Rerun it after you upgrade Rhiza, because the guide text changes with the CLI. If a folder is already there that Rhiza didn't write, it stops, and `--force` overwrites it.

Another agent can read the same text. Point it at `rhiza guide`, or paste that output into whatever instructions file it uses.

## The loop

Every task follows the same shape. The agent works on an accepted [change set](concepts/change-sets.md) and keeps to what that change set intends.

1. Read the slice. Before touching any code, run `rhiza slice --change "$RHIZA_CHANGE" --json`. It composes the context around the work in one read: the capability, its parents, the decisions and the constraints that apply. `rhiza ready` names the next accepted change set if you don't have one in hand.
2. Do the work. Write the code and its tests, the same as you would without Rhiza.
3. Write back what changed in the model. If the work added a capability, settled a decision or turned up a question, it goes back in one `rhiza apply <file> --json`, with everything in it. That's the batch path for agents. When you're the one writing the model by hand, you use plain commands instead, like `rhiza add question export-size-cap --field statement="Do exports need a size cap?"`. Anything outside the change set's intent stays out of the apply. It gets a new change set linked `discovered-from` this one, or a question.
4. Commit. The commit carries trailers, covered below.
5. Land. `rhiza land --json` makes the model changes part of the model, then `rhiza render` refreshes `.rhiza/rendered/`, and that goes in the last commit.

If the agent isn't sure about something, it should add a question and not guess. The doubt sits in the model where the owner will see it, and the agent carries on with whatever it can still finish.

## Commit trailers

Every commit an agent makes ends with trailers that tie it back to the model:

```text
Stream CSV exports

Change: 01M5PWGEF85Y185JZ2ZGCNH40N
Actor: claude-code-1a2b3c4d
```

`Change:` is the ID of the change set. `Actor:` is the handle of whoever made the commit, and it comes from `RHIZA_ACTOR` in the agent's environment. The actor has to exist as a record in the repo, so when an agent starts it writes its own actor record, and that record gets committed with its first piece of work.

When a commit builds one leaf of a larger change set, there's a third trailer, `Leaf: <handle>`. Any `Co-Authored-By:` line goes in the same paragraph, with no blank line between the trailers. Git only reads the last paragraph as trailers, so a blank line in the middle quietly drops the ones above it.

The trailers are what let `rhiza check` and the rest of the tooling say which change set, and which actor, a given piece of code came from.

## Rules the agent follows

These are in the guide too. They're what you're relying on, so they're worth knowing:

- Never edit `.rhiza/log/` by hand. Only the CLI writes it.
- Always pass `--json` and read that, not the human text.
- Propose, and leave `rhiza accept` to the owner unless your policy says otherwise.
- Never paraphrase the owner. A decision's `quote` is their exact words, and an option the agent wrote is the agent's wording, not theirs.

## Where next

- [Change sets](concepts/change-sets.md) explains the states the loop moves a change through.
- [Decisions, sources and why](concepts/why.md) covers questions and quotes.
- `rhiza guide ls` lists every topic the agents read.
