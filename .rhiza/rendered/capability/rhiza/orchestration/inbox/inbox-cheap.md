# inbox-cheap

capability · `01M6KVA76HM3KXBBGQJWP3NPGV` · parent [inbox](../inbox.md)

## acceptance

Go tests: TestInboxReadsChangeSetsOnce (internal/api) gives an inbox over 100+ change sets, an overlap and a failed run per change set with zero log reads when an index is set, and without one the fold reads the log exactly once and gives the same states, titles and overlaps as the index; TestStoreCachesUnchangedJournals (internal/run) parses an unchanged journal once over three lists and again after an append; TestUsageRecordsReadsOnlyAppended (internal/adapter/claudecode) reads only the appended bytes of a grown transcript, with the same records as a whole read, and reads a rewritten or truncated one whole. Measured on the rhiza repo (4-core VM, 2069 change sets, 7849 entries, 40 runs): rhiza inbox before 12.9 s wall (17.6 s under load), after 0.22-0.44 s with byte-identical output; an idle, paused daemon on a copy of the repo before 3.4-4.5 CPU-s per minute, after 1.2; a daemon fold after a model change before 5.5 s, after 1.8 s.

## statement

rhiza inbox and the daemon's polls stay cheap on a large repo: one inbox call reads the change sets once (from the index when there is one, else one log read) for every state, title and overlap lookup, and an idle daemon poll re-parses neither unchanged run journals nor the already-read part of a growing agent transcript

## Links

### Out

- decided-by [inbox-cheap-diagnosis](../../../../decision/inbox-cheap-diagnosis.md)

## History

- `01M6KV89S5MKWG2DHA762F8WTZ` Inbox and daemon notify are cheap: overlaps and inbox items are computed once per model change, not per call or poll
