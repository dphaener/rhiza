# milestone-3-conversation

source · `01M5WPTJ0WSM5QRW0M2WS7J4E6`

## quotes

```json
[
  {
    "by": "owner",
    "text": "Perfect.",
    "where": "milestone-3 conversation, topic 1: what milestone 3 is for"
  },
  {
    "by": "owner",
    "text": "1. Yes the model is the plan and the state",
    "where": "milestone-3 conversation, topic 2: the shape of the workflows"
  },
  {
    "by": "owner",
    "text": "2. Yes those are the four workflows",
    "where": "milestone-3 conversation, topic 2"
  },
  {
    "by": "owner",
    "text": "3. Yes those are perfect, but I do want to have an autonomous mode that skips those gates.",
    "where": "milestone-3 conversation, topic 2"
  },
  {
    "by": "owner",
    "text": "1. Correct, just as afk does, as described.",
    "where": "milestone-3 conversation, topic 3"
  },
  {
    "by": "owner",
    "text": "2. Looks correct!",
    "where": "milestone-3 conversation, topic 3: the plan workflow"
  },
  {
    "by": "owner",
    "text": "everything else looks good.",
    "where": "milestone-3 conversation, topic 4: the build workflow"
  },
  {
    "by": "owner",
    "text": "2. Yes that is correct.",
    "where": "milestone-3 conversation, topic 4"
  },
  {
    "by": "owner",
    "text": "It doesn't necessarily have to be a sub-agent. With rex it can just be a regular old agent since agents can talk to each other now.",
    "where": "milestone-3 conversation, topic 4, point 1"
  },
  {
    "by": "owner",
    "text": "Right now I'm not seeing the benefit of having an LLM conductor, unless possibly we're running in autonomous mode and there need to be decisions made on the larger scale level.",
    "where": "milestone-3 conversation, topic 4, point 1"
  },
  {
    "by": "owner",
    "text": "1. I like c, but starting with a and folding in b once the engine has proven itself",
    "where": "milestone-3 conversation, topic 4, git for parallel leaves"
  },
  {
    "by": "owner",
    "text": "2. Freely for sure, with messages recorded so we have a log to fall back on and inspect.",
    "where": "milestone-3 conversation, topic 4, how sessions talk"
  },
  {
    "by": "owner",
    "text": "1. arbiter propsal looks good",
    "where": "milestone-3 conversation, topic 5"
  },
  {
    "by": "owner",
    "text": "2. review looks good",
    "where": "milestone-3 conversation, topic 5: the review workflow"
  },
  {
    "by": "owner",
    "text": "everything looks good",
    "where": "milestone-3 conversation, topic 6: what milestone 3 covers"
  },
  {
    "by": "owner",
    "text": "2. Looks good",
    "where": "milestone-3 conversation, topic 7: milestone 3 as change sets"
  },
  {
    "by": "owner",
    "text": "Yes accept the test change and merge.",
    "where": "milestone-3 conversation, on the runner read-speed change set"
  },
  {
    "by": "owner",
    "text": "Let's do it",
    "where": "milestone-3 conversation, start"
  },
  {
    "by": "owner",
    "text": "Wait, shouldn't Rhiza be doint this?",
    "where": "milestone-3 conversation, topic 1"
  },
  {
    "by": "owner",
    "text": "Ahhhh, we don't have our default workflow yet. Gotcha. We kind of need that though, and I want it to be built out correctly. Which milestone is that in?",
    "where": "milestone-3 conversation, topic 1"
  }
]
```

## raw_sha256

1d8a2673120241d40b849933412dcecd0fc8ee71d1066fd3d13cd78b4cf4a1e1

## summary

The owner and the orchestrator (Claude Code in Rex) planned milestone 3 inside Rhiza on 2026-10-26, after every planned change set of milestone 2 had landed. Decided: milestone 3 is Rhiza's own workflows (plan, build, review, with ship switched over), so milestone 4 is planned, built and reviewed by Rhiza; the model is the plan and the state; two owner gates (plan acceptance, ship sign-off) and an autonomous mode that skips them but never bypasses permissions; the plan workflow; the build workflow with Rhiza's run engine as the conductor and full agent sessions as builders and verifiers, only a verifier writing a leaf's Evidence, one worktree first, sessions talking freely with every message logged; an arbiter called only at decision points; the review workflow; the scope, done points and eight change sets, with reopening a change set first as milestone 2's proof. Most answers were the owner accepting the orchestrator's proposals, whose wording is the orchestrator's; the owner's own words shaped the autonomous mode, full sessions, the engine as conductor, the worktree order and logged messages.

## title

Milestone 3 planning, with the owner

## Links

### In

- derived-from [milestone-3](../intent/anyone-builds/milestone-3.md)
- derived-from [m3-rhizas-own-workflows](../decision/m3-rhizas-own-workflows.md)
- derived-from [model-is-plan-and-state](../decision/model-is-plan-and-state.md)
- derived-from [four-workflows](../decision/four-workflows.md)
- derived-from [two-owner-gates](../decision/two-owner-gates.md)
- derived-from [autonomous-mode](../decision/autonomous-mode.md)
- derived-from [plan-workflow-shape](../decision/plan-workflow-shape.md)
- derived-from [build-workflow-shape](../decision/build-workflow-shape.md)
- derived-from [verifier-writes-evidence](../decision/verifier-writes-evidence.md)
- derived-from [full-agent-sessions](../decision/full-agent-sessions.md)
- derived-from [engine-is-conductor](../decision/engine-is-conductor.md)
- derived-from [one-worktree-first](../decision/one-worktree-first.md)
- derived-from [sessions-talk-freely](../decision/sessions-talk-freely.md)
- derived-from [arbiter](../decision/arbiter.md)
- derived-from [review-workflow-shape](../decision/review-workflow-shape.md)
- derived-from [m3-scope](../decision/m3-scope.md)
- derived-from [m3-change-sets](../decision/m3-change-sets.md)
- derived-from [budget-timed-alone](../decision/budget-timed-alone.md)

## History

- `01M5WPTHR3EBHDNNH1CFX7A0T0` Milestone 3 planning record: the conversation, the milestone-3 intent and its done points, and its decisions
