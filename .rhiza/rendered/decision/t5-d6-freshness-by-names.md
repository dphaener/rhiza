# t5-d6-freshness-by-names

decision · `01M5SWS2B0K2ZF8WHR0HKSC2EE`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: File names; gone → rebuild (Recommended): List log/ folders and compare entry file names with those indexed (entries are immutable, so names suffice); new files catch up incrementally, any vanished file (branch switch) triggers a full rebuild; hand-edited content is check's job.

## statement

Freshness by entry file names; new names catch up incrementally; any vanished name → full rebuild; content edits are `check`'s job.

## wording

agent of milestone-1 task 5 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-5](../source/plan-task-5.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- relates-to [fresh-read-by-folder-stats](fresh-read-by-folder-stats.md) — refines it: the stat fingerprint is a pre-check that skips the listing when nothing changed; names still decide freshness on any mismatch
- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWS22WXZMST7XXR9ZKY4EE` Task 5's plan decisions (PR #7, merge ed8d292), accepted by the owner
