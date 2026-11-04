# log-paged-qa-1

evidence · `01M6J177KMYD0283YE3PM2ZR8E` · actor [claude-code-17e5e990](../actor/claude-code-17e5e990.md)

## result

PASS

## statement

At b63f7c18 I built rhiza to a scratch path and ran it by hand on this repo's .rhiza. log 01M5TYGK7EFHKPB9JG3A2MH2G2 (three change sets) --limit 2 printed two change sets and 'next: rhiza log 01M5TYGK7EFHKPB9JG3A2MH2G2 --limit 2 --cursor <c>'; that cursor's page held the third with no next_cursor and the same target and view; the two pages concatenated equal --all. --limit 0 gave the default page; --limit -1, --limit 1001, --cursor zz and --cursor !!! exited 2 with 'next: rhiza describe log'. log --json and api call log printed the same bytes. go test of log-paged and describe-golden scripts: ok.

## Links

### Out

- verifies [log-paged](../capability/rhiza/views/paged-lists/log-paged.md)

## History

- `01M6H3A65RDP761HKE82STZSQZ` log, roadmap and why page their output with a cursor, as ls already does (replaces backlog change set 01M5SWTHZSPBNQ, which has no leaves; link the new change set to it as relates-to and withdraw the old one in the plan)
