# t4-d14-conflict-repair

decision · `01M5SWS032T0YCH3RKNK0Y1KMZ`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Rewrite the aspects (Recommended): Cleared once a later op written from the current version writes every aspect the conflict overlapped (e.g. D sets `statement` again). Precise: an unrelated edit doesn't hide it.

## statement

(Settled in review.) A landed conflict is repaired, and no longer reported, once later ops written from the target's current version have rewritten every aspect the conflicting op wrote. An unrelated edit does not clear it. A conflict whose op is a retirement, or whose target is retired, cannot be repaired this way (a retired target takes no more writes) and stays reported

## wording

agent of milestone-1 task 4 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-4](../source/plan-task-4.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [conflicts](../capability/genesis/change-sets/conflicts.md)

## History

- `01M5SWRZJR0EXNHS4GAQDAJPHZ` Task 4's plan decisions (PR #4, merge 380f2f2), accepted by the owner
