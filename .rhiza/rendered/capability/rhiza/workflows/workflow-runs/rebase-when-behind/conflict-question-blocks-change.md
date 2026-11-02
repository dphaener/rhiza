# conflict-question-blocks-change

capability · `01M6D133EKJDM0ACMC374VY9PK` · parent [rebase-when-behind](../rebase-when-behind.md)

## acceptance

A ship rebase testscript with a code conflict: the raised question has a blocks link to the change set; rhiza inbox --json lists it as a blocker of that change set; rhiza question answer <q> <words> --change <cs> writes the answer; raising it again on the same conflict writes nothing

## statement

A rebase-conflict question raised by ship rebase blocks the change set being shipped, so the inbox lists it as a blocker and rhiza question answer accepts it with --change

## Links

### Out

- decided-by [conflict-question-blocks](../../../../../decision/conflict-question-blocks.md)

### In

- verifies [conflict-question-blocks-change-verified](../../../../../evidence/conflict-question-blocks-change-verified.md)
- verifies [conflict-question-blocks-change-qa-1](../../../../../evidence/conflict-question-blocks-change-qa-1.md)
- verifies [conflict-question-blocks-change-verified-2](../../../../../evidence/conflict-question-blocks-change-verified-2.md)
- verifies [pr-189-conflict-question-blocks](../../../../../evidence/pr-189-conflict-question-blocks.md)

## History

- `01M6D12NNWQZHEWFBX60JS04AG` A run's rebase-conflict question blocks its change set, so the inbox and question answer see it
