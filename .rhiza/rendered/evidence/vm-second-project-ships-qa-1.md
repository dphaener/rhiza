# vm-second-project-ships-qa-1

evidence · `01M6TGNN8XHD6CNJ054AY50S3V` · actor [claude-code-1237c200](../actor/claude-code-1237c200.md)

## result

PASS

## statement

QA claude-code-1237c200 on VM rhiza, 2026-11-07: `rhiza daemon status` in ~/code (root daemon rhiza:3484327, mode service, Max runs 5 from rhiza.toml) lists both Project nebari and Project rhiza under the one daemon, with nebari's last start a run of change set 01M6S6R84W (Disabled button contrast) at 04:32:11Z and four rhiza runs live now. ~/orchestrator/vm-root-migration.log's 2026-11-07T04:32:31Z root daemon status shows nebari run 01M6T4EHEM1TBM0CET0SCBKHHG live beside rhiza runs 01M6T36NJ2YF and 01M6T3HSJVBP. `rhiza run status 01M6T4EHEM1TBM0CET0SCBKHHG` in ~/code/nebari: succeeded, ship-change workflow, autonomous. `gh pr view 72 -R dphaener/nebari`: MERGED at 2026-11-07T07:39:35Z, merge commit 1689c77, head cs/disabled-button-contrast; ~/code/nebari main is at 1689c77 'Disabled button contrast (#72)'. The concurrent live snapshot comes from the log, since the nebari run has finished.

## Links

### Out

- verifies [vm-second-project-ships](../capability/rhiza/orchestration/daemon/vm-second-project-ships.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
