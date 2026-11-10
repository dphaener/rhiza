# index-rebuilds-on-entry-format-change-cap

capability · `01M73T8SB2AFFX66MC9CZ2PYDC` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath ./internal/index/... : a Go test writes a log whose change set has an entry of format N+1 followed by later entries (a land included), opens the index with the entry format version seen as N and asserts the newer-format problem and the left-out entries, then opens the same index file (unchanged log, so the stat fingerprint matches) with the version seen as N+1 and asserts the change set folds landed with no problems, equal to changeset.FoldLog over the log; a second open at the same version does not rebuild; the package doc in internal/index/index.go lists an entry format change among the rebuild causes

## statement

The persistent index records the entry format version it decoded the log with, and an Open by a binary whose entry.FormatVersion differs (the fast stat-fingerprint path included) rebuilds it from nothing, so entries stored with a newer-format problem by an older binary are decoded again and their change sets fold as the log says

## Links

### Out

- decided-by [index-rebuilds-on-entry-format-change](../../../../decision/index-rebuilds-on-entry-format-change.md)
- derived-from [stale-index-format-plan](../../../../source/stale-index-format-plan.md)

### In

- verifies [index-rebuilds-on-entry-format-change-cap-verified](../../../../evidence/index-rebuilds-on-entry-format-change-cap-verified.md)
- verifies [index-rebuilds-on-entry-format-change-cap-qa-1](../../../../evidence/index-rebuilds-on-entry-format-change-cap-qa-1.md)
- verifies [pr-434-fix-the-persistent-index-keeps](../../../../evidence/pr-434-fix-the-persistent-index-keeps.md)

## History

- `01M73T319T5VT2Q37MFXTDP8T6` Fix The persistent index keeps an entry's newer-format problem after the binary learns that format, so change sets stay unlanded forever. After #424 bumped the entry format to 5 and the daemon self-updated, the three milestone-31 docs change sets (01M71G4BKFHP1E172RR6H82WFT, 01M71G4G92NEJ33C04DSKC4RMK, 01M71G4MMC4VGT41WVM4A33ZWM) merged (PRs 426-428, runs 01M73CQ3T, 01M73CQ4D, 01M73CQ4C succeeded) but folded as accepted and stayed in rhiza ready/roadmap as ready. rhiza change show 01M71G4G92 --json problems: 'entry out of order: after names 01M731FQPE9PY5EYQTA9H0C4CF, which is not an entry of this change set' (the v5 pin entry, written while the installed binary read v4), then every later entry, the land included, 'left out of the fold'. With a fresh cache (XDG_CACHE_HOME=empty dir) the same binary folds it landed with no problems. Cause: internal/index stores each listed entry file with its problem and only re-reads a file when its listing changes; a rebuild happens only on a vanished path or a schema change (internal/index/index.go doc, ~l.40), so an entry stored as newer-format is never re-decoded after an upgrade to a binary whose entry.FormatVersion covers it. Expected: the index records the FormatVersion it decoded with and re-decodes (or rebuilds) entries stored with a newer-format problem when the binary's FormatVersion rises; a test writes a v(N+1) entry, opens the index with FormatVersion N, then N+1, and the change set folds clean. Unstuck by moving ~/.cache/rhiza/index/<sha256 of /home/rhiza/code/rhiza>.db aside so it rebuilt. Side note: ~/.cache/rhiza/index holds 1580 dbs, 14G, mostly from test and run worktrees that no longer exist; nothing prunes them.
