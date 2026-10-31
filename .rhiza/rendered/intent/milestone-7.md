# milestone-7

intent · `01M686S9JF5GP5KMDQG2SVHCTZ`

## milestone

7

## statement

The owner drives Rhiza from Rex and sees what it does: unattended failures say why, everything waiting on the owner is in the inbox, new work starts from the roadmap app, the live app shows runs in progress, umbrella areas stay out of the way, and updating Rhiza is one command

## Links

### Out

- decided-by [m7-picked-by-orchestrator](../decision/m7-picked-by-orchestrator.md)

### In

- serves change set `01M686RRCJSQ57M82FENWT2HNM` (When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why)
- serves change set `01M686RS47EK7F4H2WY15TNAQ5` (rhiza daemon status names each accepted change set the daemon does not start and why)
- serves change set `01M686RSVTRCWEV9GPT2V16M9Y` (The inbox lists open questions that block proposed and accepted change sets and the inbox app answers them in place, through one question.answer operation the roadmap app shares)
- serves change set `01M686RTJM0TBSN6C1FXES93WN` (Umbrella change sets, which serve an intent marked umbrella, are grouped apart in the roadmap, are never ready, and collapse to one row in the roadmap app)
- serves change set `01M686RVAD1EHGERBBDJN3S45G` (The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent)
- serves change set `01M686RW277XSQ20MYCE580GG8` (A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run)
- serves change set `01M686RWSEKXXBCNC7RZMT07HB` (make update rebuilds rhiza in place, rewrites the Rex script and reloads Rex in one command)
- serves change set `01M696Q0P1JNY00XQ1R312XVRJ` (The model app's read timing test waits until go test runs nothing else before it times, as TestReadBudget does, so machine load no longer fails CI or a run)
- serves change set `01M69CDDZAS7JREQN9AWKPZ7XY` (make test gives each package 20 minutes instead of go test's default 10, so cmd/rhiza's growing script suite stops timing out on CI's two-core runner)

## History

- `01M686RDMZR9FBZXN7DZWRWQW9` Plan milestone 7: drive Rhiza from Rex and see what it does, picked by the orchestrator under the owner's standing instruction
