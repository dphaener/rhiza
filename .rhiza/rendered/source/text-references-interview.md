# text-references-interview

source · `01M6HP1KE2PY062G1N62E1H7BC`

## quotes

```json
[]
```

## raw_sha256

1305d7191606ff8abaa09160a1ba6876bb4b27a6a74387a24183eab50223dbe6

## session

ed7e1e54-977b-4520-af35-01a1a35b0f55

## summary

Autonomous plan run for change set 01M6HNNWFR1J2AHQRNKFTPME46; no owner took part. The planner raised four choices: which intent to serve and how the leaf comes over, the stored form and the fields scanned, what a write does with an unresolvable reference, and how references are shown and checked. The arbiter settled all four in Decision text-references-shape. The change set serves model-tells-truth. Capability text-references was carved from the Area model umbrella 01M5SWSWDVB237, keeping its ID and its decided-by link to references-by-id, and the umbrella, left empty, was withdrawn. References are written as @handle or @<ID prefix>, starting at the start of the text or after whitespace, ( or [. A write stores each as @{<full ID>} in every string field except quote, quote_where, wording and a source's fields; @@ is a literal @. A write with an unresolvable reference is refused whole with not_found. Every read surface shows @<current handle>, or @{<full ID>} when the reference does not resolve. rhiza check fails a dangling reference and warns on a retired one. The plan has four leaves under text-references: text-ref-syntax (a pure internal/textref package), text-ref-write, text-ref-display and text-ref-check, the last three depending on the first.

## title

Text references plan interview (autonomous)

## Links

### In

- derived-from [text-references-shape](../decision/text-references-shape.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
