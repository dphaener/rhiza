# reopen-format-verified

evidence · `01M5WT1Q51QNZFYAZXN1CGBGX6` · actor [claude-code-b5610794](../actor/claude-code-b5610794.md)

## result

PASS under make check at commit afc1020 on branch cs/reopening-a-change-set, run by an independent verifier

## statement

Corpus and golden tests read, fold and render the unchanged format-1 and format-2 cases and a new format-3 reopen case, and entry tests show reopen needs format 3 with each version stamped as before.

## tests

go test: TestCorpusFolds (internal/changeset), TestCorpus and TestCorpusCorrupt (internal/store), TestCorpus (internal/model), TestGolden (internal/render), each over format-1/*, format-2/carve and format-3/reopen; TestTypes, TestVersions (1, 2, 3), TestReopenNeedsFormat3 (internal/entry); TestVersionsWithoutCarve, TestReopenNeedsFormat3 (internal/changeset); git diff --stat main -- testdata/corpus/format-1 testdata/corpus/format-2 is empty

## Links

### Out

- verifies [reopen-format](../capability/rhiza/change-sets/reopen/reopen-format.md)

## History

- `01M5WPTRX1DDQ7BAGDF7T29QPW` Reopening a change set: rhiza reopen lets a change set landed only on its branch take more writes and land again, so no session deletes a log entry to fix its work.
