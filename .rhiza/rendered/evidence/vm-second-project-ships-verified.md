# vm-second-project-ships-verified

evidence · `01M6TGDDF04WDCP9CB32S1Y4EM` · actor [claude-code-a54ef720](../actor/claude-code-a54ef720.md)

## result

PASS

## statement

Verifier claude-code-a54ef720 on VM rhiza, 2026-11-07: (1) ~/orchestrator/vm-root-migration.log line 82, section '2026-11-07T04:32:31Z rhiza daemon status at the root (~/code)', shows one root daemon's status listing nebari run 01M6T4EHEM1TBM0CET0SCBKHHG running (change set 01M6S6R84W7GD1J26JESF0ZH18, Disabled button contrast) and, in the same output, rhiza runs 01M6T36NJ2YFVFPT8MS9YK8TDW and 01M6T3HSJVBP341H5C7V78K2T4 running. rhiza daemon status now shows that root daemon (rhiza:3484327, mode service, Max runs 5 from rhiza.toml). (2) gh pr view 72 -R dphaener/nebari: state MERGED, mergedAt 2026-11-07T07:39:35Z, merge commit 1689c77b26bf1810d5106a8c4d7938fa166b2a4b, title 'Disabled button contrast', head cs/disabled-button-contrast; ~/code/nebari has 1689c77 'Disabled button contrast (#72)'. The leaf commit 0eedbf44 only adds the builder's actor record.

## Links

### Out

- verifies [vm-second-project-ships](../capability/rhiza/orchestration/daemon/vm-second-project-ships.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
