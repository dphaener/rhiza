# t10-d2-source-add-shape

decision · `01M5SWSRN9QEWG4H8BDCZR2TR9`

## accepted_by

owner

## settled_by

Agent's option, accepted by the owner: Handle + JSON file + --raw (Recommended): `genesis source add <handle> <file.json> --raw <transcript>`. The file holds {title, summary, quotes:[{text, by, where}]}. --raw is hashed (sha256) and never stored; `--raw-sha256` covers raw material that's already gone. Writes through a change set or --record, so --dry-run works.

## statement

`genesis source add <handle> <file.json> --raw <path>` (or `--raw-sha256 <hex>`). The file is `{title, summary, quotes:[{text, by, where}]}`. Raw material is hashed (sha256) and never stored. It writes through a change set or `--record`, and `--dry-run` works

## wording

agent of milestone-1 task 10 (shipit plan), accepted by the owner

## Links

### Out

- derived-from [plan-task-10](../source/plan-task-10.md)
- derived-from [task-interviews](../source/task-interviews.md)

### In

- decided-by [source-add](../capability/rhiza/lineage/source-add.md)

## History

- `01M5SWSRAVRYJ1AZH6RT1GF122` Task 10's plan decisions (PR #8, merge 3e17d7b), accepted by the owner
