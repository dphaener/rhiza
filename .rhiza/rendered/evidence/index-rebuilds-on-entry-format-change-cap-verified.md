# index-rebuilds-on-entry-format-change-cap-verified

evidence · `01M73VA6G5HDYMRAAM4C00FN68` · actor [claude-code-170995b7](../actor/claude-code-170995b7.md)

## result

PASS

## statement

At 7ea02cb9 I ran go test -trimpath -count=1 ./internal/index/... (ok) and TestEntryFormatChangeRebuilds -v (PASS). The test bumps the Propose entry's v to entry.FormatVersion+1 in a log with Propose/Accept/Land after it. On the first open it asserts exactly one newer-format log problem on that path, that the change set is stored as draft (later entries left out), and checkModel. It writes the entry back in place and asserts the open at the old version is statFresh with read 0. Then it bumps entryFormat and asserts a rebuild (read 4), no problems, checkModel (equal to changeset.FoldLog over the log) and state landed. A further open asserts read 0 and fresh. In a scratch worktree, disabling the meta entry_format comparison in open() makes the test fail with stats statFresh:true rebuilt:false read:0, so the fast stat-fingerprint path is what gets covered. The package doc in internal/index/index.go names a change of entry format version (entry.FormatVersion) among the rebuild causes. Both attempt-1 findings (f1 stat path, f2 problem/left-out assertions) are fixed.

## Links

### Out

- verifies [index-rebuilds-on-entry-format-change-cap](../capability/rhiza/views/fast-reads/index-rebuilds-on-entry-format-change-cap.md)

## History

- `01M73T319T5VT2Q37MFXTDP8T6` Fix The persistent index keeps an entry's newer-format problem after the binary learns that format, so change sets stay unlanded forever. After #424 bumped the entry format to 5 and the daemon self-updated, the three milestone-31 docs change sets (01M71G4BKFHP1E172RR6H82WFT, 01M71G4G92NEJ33C04DSKC4RMK, 01M71G4MMC4VGT41WVM4A33ZWM) merged (PRs 426-428, runs 01M73CQ3T, 01M73CQ4D, 01M73CQ4C succeeded) but folded as accepted and stayed in rhiza ready/roadmap as ready. rhiza change show 01M71G4G92 --json problems: 'entry out of order: after names 01M731FQPE9PY5EYQTA9H0C4CF, which is not an entry of this change set' (the v5 pin entry, written while the installed binary read v4), then every later entry, the land included, 'left out of the fold'. With a fresh cache (XDG_CACHE_HOME=empty dir) the same binary folds it landed with no problems. Cause: internal/index stores each listed entry file with its problem and only re-reads a file when its listing changes; a rebuild happens only on a vanished path or a schema change (internal/index/index.go doc, ~l.40), so an entry stored as newer-format is never re-decoded after an upgrade to a binary whose entry.FormatVersion covers it. Expected: the index records the FormatVersion it decoded with and re-decodes (or rebuilds) entries stored with a newer-format problem when the binary's FormatVersion rises; a test writes a v(N+1) entry, opens the index with FormatVersion N, then N+1, and the change set folds clean. Unstuck by moving ~/.cache/rhiza/index/<sha256 of /home/rhiza/code/rhiza>.db aside so it rebuilt. Side note: ~/.cache/rhiza/index holds 1580 dbs, 14G, mostly from test and run worktrees that no longer exist; nothing prunes them.
