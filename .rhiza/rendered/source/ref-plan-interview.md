# ref-plan-interview

source · `01M6VYKQN0EVFE0E26R3S8630H`

## quotes

```json
[]
```

## raw_sha256

041f9daa74865932028c9f4981b51e263672f2ce45430a3d370c733b9f1aeedf

## session

0e06ee7e-7cb8-432f-97f8-01a1b7e43938

## summary

Autonomous plan run, so no owner took part. The planner read the request queue (internal/api/request.go, internal/daemon/requests.go), the change set open body (internal/changeset/body.go), the event hook (internal/api/daemon_events.go) and ship open-pr (internal/api/ship_openpr.go). It raised four choices to the arbiter in one question: which intent to serve with no open milestone, where the ref lives, the length limit, and what to do about decision workflows-are-the-rules (PR #353, not landed). The arbiter ruled in Decision request-ref-shape: serve milestone-25. Store the ref as an optional field on the change set's open entry, under a new entry format version, and also on the request. Limit it to 2048 bytes after trimming and validate nothing else. Escape it wherever it is shown: Go-quoted in human output, a plain string in JSON, an escaped inline code span in the PR body. Don't link the unlanded PR #353 decision. The planner then split the work into five leaves: request ref, change set ref, plan run carries ref, event payload ref and PR body ref.

## title

Plan interview: requests carry an opaque ref (autonomous)

## Links

### In

- derived-from [request-ref-shape](../decision/request-ref-shape.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
