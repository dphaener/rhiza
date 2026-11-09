# pr-404-fix-rhiza-change-show

evidence · `01M6Z1N4TZKVKRFXZ0R8H5P9BB` · actor [dphaener](../actor/dphaener.md)

## head

50617ee7367bca08ab1eb6c45e91b30cb44040b3

## pr

https://github.com/dphaener/rhiza/pull/404

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #404 (https://github.com/dphaener/rhiza/pull/404) passed CI at 50617ee

## status

verified

## Links

### Out

- verifies [shown-text-canonical](../capability/rhiza/model/text-references/shown-text-canonical.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
