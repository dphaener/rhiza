# request-ref-shape

decision · `01M6VYJMHYEJ0FXDWPSR5ZJ27W`

## reasoning

Read the question, the change set's intent and slice (no nodes yet), rhiza health (0 open milestones; daemon builds only change sets under a milestone), milestone-25's statement and its 37 serving change sets, and the arbiter decision on the same question for 01M6Q58 (ci-speed-shape), which served milestone-25 for an owner-requested change; the roadmap shows no other open milestone fits. For the ref's home I read internal/api/request.go (requests live in clone-local .git/rhiza/requests.json), internal/changeset/body.go (Open already carries an optional Title with a format-version bump, the established pattern for an optional open field) and the eventCommand setting in internal/api/daemon.go (it runs per lifecycle move read from the log), so only the log reaches every consumer the intent names. 2048 bytes covers real issue and ticket URLs with margin while bounding log and PR-body size; the intent forbids validating anything beyond length, so safety comes from escaping on display. gh pr view 353 shows it OPEN and unmerged, so a link would be refused; the feature does not need that decision's code to be correct, so blocking on it would only stall the daemon.

## statement

(1) Intent: change set 01M6VY7 serves milestone-25 ('The backlog is empty: every change set planned from the owner's queued requests and the backlog is built and shipped by the daemon on its own'), as the other owner-requested change sets planned from requests do (e.g. 01M6Q58, 01M6PQTK, 01M6RXNQ serve it); link 01M6VY7 serves milestone-25 and no other milestone. Serving it reopens milestone-25 on the roadmap, which is what lets the daemon build this. (2) Where the ref lives: in the log, as an optional `ref` field on the change set's open entry body, beside `title` and built the same way (internal/changeset Open gains Ref; encoded only when set; an open body with a ref is stamped with a new entry format version so an older binary refuses it rather than dropping it silently). It is not kept only in the clone-local .git/rhiza/requests.json, because rhiza.eventCommand fires from lifecycle moves read out of the log on whichever clone's daemon sees them, and change show must show it on any clone. The request queue also stores the ref on the request (RequestRequest/RequestInfo gain Ref, shown by request ls). `rhiza change new` gains a --ref flag; the plan run that serves a request opens its change set with the request's ref (the daemon hands the request's ref to the run, and the run passes it to change new). The event payload, change show and the ship PR body read it from the folded change set. (3) Length limit: at most 2048 bytes after trimming surrounding whitespace (an empty ref after trimming means none); a longer one is a usage error naming the limit. Nothing else is validated: no URL parsing, no charset or newline rules. Because it is untrusted, every place that shows it escapes it rather than validating it: human output prints it Go-quoted (strconv.Quote-style) so control characters and terminal escapes cannot act; JSON carries it as a plain string; the ship PR body shows it inside an inline code span with backticks and newlines escaped so it cannot inject markdown or mentions. (4) PR #353 not landed: the plan does not link to decision workflows-are-the-rules and does not wait for it. This change set's behaviour does not depend on it (Rhiza enforces only the length limit and stores the ref as opaque text, which already fits that principle), so it names the principle in prose only, with no @ reference, no depends-on and no question parked on it. A later change set may add the link once #353 lands.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-0e06ee7e-1](../question/question-0e06ee7e-1.md)
- derived-from [ref-plan-interview](../source/ref-plan-interview.md)

### In

- decided-by [request-ref](../capability/rhiza/orchestration/daemon/plan-requests/request-ref.md)
- decided-by [change-set-ref](../capability/rhiza/change-sets/change-set-ref.md)
- decided-by [plan-run-carries-ref](../capability/rhiza/orchestration/daemon/plan-requests/plan-run-carries-ref.md)
- decided-by [event-payload-ref](../capability/rhiza/hooks/model-event-hooks/event-payload-ref.md)
- decided-by [ship-pr-body-ref](../capability/ship-through-rhiza/ship-pr-body-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
