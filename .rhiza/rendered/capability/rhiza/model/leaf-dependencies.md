# leaf-dependencies

capability · `01M5WR3Z4YXYHDRPGQB1VKJ6ZG` · parent [model](../model.md)

## acceptance

A testscript writes a change set whose three leaves link depends-on in a chain, sees change show list them in that order, and has a fourth link closing the loop refused with nothing written.

## statement

A change set's leaves link depends-on each other to say which must be built first: the links never form a loop, and the leaves read back in an order that respects them.

## Links

### Out

- decided-by [build-workflow-shape](../../../decision/build-workflow-shape.md)

### In

- verifies [pr-36-leaf-dependencies-and](../../../evidence/pr-36-leaf-dependencies-and.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
