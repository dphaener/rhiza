# review-01m6hnnv-r3-f1

question · `01M6HPT5MBMN2AP7QDHDZ74HP0`

## statement

Plan review round 3 left finding f1 (contradiction on text-ref-write): The leaf names only rhiza add, set and apply, but the engine writes free agent text through the same buildAdd: a session's question (internal/api/run.go:1939), check and leaf questions (run_checks.go:251, run_leaves.go:268, run_review.go:1018), interview questions (session_interview.go:158), rebase-conflict questions (ship_rebase.go:438), question.answer's Decision (question.go:85) and notes (note.go:82). If resolution lives in buildAdd, rule (3) refuses those writes whole when the text holds an @@word that names no node: this plan's own question question-ed7e1e54-1 reads 'what does a write do when @@handle names no node', so the engine's ask would fail with not_found and the run would error. ship_rebase.go:422 also dedupes by comparing the stored statement to the raw input, which stops matching once a reference is rewritten to @@{ID}. If resolution lives only in the CLI requests, engine-written @@handle text is stored literal and later shown as @@@handle. The plan must say which writes resolve references and what an engine-authored write does with an unresolvable one (e.g. store it as a literal @ instead of refusing), and the acceptance should cover a run question whose text holds an unresolvable @@word. The critic confirmed it: Engine writes of free text go through buildAdd: session questions (run.go:1939), check/leaf/review questions (run_checks.go:252, run_leaves.go:271, run_review.go:1020), interview questions (session_interview.go:158), rebase-conflict questions (ship_rebase.go:438), question.answer's Decision (question.go:85, whose statement is not a verbatim field) and notes (note.go:82). text-references-shape (2)/(3) apply to any node add or set and refuse the whole write on an unresolvable reference, while text-ref-write names only CLI add/set/apply; so either engine asks holding an @@word that names no node fail with not_found mid-run, or (if resolution is CLI-only) the decision's 'a write' is wrong and that text later shows as @@@word. ship_rebase.go:422 also compares the stored statement to raw input, which breaks once references are rewritten. The plan must say which writes resolve and what an engine-authored write does with an unresolvable reference, with acceptance for it. How should it be settled?

## Links

### Out

- blocks [text-ref-write](../capability/rhiza/model/text-references/text-ref-write.md)

### In

- answers [text-ref-engine-writes](../decision/text-ref-engine-writes.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
