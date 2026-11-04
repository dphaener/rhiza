# anchors-plan-session

source · `01M6HK8FVXEQSREWWBVVX2RHYF`

## quotes

```json
[]
```

## raw_sha256

95fd5bf0f47da14975a4b6ca6214f8c843523f2f9a749b19fede25f52fe3dba7

## session

298a341f-f050-4419-bc88-01a1a32ec55a

## summary

Autonomous plan run for change set 01M6HJXAJ6GYA2RTAGVK6JR1N9; no owner took part. The planner proposed four choices (which intent to serve, what the fingerprint hashes, how anchors are written and re-fingerprinted, and what check resolves against and how severe it is), and the arbiter settled them in Decision anchor-fingerprint-shape. The change set serves model-tells-truth. An anchor carries path, an optional Go symbol (top-level name or Type.Method, .go files only) and a fingerprint, 'sha256:<hex>' of the declaration's source bytes or else of the whole file, read from the working tree. Writing an anchor computes the fingerprint and refuses a path or symbol that doesn't resolve. A separate operation re-fingerprints an anchor after review. rhiza check resolves every live anchor against the working tree, naming the anchor and its capability in each finding: a missing path or a missing symbol is a broken anchor and fails, as the landed Decision ci-fails-on requires; a changed fingerprint alone only warns. The anchors and fingerprints leaves are carved out of the Area integrity umbrella 01M5SWT0K9TV1R and split into four buildable leaves.

## title

Code anchors with fingerprints, planned autonomously

## Links

### In

- derived-from [anchor-fingerprint-shape](../decision/anchor-fingerprint-shape.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
