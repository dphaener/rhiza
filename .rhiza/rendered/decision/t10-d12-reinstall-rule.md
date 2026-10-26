# t10-d12-reinstall-rule

decision · `01M5SWSTTBSS7KPMZJ44MYB3AK`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Overwrite ours, conflict on foreign (Recommended): If plugin.json names `genesis` and carries Genesis's marker, rewrite every file atomically: upgrades just work, and same bytes are reported unchanged. A `genesis/` folder that isn't ours exits 3 with next `--force`. Files Genesis doesn't own inside the folder are left alone.

## statement

Re-install: a folder that is ours (plugin.json name `genesis` + Genesis marker) is rewritten atomically, and identical bytes are reported unchanged. A foreign folder exits 3 with next `--force`. Files we don't own are left alone

## wording

agent of milestone-1 task 10 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-10](../source/plan-task-10.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [guide-install](../capability/genesis/guide/guide-install.md)

## History

- `01M5SWSRAVRYJ1AZH6RT1GF122` Task 10's plan decisions (PR #8, merge 3e17d7b), accepted by the owner
