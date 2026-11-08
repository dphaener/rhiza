# pr-357-requests-carry-an-opaque

evidence · `01M6W39RBJ9XXEJCXRYZ0WMZ5T` · actor [dphaener](../actor/dphaener.md)

## head

7cdefa5947dd9ec6dd44dca6b0b26b040a961cd5

## pr

https://github.com/dphaener/rhiza/pull/357

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #357 (https://github.com/dphaener/rhiza/pull/357) passed CI at 7cdefa5

## status

verified

## Links

### Out

- verifies [request-ref](../capability/rhiza/orchestration/daemon/plan-requests/request-ref.md)
- verifies [change-set-ref](../capability/rhiza/change-sets/change-set-ref.md)
- verifies [plan-run-carries-ref](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-carries-ref.md)
- verifies [event-payload-ref](../capability/rhiza/hooks/model-event-hooks/event-payload-ref.md)
- verifies [ship-pr-body-ref](../capability/ship-through-rhiza/ship-pr-body-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
