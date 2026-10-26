# fresh-read-by-folder-stats

decision · `01M5V5BSMKWHFT3HCYNXC8YFFP`

## refines

t5-d6-freshness-by-names

## statement

A fresh index open first compares a fingerprint of the names in log/ and each change set folder's stat (modification time, mode, change time, inode and device) with the one stored at the last catch-up, and skips listing the folders when they match; any mismatch, or a folder modified within the last two seconds, falls back to deciding freshness by names.

## wording

agent of change set 01M5SWTJ, in autonomous mode

## Links

### Out

- relates-to [t5-d6-freshness-by-names](t5-d6-freshness-by-names.md) — refines it: the stat fingerprint is a pre-check that skips the listing when nothing changed; names still decide freshness on any mismatch

### In

- decided-by [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWTJY0AZP9TR3V1QJ54HMY` planGather (behind roadmap and ready) and the index's fresh read get faster, so CI has headroom under 25 ms again.
