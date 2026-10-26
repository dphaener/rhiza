# leaf-depends-on-acyclic

capability · `01M5WR3ZBTGE5H3YSXP61HGY9C` · parent [leaf-dependencies](../leaf-dependencies.md)

## acceptance

A testscript links a depends-on b and b depends-on c, then c depends-on a by rhiza link and by rhiza apply: both exit 3 with the loop named and the log entry count unchanged. A Go check test folds a loop made by two concurrent change sets and gets one fail finding on rule depends-on naming the three capabilities.

## statement

depends-on between capabilities never forms a loop, as between change sets: rhiza link and rhiza apply refuse a depends-on that would close one, and rhiza check fails on each loop already in the log.

## Links

### Out

- decided-by [build-workflow-shape](../../../../decision/build-workflow-shape.md)
- decided-by [capability-loops-on-depends-on](../../../../decision/capability-loops-on-depends-on.md)

### In

- depends-on [leaf-order](leaf-order.md)
- verifies [leaf-depends-on-acyclic-verified](../../../../evidence/leaf-depends-on-acyclic-verified.md)
- verifies [pr-36-leaf-dependencies-and](../../../../evidence/pr-36-leaf-dependencies-and.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
