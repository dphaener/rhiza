# rebase-conflict-text-references-resolved

decision · `01M6MZ6XVSYC4X1VXDA2PW3M70`

## statement

Resolve the text-references rebase conflicts by merging origin/main into cs/text-references, keeping both main's unclaimed check rule (#248) and the branch's text-ref rule (order: evidence-actor, unclaimed, anchor, text-ref), and regenerating .rhiza/rendered and describe-golden.txtar instead of hand-merging them; done in merge commits 19675ddb and 841f1008, scoped tests and rhiza check green

## wording

stated by the orchestrator acting as arbiter

## Links

### Out

- answers [rebase-conflict-text-references](../question/rebase-conflict-text-references.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
