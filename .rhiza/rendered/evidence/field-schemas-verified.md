# field-schemas-verified

evidence · `01M69ZTY3DF3P22MMXAY8J4J0H` · actor [claude-code-0270f381](../actor/claude-code-0270f381.md)

## result

PASS

## statement

At 553e1f53 I ran go test ./internal/op -count=1 -v: TestEveryFieldInRepoLogIsKnown walks every .rhiza/log entry's node.create fields and node.update set keys (shapes checked against the real log) and finds each known for its kind; TestBogusFieldUnknown shows bogus unknown and result known on evidence; TestKnownFieldsSortedAndCoverEveryKind passes. Reading internal/op/fields.go, every field the acceptance names (intent statement/milestone/umbrella; capability statement/acceptance/built_in; evidence statement/result/tests/test/status/head/pr/command/observed; decision, source, question, constraint, actor fields) is listed. make check passes on a clean detached worktree of the commit.

## Links

### Out

- verifies [field-schemas](../capability/rhiza/model/node-kinds/field-schemas.md)

## History

- `01M5SWTH32C73BQET8E0708F9F` Each core node kind has a field schema, so its fields are known and checked.
