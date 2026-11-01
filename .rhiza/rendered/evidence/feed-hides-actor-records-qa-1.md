# feed-hides-actor-records-qa-1

evidence · `01M6C0RXVS45XGC6G7Z4PHR3B4` · actor [claude-code-6544b34b](../actor/claude-code-6544b34b.md)

## result

PASS

## statement

QA at the branch head with a binary built by go build -o $TMPDIR/rhiza-qa ./cmd/rhiza. On this repo's own model, rhiza feed --limit 1000 listed 138 change sets. feed --all listed 1000 and its cursor page listed 286 more, 1286 in all. The 1148 hidden ones were all actor records ('Record agent actor', 'Record actor', 'add actor ...'), and the shown ones kept their --all order. Paging the default feed with --limit 50 gave full pages of 50, 50 and 38, which adds up to 138. The --all next: line carries --all. In a scratch git repo (RHIZA_* unset) I made capability and intent records. feed hid the 'Record actor' change set and --all showed it. A mixed apply that adds an actor and a capability was shown by default. An actor-only update (set bot-two --field name=... --record) was hidden by default and shown with --all. tui app changes --once listed only the non-actor change sets. A second look listed none. .git/rhiza/feed.json's seen list held only the shown IDs, not the hidden actor ones.

## Links

### Out

- verifies [feed-hides-actor-records](../capability/rhiza/views/change-feed/feed-command/feed-hides-actor-records.md)

## History

- `01M6BH7XPH07NSDTBKE2RGX949` rhiza feed and the changes app leave out actor-only change sets unless --all is given
