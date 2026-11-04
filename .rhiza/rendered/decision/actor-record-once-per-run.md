# actor-record-once-per-run

decision · `01M6M2XXMPAABRE5CHJTTA4F0R`

## reasoning

m4-runs-hold-up says runs finish 'with no swept-up actor records', and the earlier actor-record fix 01M602TH served milestone-4, so this follow-up belongs to the same intent. The question names the cause: recordAgentActor (internal/api/session.go:668) writes only in the run worktree, and the verifier's checkout (the session record's Dir) lacks it, so the agent hand-adds a second record. session.go already has the mechanism to fix that in worktreeWorkspace: entry files are immutable and named by ID, so a copy is the same file a later merge brings. Copying the one record into the checkout removes the cause; teaching check/land to merge identical duplicates would only hide a duplicated write and widens the model's rules, which the change set's intent offers only as a fallback. Matching on the session field makes stamping idempotent if a session is ever started twice.

## statement

Change set 01M6M2R7 serves m4-runs-hold-up (milestone-4). Fix it at the source: session start writes the agent's actor record once, in the run worktree as now, and when the session has its own checkout (a verifier's dir) it copies that same record change set's log entries into the checkout, so the actor resolves there and no agent ever hand-adds it. Because the copies are the same immutable files with the same IDs, committing them from both places merges into one record, not a second change set with a duplicate handle. recordAgentActor also reuses an existing actor whose session field matches instead of adding a new one, so a repeat stamp is a no-op. rhiza check and land do not learn to merge duplicate actor records.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-7346bd67-1](../question/question-7346bd67-1.md)
- derived-from [actor-once-per-run-plan](../source/actor-once-per-run-plan.md)

### In

- decided-by [actor-record-in-checkout](../capability/rhiza/runtime/actor-stamping/actor-record-in-checkout.md)
- decided-by [actor-stamp-reuses-session](../capability/rhiza/runtime/actor-stamping/actor-stamp-reuses-session.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
