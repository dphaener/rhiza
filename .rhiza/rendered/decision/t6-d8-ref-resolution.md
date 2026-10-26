# t6-d8-ref-resolution

decision · `01M5SWS6Y9WXF65KNC3HBCWXQN`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: ID/prefix/handle, unique per repo (Recommended): A ref is a full ID, a unique ID prefix (4+ chars), or a handle of a live (not retired) node. Handles are unique across live nodes in the to-be view; `add` or a rename to a taken handle fails (exit 3). No match: exit 4. Several matches: exit 2, listing the candidates.

## statement

A ref is a full ID, a unique ID prefix (4+ chars), or the handle of a live node in the to-be view. Handles are unique across live nodes: `add` or a rename to a taken handle fails with exit 3. No match: exit 4. Several matches: exit 2, listing the candidates

## wording

agent of milestone-1 task 6 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-6](../source/plan-task-6.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [handles](../capability/genesis/model/handles.md)

## History

- `01M5SWS5QKWPQ3FS357EAV60D2` Task 6's plan decisions (PR #5, merge 6d5c61e), accepted by the owner
