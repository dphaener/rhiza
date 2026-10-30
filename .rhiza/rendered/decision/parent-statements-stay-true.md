# parent-statements-stay-true

decision · `01M651AV8PXAPH0B3H28VGFET5`

## reasoning

The question is whether the decision and the parent statements contradict. run-status-tokens says run status shows each session's journaled token use; session-journal says the journal holds token use and status shows it. Neither says status shows only journaled use, and both stay true for sessions that have a journaled reading. The decision adds read-time counts only where no reading is journaled, and it says they are not journaled, so the journal statement stays exact. Amending the parents would put a narrowing or widening into a change set whose intent is only that status never says 'no token use recorded' for a session that used tokens. The arbiter may only write this one decision, so no amendment is made here; if the owner wants the parent statements to name the fallback, that is a separate change set.

## statement

The parent statements of run-status-tokens and session-journal stay as written and true. The transcript-read counts of printout-reads-transcript-when-unjournaled are shown only for a session with no journaled reading, so run status shows journaled use for every session that has one, and the counts read at print time are an extra line for a session without one. Neither parent is amended in this change set, and no leaf may say the journal holds those read-time counts.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m65027-r3-f1](../question/review-01m65027-r3-f1.md)

## History

- `01M61GW8RKSBM4BT1B8B3FR5JH` A run's final printout and run status show each session's token use, never 'no token use recorded' for sessions that used tokens.
