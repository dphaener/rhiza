# review-code-check-in-rex-qa-1

evidence · `01M66DH35AB8MH4JXYPZTT9A9G` · actor [claude-code-31db92cd](../actor/claude-code-31db92cd.md)

## result

PASS

## statement

Built ./cmd/rhiza at 5dd508d and ran the builtin review inside a real Rex block on a scratch repo whose leaf beta had no independent Evidence. Real Rex got, read back with rex events --json, the root record working 'step code', then error 'run failed: step "code" ended with outcome "fail"'. The run exited 1, and rhiza run status listed 'step code: fail (leaf beta: no Evidence verifies it)'.

## Links

### Out

- verifies [review-code-check-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-code-check-in-rex.md)

## History

- `01M602TGXM184GVWJG4NE98P6Y` Rex program status shows a review run's live fixer, verifier and change-set fixer sessions, its open questions and its code check, as rhiza run status already does
