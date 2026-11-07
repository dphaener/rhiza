# m29-nebari-ships-after-import

decision · `01M6RYBZEMWR88WVV3WB1PPJJ5`

## statement

vm-second-project-ships waits for Nebari's real import (milestone 27's import procedure run on Nebari) and then ships a real Nebari change set through the root daemon; it does not use a Nebari given only rhiza init.

## wording

Wait for the import: the check waits until Nebari has been imported into Rhiza, then ships a real change set; the agent's option, picked by the owner

## Links

### Out

- answers [m29-nebari-model-first](../question/m29-nebari-model-first.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
