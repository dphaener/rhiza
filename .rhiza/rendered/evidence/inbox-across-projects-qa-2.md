# inbox-across-projects-qa-2

evidence · `01M6SYP6XH17CND694T20J311S` · actor [claude-code-fea95113](../actor/claude-code-fea95113.md)

## result

PASS

## statement

QA cycle 2 at 23ada88d with a scratch build (go build -trimpath ./cmd/rhiza), RHIZA_* env cleared. On the real root ~/code, rhiza inbox --json lists rhiza and nebari failed runs each with its project field (text shows project: lines); --project nebari lists only nebari items. In a scratch root (HOME and git config scratch) with rhiza.toml registering alpha holding a question run: at the plain root the question is listed tagged alpha; after git init + rhiza init of the root itself, rhiza inbox lists the root own empty inbox, --all lists alpha question, --project alpha lists it. The cycle-1 major (own-model root merged every project) is fixed. inbox-projects testscript passes.

## Links

### Out

- verifies [inbox-across-projects](../capability/rhiza/orchestration/inbox-across-projects.md)

## History

- `01M6RN3VGZ7D1JK4FRPNX21WHX` Cross-project CLI: at the root or with --project, rhiza daemon status, inbox, roadmap and health cover every project or one
