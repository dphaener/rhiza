# build-workflow-plan

source · `01M5YX1W7CJR6Q5S775HFJ52CN`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I like 1, but let's be smart about splitting panes. For instance, I am using a monitor in portait mode right now so a different split strategy may be desired.",
    "where": "build workflow plan session, topic 5: panes"
  }
]
```

## raw_sha256

24765456d6fcc122838d7f707ff32e5189fa93bc7cfe8b77f4bed4d273221040

## session

ab54d544-8c72-4c39-ab01-01a17dc7f611

## summary

The plan step of a plan run on change set 6 of milestone 3 (The build workflow), 2026-10-27. Starting from the settled milestone-3 decisions (build-workflow-shape, verifier-writes-evidence, full-agent-sessions, engine-is-conductor, one-worktree-first) and the one leaf already written (run-program-status), the planner put six topics to the owner, each through AskUserQuestion with a recommendation; the owner picked the recommended option in each. Settled: (1) the per-leaf loop is a new `leaves` step kind in the workflow file, with builder and verifier sub-tables, conducted and journaled by the engine; (2) each verifier runs in a throwaway detached worktree at the leaf's commit, its Evidence written into the build worktree's log, so verifiers stay parallel; (3) 3 builder attempts in all, a model ladder per attempt (default sonnet, opus, opus; verifier opus), the Claude adapter gaining a model flag, and each retry a fresh builder session that reads the verifier's journaled findings; (4) a retry adds a fix commit with the same Leaf: trailer, and the engine commits only the Evidence entry's path as its own chore(rhiza) commit; (5) one tab per build, builder pane and verifier panes, a pane closed when its session passes and failed or blocked panes kept until the run ends; on the owner's typed ask, splits follow the pane's shape in pixels (along the longer side, the builder keeping half), with a `split` override in build.toml; (6) a failed leaf lets independent leaves carry on and the run ends failed, a resume after the leaf's question is answered gives it a fresh round (in autonomous mode the arbiter answers and the engine gives one fresh round); needed outside work is a question, noticed-only work a discovered-from draft change set; the run ends with rhiza run build-check and does not land. The picks are captured as Decisions by session interview; the leaves link the landed milestone-3 decisions.

## title

The build workflow's leaves, planned with the owner

## Links

### In

- derived-from [build-workflow](../capability/rhiza/workflows/built-in-workflows/build-workflow.md)
- derived-from [build-leaves-step](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step.md)
- derived-from [build-model-ladder](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-model-ladder.md)
- derived-from [build-leaf-order](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-order.md)
- derived-from [build-leaf-commit](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-commit.md)
- derived-from [build-verifier-checkout](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-verifier-checkout.md)
- derived-from [build-evidence-commit](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-evidence-commit.md)
- derived-from [build-verifier-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-verifier-findings.md)
- derived-from [build-retries](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-retries.md)
- derived-from [build-leaf-fails](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-fails.md)
- derived-from build-panes (retired)
- derived-from [build-check](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-check.md)
- derived-from [build-guide](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide.md)
- derived-from [build-workflow-file](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-workflow-file.md)
- derived-from [answer-ab54d544-1](../decision/answer-ab54d544-1.md)
- derived-from [answer-ab54d544-2](../decision/answer-ab54d544-2.md)
- derived-from [answer-ab54d544-3](../decision/answer-ab54d544-3.md)
- derived-from [answer-ab54d544-4](../decision/answer-ab54d544-4.md)
- derived-from [answer-ab54d544-5](../decision/answer-ab54d544-5.md)
- derived-from [answer-ab54d544-6](../decision/answer-ab54d544-6.md)
- derived-from [answer-ab54d544-7](../decision/answer-ab54d544-7.md)
- derived-from [answer-ab54d544-8](../decision/answer-ab54d544-8.md)
- derived-from [answer-ab54d544-9](../decision/answer-ab54d544-9.md)
- derived-from [answer-ab54d544-10](../decision/answer-ab54d544-10.md)
- derived-from [answer-ab54d544-11](../decision/answer-ab54d544-11.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
