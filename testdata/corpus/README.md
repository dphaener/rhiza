# The log corpus

Sample `.genesis` logs from every released format version.

Genesis never compacts or rewrites its log, so a repository written by any
released Genesis must still load in every later one, forever. This corpus is
that promise, checked on every test run.

## Layout

One folder per format, then one folder per case:

```
format-1/
  <case>/
    .genesis/
      repo.json
      log/<change_set>/<id>.json
format-2/
  <case>/        (same layout)
```

## Rules

- Format 1 may be amended until the first tagged release. After that, cases
  are only added, never edited or removed. The same holds for every later
  format from its first release.
- The corpus holds only valid logs, valid down to the typed body of every
  entry and every operation inside them. Every case folds with zero
  problems. Invalid inputs belong in table tests.
- The files are fixtures, written once. Do not regenerate them.

## Tests

`internal/store/corpus_test.go` opens every case, reads its log and expects
zero problems, checks every `.json` file is canonical, and checks that the
format-1 cases use every entry type at least once.

Format 2 adds the optional `carve` field on an ops entry. Only carve entries
are written as v2; every other entry in a format-2 case is still v1, as in a
real repository.

`internal/changeset/corpus_test.go` folds every case with
`changeset.FoldLog`, expects zero problems, and checks each change set's
final state (landed, rejected, withdrawn, superseded and so on). A new case
fails that test until it is given an expectation there.
