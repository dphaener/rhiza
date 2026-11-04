# why-paged-verified

evidence · `01M6J07KYTYVRSGCVMQS042551` · actor [claude-code-8f8c6422](../actor/claude-code-8f8c6422.md)

## result

PASS

## statement

Verified at c4ffa768. Ran go test ./cmd/rhiza -run 'TestScripts/(why-paged|describe-golden)$' -count=1: ok. Read why-paged.txtar: cache gets 3 decisions, 2 sources and 2 change sets; walk.sh walks every --json page at --limit 1, 2 and 3, and the concatenated (list, ID) entries cmp equal to --all in order. Every page carries node and view, only the last has no next_cursor, and pages of 2 split decisions/sources and sources/change sets. why --limit 3 --json cmp equals api call why. --limit -1, --limit 1001 and a bad cursor exit 2. describe-golden has limit, cursor, all and next_cursor, and the golden test passes. Read the lineage.go diff: page() counts one sequence across the three lists; the cursor is base64 of 'list offset' and strictly decoded. By hand with a scratch build in a scratch repo (3 decisions, 2 sources, 2 change sets): cursor walks at --limit 1..8 equal --all in order; api call why and why --json print the same bytes; a forged out-of-range cursor exits 2. With RHIZA_* unset, fmt-check, self-check and go vet pass. go test ./... passes except the read-budget timing tests (modelapp TestGeneratedModel, api ready budget) at load average ~20 from other runs; go test ./internal/api -skip Budget: ok. Attempt 1's finding f1 (no in-order compare with --all) is fixed.

## Links

### Out

- verifies [why-paged](../capability/rhiza/views/paged-lists/why-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
