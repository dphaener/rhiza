# plan-task-2

source · `01M5SWRGNF2CMZJE8MXG1JQT40`

## quotes

```json
[]
```

## raw_sha256

24abaf3007be2d1126e1c79098a2b4b8a1a7a9bf5481d1316314547c9532fe1b

## summary

The shipit work plan for task 2: the log's foundations as library packages (internal/id, canon, entry, store) plus a format-1 corpus, all test-first, with no new CLI commands. Decided (D1-D19, agent wording accepted by the owner): library only, with repo.json, Init and Open in the library and genesis init left to task 6; IDs are ULIDs in canonical upper-case Crockford base32, from oklog/ulid/v2; TDD with table tests, rapid properties and a fuzzed entry decoder; format 1 allows only integer JSON numbers; strings escape only what JSON requires, raw UTF-8, no HTML escaping, no normalisation, invalid UTF-8 rejected; reading is strict, so bytes must be exactly canonical and duplicate keys and unknown envelope fields are rejected; body is any JSON object, typed bodies being task 3's; writes are atomic and never overwrite (temp file, fsync, hard link, fsync the folder); the reader collects every per-file problem, keeps good entries and returns a stable order; the corpus may be amended until the first tagged release and is add-only after; the envelope rules (v, required ULIDs, UTC second-precision at, a sorted unique after list, nine entry types, an object body); repo.json's canonical shape, with Init refusing an existing .genesis/; no check between a ULID's time and at; the reader reports every stray file or folder under log/ and skips only dot-prefixed temp files; the corpus lives under testdata/corpus/format-1 and loads with zero problems. Out of scope: CLI commands, folding, handles, fractional keys, the index, projection, check. Nothing left open.

## title

Shipit plan, milestone 1 task 2: log foundation

## Links

### In

- derived-from [node-ids](../capability/rhiza/model/node-ids.md)
- derived-from [operation-log](../capability/rhiza/model/operation-log.md)
- derived-from [format-compat](../capability/rhiza/model/format-compat.md)
- derived-from [t2-d1-library-only](../decision/t2-d1-library-only.md)
- derived-from [t2-d2-ulid-ids](../decision/t2-d2-ulid-ids.md)
- derived-from [t2-d3-oklog-ulid](../decision/t2-d3-oklog-ulid.md)
- derived-from [t2-d4-tdd-rapid-fuzz](../decision/t2-d4-tdd-rapid-fuzz.md)
- derived-from [t2-d5-integer-numbers](../decision/t2-d5-integer-numbers.md)
- derived-from [t2-d6-minimal-escaping](../decision/t2-d6-minimal-escaping.md)
- derived-from [t2-d7-strict-reading](../decision/t2-d7-strict-reading.md)
- derived-from [t2-d8-opaque-body](../decision/t2-d8-opaque-body.md)
- derived-from [t2-d9-atomic-write](../decision/t2-d9-atomic-write.md)
- derived-from [t2-d10-collect-problems](../decision/t2-d10-collect-problems.md)
- derived-from [t2-d11-corpus-amendable](../decision/t2-d11-corpus-amendable.md)
- derived-from [t2-d12-task-2-out-list](../decision/t2-d12-task-2-out-list.md)
- derived-from [t2-d13-envelope-rules](../decision/t2-d13-envelope-rules.md)
- derived-from [t2-d14-four-packages](../decision/t2-d14-four-packages.md)
- derived-from [t2-d15-repo-json](../decision/t2-d15-repo-json.md)
- derived-from [t2-d16-no-time-check](../decision/t2-d16-no-time-check.md)
- derived-from [t2-d17-stray-files](../decision/t2-d17-stray-files.md)
- derived-from [t2-d18-corpus-location](../decision/t2-d18-corpus-location.md)
- derived-from [t2-d19-task-2-acceptance](../decision/t2-d19-task-2-acceptance.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
