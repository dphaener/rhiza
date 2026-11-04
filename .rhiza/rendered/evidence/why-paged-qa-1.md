# why-paged-qa-1

evidence · `01M6J18FSSEMEZWZMM5WCXMWP9` · actor [claude-code-17e5e990](../actor/claude-code-17e5e990.md)

## result

PASS

## statement

At b63f7c18 with a scratch build on this repo, why 01M5SWRSPVMK689DD23KH7KP0Z (43 decisions, 5 sources, 1 change set) paged with --limit 1, 7, 42, 44 and 48 (splitting at and across list boundaries): every page carried the --all node and view, the concatenated decisions, sources and change sets equal --all in order, and the last page had no next_cursor. --limit 0 equals the default; --all ignores --limit; --limit -1, --limit 1001 and --cursor zz exited 2. why --json and api call why printed the same bytes. why-paged.txtar passes.

## Links

### Out

- verifies [why-paged](../capability/rhiza/views/paged-lists/why-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
