# capability-loops-on-depends-on

decision · `01M5X0MBH7KJYYTMDA69EF2KMM`

## statement

A depends-on loop among capabilities is refused at write time naming the shortest path back, and reported by rhiza check on the depends-on rule, once per loop, naming its capabilities.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### In

- decided-by [leaf-depends-on-acyclic](../capability/rhiza/model/leaf-dependencies/leaf-depends-on-acyclic.md)

## History

- `01M5WPTV0NX6YECXCPMHSEJAWY` Leaf dependencies and Evidence's actor: leaves link depends-on each other, Evidence names the actor that wrote it, and rhiza check warns when that actor also wrote the leaf's code.
