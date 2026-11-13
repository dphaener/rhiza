# Build retry

A retry adds to the builder topic; read that first. It holds when `RHIZA_ATTEMPT` is
above 1, and only what a retry adds is here.

## Inspect the working tree

Earlier attempts left partial work. Before you write anything, inspect the working tree
and the earlier attempts' commits for it: what is built, what is half done, what the
verifier saw. Build on it; do not redo it.

## Read the findings

Read the previous attempt's findings: in the prompt when they are there, else
`rhiza run status --json`. Address every finding exactly as written, one by one. Do not
re-diagnose, narrow or skip one.

## Commit

Add one fix commit with the same `Leaf:` trailer as before. Never rewrite history: no
amend, rebase or force push over the earlier commits.
