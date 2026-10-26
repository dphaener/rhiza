# leaf-order

capability · `01M5WR3ZJQCG6MCWRGY69NWJJY` · parent [leaf-dependencies](../leaf-dependencies.md)

## acceptance

A testscript adds leaves c, b and a in that order with c depends-on b and b depends-on a; change show --json lists a, b, c with b's needs [a] and c's needs [b]; leaves with no links come out in tree order; describe-golden is regenerated for the new response field.

## statement

rhiza change show --json lists the change set's leaves (the capabilities it adds that have no live children) in build order: each leaf after every leaf it depends-on, ties in tree order, each with the handles it depends on.

## Links

### Out

- depends-on [leaf-depends-on-acyclic](leaf-depends-on-acyclic.md)
- decided-by [build-workflow-shape](../../../../decision/build-workflow-shape.md)
- decided-by [leaf-needs-every-target](../../../../decision/leaf-needs-every-target.md)

### In

- verifies [leaf-order-verified](../../../../evidence/leaf-order-verified.md)
- verifies [pr-36-leaf-dependencies-and](../../../../evidence/pr-36-leaf-dependencies-and.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
