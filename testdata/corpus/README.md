# The log corpus

Sample `.rhiza` logs from every released format version. The cases were
written when the project was called Genesis and moved here with `git mv` at the
rename, so the entries inside them still say genesis; their bytes are unchanged.

Rhiza never compacts or rewrites its log, so a repository written by any
released Genesis or Rhiza must still load in every later one, forever. This
corpus is that promise, checked on every test run.

## Layout

One folder per format, then one folder per case:

```
format-1/
  <case>/
    .rhiza/
      repo.json
      log/<change_set>/<id>.json
format-2/
  <case>/        (same layout)
format-3/
  <case>/        (same layout)
format-4/
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
zero problems, checks every `.json` file is canonical, and checks that every
entry type is used by some case of the earliest format that has it: reopen
in format 3, every other type in format 1.

Format 2 adds the optional `carve` field on an ops entry. Only carve entries
are written as v2; every other entry in a format-2 case is still v1, as in a
real repository.

Format 3 adds the reopen entry type, which moves a landed change set back
to accepted. Only reopen entries are written as v3; the format-3 case lands,
reopens and lands again, and every other entry in it is v1.

Format 4 adds the optional `ref` field on an open entry: an opaque external
reference, such as an issue URL. Only an open with a ref is written as v4;
the format-4 case opens one change set with a ref and lands it, and every
other entry in it is v1.

`internal/changeset/corpus_test.go` folds every case with
`changeset.FoldLog`, expects zero problems, and checks each change set's
final state (landed, rejected, withdrawn, superseded and so on). A new case
fails that test until it is given an expectation there.
