# review-workflow-plan

source · `01M5Z0BYDW9Q25H145596AEK0M`

## quotes

```json
[
  {
    "by": "agent (accepted by the owner)",
    "text": "New `checks` step kind (Recommended); Code first, then five agents (Recommended); Opus for all (Recommended)",
    "where": "topic 1: how the checks run"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Extend run findings (Recommended); Critic per serious check (Recommended); Reported, not fixed (Recommended)",
    "where": "topic 2: findings and the critic"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Fixer, then verifier per leaf (Recommended); One change-set fixer (Recommended); Questions; run fails (Recommended)",
    "where": "topic 3: fixes and the cycle cap"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Approve or request changes (Recommended); Rendered here, posted by ship (Recommended); QA Evidence + proposed rules (Recommended)",
    "where": "topic 4: sign-off, summary, model"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Throwaway worktree each (Recommended); Merge-base with main (Recommended)",
    "where": "topic 5: where the checks work"
  }
]
```

## raw_sha256

80e8a0b60336e91cd1c2dac2bea37fdeb3469f12f058b864af51ceeb9a6f4c9c

## session

a5ffc966-b55d-4564-aae8-01a17dee91f5

## summary

The plan session for change set 01M5WPTZ1V0JVVW8FPX0QYMRWJ (the review workflow) on 2026-10-27, within milestone 3's review-workflow-shape decision. The owner answered five topics by picking the agent's recommended option every time, typing nothing. Decided: a new `checks` step kind runs the agent checks as fresh sessions in parallel, each in a pane, and re-runs only the failed ones. The built-in review starts with a code check (`rhiza check`, `rhiza run build-check`, and the repo's own command through a workflow override), then runs five agent checks (compliance, quality, hands-on QA with the built binary, scope, model). Checks and critics run on opus. `rhiza run findings` and `rhiza run verdict` gain a `check` field and severities (blocker, major, minor, nit), each finding optionally naming a leaf. A critic runs only for a check with a blocker or major finding, can't add findings, and a check fails only on a confirmed blocker or major. Minors and nits are reported, not fixed. A confirmed finding goes to a fixer per leaf, whose fix commit carries the leaf's `Leaf:` trailer, and then the leaf's verifier runs again. Findings naming no leaf go to one change-set fixer with a `Review:` trailer, and scope findings needing outside work become questions or discovered-from drafts. After 2 fix cycles the surviving findings become questions and the run fails before sign-off; a resume after an answer gets one fresh cycle, and the arbiter answers in autonomous runs. The owner signs off with `rhiza run approve` or `rhiza run request-changes <file>` (owner findings, already confirmed, one extra cycle outside the cap), and autonomous runs skip the wait. `rhiza run summary` renders the verdicts as markdown; the PR body is change set 8's. The QA check records Evidence on the leaves it exercised, the model check proposes missing constraints in a discovered-from draft, and dropped findings stay in the journal. Each check and critic works in a throwaway detached worktree at the branch head and reviews from the merge-base with the default branch. Left open: what a failing code check does.

## title

Review workflow planning, with the owner

## Links

### In

- derived-from [built-in-review](../capability/rhiza/workflows/built-in-workflows/built-in-review.md)
- derived-from [answer-a5ffc966-1](../decision/answer-a5ffc966-1.md)
- derived-from [answer-a5ffc966-2](../decision/answer-a5ffc966-2.md)
- derived-from [answer-a5ffc966-3](../decision/answer-a5ffc966-3.md)
- derived-from [answer-a5ffc966-4](../decision/answer-a5ffc966-4.md)
- derived-from [answer-a5ffc966-5](../decision/answer-a5ffc966-5.md)
- derived-from [answer-a5ffc966-6](../decision/answer-a5ffc966-6.md)
- derived-from [answer-a5ffc966-7](../decision/answer-a5ffc966-7.md)
- derived-from [answer-a5ffc966-8](../decision/answer-a5ffc966-8.md)
- derived-from [answer-a5ffc966-9](../decision/answer-a5ffc966-9.md)
- derived-from [answer-a5ffc966-10](../decision/answer-a5ffc966-10.md)
- derived-from [answer-a5ffc966-11](../decision/answer-a5ffc966-11.md)
- derived-from [answer-a5ffc966-12](../decision/answer-a5ffc966-12.md)
- derived-from [answer-a5ffc966-13](../decision/answer-a5ffc966-13.md)
- derived-from [answer-a5ffc966-14](../decision/answer-a5ffc966-14.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
