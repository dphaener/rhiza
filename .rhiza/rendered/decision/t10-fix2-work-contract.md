# t10-fix2-work-contract

decision · `01M5SWSVKBBK3P6K3AYTEM6AY3`

## authorized_by

owner (the orchestrator decided under the owner's authorization)

## settled_by

decided by the orchestrator under the owner's authorization, in its sign-off review of task 10; applied in PR #8 before merge

## statement

genesis.md core loop: the work contract. Agents work on an accepted change set (`ready` requires it). The contract: read with `slice`; write back only revisions within the change set's intent; anything outside it becomes a new change set linked `discovered-from`, or a question; finish with `genesis land` as the last commit. This replaces "Then `genesis propose`" (propose stays for a draft the agent opens itself); change-sets.md agrees.

## wording

the orchestrator's review at sign-off, applied by task 10's agent

## Links

### Out

- derived-from [plan-task-10](../source/plan-task-10.md)
- derived-from [founding-conversation](../source/founding-conversation.md)

### In

- decided-by [skills](../capability/rhiza/guide/skills.md)

## History

- `01M5SWSRAVRYJ1AZH6RT1GF122` Task 10's plan decisions (PR #8, merge 3e17d7b), accepted by the owner
