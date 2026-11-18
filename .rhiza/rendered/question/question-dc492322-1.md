# question-dc492322-1

question · `01M7QDYT3TTWCZ7MXC0S1WWB2F`

## session

dc492322-38b6-41d5-bece-01a1eedd8716

## statement

Change set 01M7M2FB (rhiza run spawn can return): should returning once the child starts be opt-in on rhiza run spawn (a new --no-wait flag; waiting stays the default) or the new default (spawn always returns at start)? / Change set 01M7M2FB: who sends the conductor its wake-up message (rhiza session message) when a --no-wait child ends? / Change set 01M7M2FB: the intent says a child's 'blocked question' also wakes the conductor. In a conductor step the arbiter already answers a child's question. What should wake the conductor?

## Links

### In

- answers [spawn-no-wait-wake](../decision/spawn-no-wait-wake.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
