# usage-window-tokens-verified

evidence · `01M6HYYPBCGZ8F7J2R7BFK27P4` · actor claude-code-a463ece8 (retired)

## result

PASS

## statement

At b67b44b7 ran go test ./internal/adapter/claudecode -run TestUsageRecords -count=1 -v: TestUsageRecords and TestUsageRecordsNoConfig PASS. Read the test: a fake CLAUDE_CONFIG_DIR with transcripts in two project dirs (-a, -b); message m1 repeated on two content-block lines (counted once, 60 = input+output+cache-creation, its 5000 cache reads excluded); a helper transcript under s1/subagents with sidechain lines and 9000 cache reads (counted 300); records at 09:59:59 and 12:00 outside [10:00,11:00) dropped; a user line and a non-JSON line skipped; m3 written out of order. reflect.DeepEqual checks the exact four records {10:05,7},{10:10,60},{10:20,300},{10:40,3} in time order. A file outside the config dir, plus a symlink to it inside projects, would add a 999-token record at 10:30 if read, so the exact match shows nothing outside is read (WalkDir with d.Type().IsRegular() skips symlinks). Read the diff: an adapter.UsageReader interface with UsageRecord{At,Tokens}; claudecode.UsageRecords walks <config>/projects ($CLAUDE_CONFIG_DIR, else $HOME/.claude), dedupes by message id at the earliest timestamp, and filters to [from,to). make fmt-check vet self-check exit 0. Full go test ./...: every package passed except internal/api (read-timing budget test) and internal/slice and cmd/rhiza TestScripts, which hit the 10m package timeout at load average about 28 from parallel agents; internal/slice then passed with -timeout 40m. Nothing outside internal/adapter references the new code, so those failures are from machine load, not this leaf.

## Links

### Out

- verifies [usage-window-tokens](../capability/rhiza/orchestration/usage-governor/usage-window-tokens.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
