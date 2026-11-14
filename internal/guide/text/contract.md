# Agent contract

What any agent in a run must do, and the gate that checks it. Your role topic adds to
this page; it never removes from it.

## Obligations and their gates

- Write the model only through `rhiza` (never edit `.rhiza/log`): `rhiza check` and the fold refuse a hand-edited log.
- A builder's or fixer's commit carries `Change:`, `Actor:` and `Leaf:` (or `Review:`) trailers in the final paragraph: the build gate refuses a commit without them.
- Every leaf ends with Evidence written as the verifier: the leaf is not done without it.
- Report findings only through `rhiza run findings`: the review gate reads nothing else.
- Ask the owner only through `rhiza` questions: an unrecorded question is never answered.
- End with `rhiza session done` as its own command: the run waits for it and fails the attempt without it.

## Invoke placeholders

Every `invoke` or `handback` template in a workflow file takes `{change}`, `{slug}`,
`{title}`, `{run}`, `{step}`, `{brief}`, `{slice}`, `{handback}` (the session's file
paths), `{worktree}`, `{leaf}` and `{attempt}`; `{leaf}` and `{attempt}` render empty
where the session has none (an agent step, a check or critic, and `{leaf}` for a fixer
whose findings name no leaf). A checks step's fixer adds `{cycle}`; a checks step's check
and critic add `{check}` and `{cycle}`. `prompt` is an alias of `invoke`; `rhiza workflow check` refuses an
unknown placeholder.

## Environment

`RHIZA_CHANGE`, `RHIZA_RUN`, `RHIZA_ACTOR`, `RHIZA_SESSION`, `RHIZA_LEAF`,
`RHIZA_ATTEMPT`, `RHIZA_ROLE`, `RHIZA_CHECK`, `RHIZA_COMMIT`, `RHIZA_ROOT` and
`RHIZA_AUTONOMOUS` are set when they apply to your role.
`RHIZA_BRIEF`, `RHIZA_SLICE` and `RHIZA_HANDBACK` name the context files below.

## Context files

In the session's own directory under the runs directory:

- `brief.md`: your task, built from the model (`RHIZA_BRIEF`).
- `slice.json`: `rhiza slice --json` for the change set when the session started (`RHIZA_SLICE`).
- `handback.md`: why the session was handed back, only after a hand-back (`RHIZA_HANDBACK`).

A variable is set only when its file exists; a session outside a run has none.
