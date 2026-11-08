# question-cf4d0c69-1

question · `01M6XHQ6HBXZXWDW1CZWN4ZA2P`

## session

cf4d0c69-229d-4917-b1e9-01a1bb16ffc1

## statement

Change set 'Fix Ship-change stops on a' (PR 383) can never land: it updated a node (rename, new statement), then retired it in a later entry. `rhiza retire` bases the retire on the node's landed version (Node.Version skips the change set's own ops), so at land model.checkBase sees the retire as stale against the change set's own earlier update and refuses (exit 3). rhiza check only warns (rule in-flight), so CI passed and the watchdog caught it at land. How should this change set fix it?

## Links

### In

- answers [retire-supersedes-own-pending-edits](../decision/retire-supersedes-own-pending-edits.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
