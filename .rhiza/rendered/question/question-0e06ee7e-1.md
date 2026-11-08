# question-0e06ee7e-1

question · `01M6VYDNN7EMSMMZ8Z6K5Q41AZ`

## session

0e06ee7e-7cb8-432f-97f8-01a1b7e43938

## statement

Change set 01M6VY7 (requests carry an opaque --ref) serves no intent, and `rhiza health` reports 0 open milestones. Which intent should it serve so the daemon will build it? / Where should a request's --ref live on the change set its plan run opens? / What length limit should Rhiza enforce on a --ref? / Decision workflows-are-the-rules (PR #353) is not landed, so this plan cannot link to it. What should the plan do?

## Links

### In

- answers [request-ref-shape](../decision/request-ref-shape.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
