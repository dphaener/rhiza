# m29-second-project-waits-for-owner

decision · `01M6SPER99F36FNKYMN2008FD7`

## reasoning

The leaf's acceptance is operational: a Nebari run and a rhiza run live under the root daemon and the Nebari PR merging. All three attempts failed for the same reason, not a build defect: Nebari has only draft change sets (notes on m29-which-nebari-change-set list all twelve), rhiza.autonomy.plan = owner, and milestone and change set planning needs the owner. Shipping one also merges a PR in a second product, an outward act the owner has not chosen. The infrastructure the leaf exercises is already in place (vm-root-installed verified: root daemon up, both projects registered), so the leaf's statement and acceptance are still right and should not be weakened; what is missing is an owner input that m29-which-nebari-change-set already asks for. Rebuilding it again would only add more notes, so the right ruling is to wait for that answer rather than change the leaf or have an agent pick a Nebari change set (e.g. Disabled button contrast) in the owner's place.

## statement

vm-second-project-ships is not rebuilt by builders: it stays as written and stays blocked on m29-which-nebari-change-set until the owner picks and plans a Nebari change set. No agent picks, plans or accepts a Nebari change set to satisfy it; the leaf is retried only after that change set is accepted, and the rest of the change set proceeds without it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6snq9-1](../question/leaf-failed-01m6snq9-1.md)

## History

- `01M6RN474B5EXMX8NYHE7F62RV` This VM on one root: rhiza and nebari registered under one root run by one root daemon service, replacing rhiza-rhiza.service
