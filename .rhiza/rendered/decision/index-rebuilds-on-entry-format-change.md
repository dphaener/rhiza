# index-rebuilds-on-entry-format-change

decision · `01M73T7X227V5V26TS8KX1ADAD`

## reasoning

internal/index/index.go already has one recovery mechanism: a rebuild that empties every table and catches up from nothing, triggered by a vanished path or a schema-version mismatch, and the package doc calls the index disposable. Keying a rebuild on the decoded FormatVersion reuses that path with one stored value and one comparison, so it cannot drift from FoldLog the way a partial re-decode of only newer-format problem rows could (re-decoding would also have to refold every touched change set and recompute both views, i.e. most of a rebuild anyway). Format bumps are rare (one per release like #424), so a full rebuild's cost is paid seldom. Rebuilding on any difference, not only a rise, also covers a downgrade then upgrade, where a lower binary stored problems for entries it could not read. It meets the intent's acceptance test: write a v(N+1) entry, open with FormatVersion N (problem stored), then N+1 (rebuilt, folds clean). The 14G of orphan dbs is a real but separate issue (cache hygiene across worktrees, not wrong folds); the intent lists it only as a side note, and widening this fix would delay unsticking landed change sets, so it is left for its own change set.

## statement

The persistent index stores the entry.FormatVersion it decoded with (a meta row beside the listing fingerprint) and, on Open, treats a stored version different from the binary's as stale and rebuilds from nothing, the same path a vanished file or another schema version takes; it does not selectively re-decode problem rows. Pruning index databases of worktrees that no longer exist is out of this change set's scope and is not done here; it belongs to a separate change set.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-76c5b90e-1](../question/question-76c5b90e-1.md)
- derived-from [stale-index-format-plan](../source/stale-index-format-plan.md)

### In

- decided-by [index-rebuilds-on-entry-format-change-cap](../capability/rhiza/views/fast-reads/index-rebuilds-on-entry-format-change-cap.md)

## History

- `01M73T319T5VT2Q37MFXTDP8T6` Fix The persistent index keeps an entry's newer-format problem after the binary learns that format, so change sets stay unlanded forever. After #424 bumped the entry format to 5 and the daemon self-updated, the three milestone-31 docs change sets (01M71G4BKFHP1E172RR6H82WFT, 01M71G4G92NEJ33C04DSKC4RMK, 01M71G4MMC4VGT41WVM4A33ZWM) merged (PRs 426-428, runs 01M73CQ3T, 01M73CQ4D, 01M73CQ4C succeeded) but folded as accepted and stayed in rhiza ready/roadmap as ready. rhiza change show 01M71G4G92 --json problems: 'entry out of order: after names 01M731FQPE9PY5EYQTA9H0C4CF, which is not an entry of this change set' (the v5 pin entry, written while the installed binary read v4), then every later entry, the land included, 'left out of the fold'. With a fresh cache (XDG_CACHE_HOME=empty dir) the same binary folds it landed with no problems. Cause: internal/index stores each listed entry file with its problem and only re-reads a file when its listing changes; a rebuild happens only on a vanished path or a schema change (internal/index/index.go doc, ~l.40), so an entry stored as newer-format is never re-decoded after an upgrade to a binary whose entry.FormatVersion covers it. Expected: the index records the FormatVersion it decoded with and re-decodes (or rebuilds) entries stored with a newer-format problem when the binary's FormatVersion rises; a test writes a v(N+1) entry, opens the index with FormatVersion N, then N+1, and the change set folds clean. Unstuck by moving ~/.cache/rhiza/index/<sha256 of /home/rhiza/code/rhiza>.db aside so it rebuilt. Side note: ~/.cache/rhiza/index holds 1580 dbs, 14G, mostly from test and run worktrees that no longer exist; nothing prunes them.
