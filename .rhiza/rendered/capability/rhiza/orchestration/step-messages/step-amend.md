# step-amend

capability · `01M5W98BZJ7N8NSBPBMMNWV70F` · parent [step-messages](../step-messages.md)

## acceptance

A testscript with a fake Rex amends a step: the brief file changes, the journal holds one amendment entry with the sender's actor, and one pointer line naming the file is sent.

## statement

rhiza can amend a running step: it edits the step's brief or plan file, records the amendment and its sender's actor in the run's journal, and sends the step a one-line pointer to the changed file.

## Links

### Out

- decided-by [amend-is-brief-edit-plus-pointer](../../../../decision/amend-is-brief-edit-plus-pointer.md)

### In

- verifies [step-amend-script](../../../../evidence/step-amend-script.md)

## History

- `01M5SWTQ11MR7CP75QWPYFBQFQ` Genesis amends a running step and sends it a message.
