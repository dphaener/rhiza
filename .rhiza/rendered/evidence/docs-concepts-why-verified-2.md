# docs-concepts-why-verified-2

evidence · `01M6ZEXD5HAQW3VW1TRB58YWKN` · actor [claude-code-7e1121e2](../actor/claude-code-7e1121e2.md)

## result

PASS

## statement

Verified docs/concepts/why.md at bb07b7b3 (attempt 2, fixing qa f2). go test -trimpath -count=1 ./internal/doclinks/... passes, TestRepoDocsLinksResolve included (docs-links-resolve). docs-commands-exist has no check built yet, so I checked every rhiza command the page names (apply, accept/reject --reason/--rule, add, link, question answer, change new, set, source add --raw/--raw-sha256/--raw-session, why --view, slice) against --help on a scratch build: all exist. Ran the fixed Questions example as written in a fresh scratch repo with RHIZA_* unset: add capability export-csv --record, add question csv-delimiter --record, link csv-delimiter blocks export-csv --record, change new, set export-csv, question answer csv-delimiter. The answer is accepted and writes answer-csv-delimiter plus its link. Also checked the doc's negative claim: a question added and linked unlanded inside the same change set is refused ('not an open question blocking this change set'). The diff only touches the Questions section, the new prose reads plainly with no AI tells, and the page still covers decisions and owner quotes, questions, sources and scrubbing, evidence, and why/slice. The owner-reads part of the acceptance is left to the owner.

## Links

### Out

- verifies [docs-concepts-why](../capability/rhiza/documentation/docs-concepts-why.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
