# import-friction-recorded-verified

evidence · `01M6RPKYFCN474STBPZDKAF559` · actor [claude-code-8bd25498](../actor/claude-code-8bd25498.md)

## result

PASS

## statement

At eb79786e, internal/guide/text/import.md has a '## Record friction' section stating the rule (record every step you could not follow, or had to work around, as a question naming the step, so the import can improve; never work around one silently) and the question form `rhiza add question <handle> --field 'statement=Import step 3 (record built capabilities): <what, and how>'`. cmd/rhiza/testdata/script/import-friction-recorded.txtar checks that the form's example appears in `rhiza guide import` output, and runs the form in a fresh repo (fold-ok). Ran go test -trimpath -count=1 ./internal/guide/... (PASS: TestCaps now passes, so the attempt-1 finding about the 800-token cap is fixed) and go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(import-.*|guide)$' (PASS: import-friction-recorded, import-records-built, import-autonomous-decisions, import-retires-docs, import-unbuilt-as-drafts and guide). Reading the diff, the tightening keeps the other import steps' rules intact.

## Links

### Out

- verifies [import-friction-recorded](../capability/rhiza/guide/import-friction-recorded.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
