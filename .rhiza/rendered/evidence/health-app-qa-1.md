# health-app-qa-1

evidence · `01M6DG4CHCH2A54BA8F0BVH9VS` · actor [claude-code-6e06158d](../actor/claude-code-6e06158d.md)

## result

PASS

## statement

At ee31e06b, QA ran the real app. In a scratch repo, rhiza tui app health --once printed the same bytes as rhiza health --json. With free disk under the minimum it printed state bad and exited 0. rhiza rex install into a temp REX_CONFIG wrote rhiza_app_health running tui app health with ["health"] = "t", and no key letter is used twice. QA launched rhiza tui app health in a real Rex block with a dirty checkout. It showed 'Health: warn' with each part marked [ok]/[warn], the cursor on daemon and the footer 'j/k move · enter open app · r reload · q quit'. After cleaning the checkout, r reloaded it and main became [ok]. j moved the cursor to inbox, and enter opened the inbox app ('Nothing waits on you', inbox footer). The healthapp unit tests and tui-health-app/rex-install scripts pass.

## Links

### Out

- verifies [health-app](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-app.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
