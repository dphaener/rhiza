# index-rebuilds-on-entry-format-change-cap-qa-1

evidence · `01M73VNAW0GEFVVEA8JKBQGHHY` · actor [claude-code-02486e7a](../actor/claude-code-02486e7a.md)

## result

PASS

## statement

Hands-on with built binaries at 58750e43 against a clone of this repo's real log (change set 01M71G4G92, whose pin entry 01M731FQPE is format 5), each with its own XDG_CACHE_HOME. I built four binaries: this branch and its merge-base b4130875, each at entry.FormatVersion 5 and patched to 4 (standing for the pre-#424 binary). Merge-base: the FV4 binary's rhiza change show 01M71G4G92 --json gives accepted with 6 problems ('entry out of order: after names 01M731FQPE9PY5EYQTA9H0C4CF...'), and the upgraded FV5 binary over that same index still gives accepted with 6 problems, twice, so the bug reproduces. This branch: FV4 gives the same 6 problems, then FV5 over the same index gives landed with 0 problems, and again on a second open. The index db's meta entry_format row reads 5 after that and user_version 6. Downgrading to FV4 rewrites it to 4 (6 problems again), and FV5 afterwards gives landed with 0 problems. A schema-5 index left by the merge-base binary, opened by this branch's FV5 binary, is rebuilt (user_version 6, entry_format 5) and folds landed. rhiza ready no longer lists 01M71G4G92. A warm open takes about 0.1s. go test -trimpath ./internal/index/... ok.

## Links

### Out

- verifies [index-rebuilds-on-entry-format-change-cap](../capability/rhiza/views/fast-reads/index-rebuilds-on-entry-format-change-cap.md)

## History

- `01M73T319T5VT2Q37MFXTDP8T6` Fix The persistent index keeps an entry's newer-format problem after the binary learns that format, so change sets stay unlanded forever. After #424 bumped the entry format to 5 and the daemon self-updated, the three milestone-31 docs change sets (01M71G4BKFHP1E172RR6H82WFT, 01M71G4G92NEJ33C04DSKC4RMK, 01M71G4MMC4VGT41WVM4A33ZWM) merged (PRs 426-428, runs 01M73CQ3T, 01M73CQ4D, 01M73CQ4C succeeded) but folded as accepted and stayed in rhiza ready/roadmap as ready. rhiza change show 01M71G4G92 --json problems: 'entry out of order: after names 01M731FQPE9PY5EYQTA9H0C4CF, which is not an entry of this change set' (the v5 pin entry, written while the installed binary read v4), then every later entry, the land included, 'left out of the fold'. With a fresh cache (XDG_CACHE_HOME=empty dir) the same binary folds it landed with no problems. Cause: internal/index stores each listed entry file with its problem and only re-reads a file when its listing changes; a rebuild happens only on a vanished path or a schema change (internal/index/index.go doc, ~l.40), so an entry stored as newer-format is never re-decoded after an upgrade to a binary whose entry.FormatVersion covers it. Expected: the index records the FormatVersion it decoded with and re-decodes (or rebuilds) entries stored with a newer-format problem when the binary's FormatVersion rises; a test writes a v(N+1) entry, opens the index with FormatVersion N, then N+1, and the change set folds clean. Unstuck by moving ~/.cache/rhiza/index/<sha256 of /home/rhiza/code/rhiza>.db aside so it rebuilt. Side note: ~/.cache/rhiza/index holds 1580 dbs, 14G, mostly from test and run worktrees that no longer exist; nothing prunes them.
