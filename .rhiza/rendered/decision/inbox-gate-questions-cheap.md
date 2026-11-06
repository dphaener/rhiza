# inbox-gate-questions-cheap

decision · `01M6RXVTQJ0T8RYE8JS3P2MQC9`

## reasoning

The question has four parts. Milestone-26's Inbox app change set (01M6MYN9) names 'answer a question or gate' as the app's job, so this fix finishes that promise; milestone-28 (web UI) only gains speed from the same shared read. internal/api/inbox.go:1067 sets c.IndexPath = "" and then calls w.model(cs.ID) inside the per-handle loop over plan, which matches the profile (ReadLog 13.5 s, 80% in inboxGateQuestions). internal/api/run_model_cache.go shows the run worktree already has its own index at index.DefaultPath(root), kept current by the driver, so reusing it is the smallest fix that fits milestone-18/20's 'reads stay fast' and adds no new cache. The app already has openGate/answerGate, but it falls back to show() (a non-scrollable InfoPanel with no action) whenever Questions comes back empty. That happens when inboxGateQuestions swallows an error and returns nil, which explains the owner's 'couldn't scroll it and couldn't really take any action'. So the fix also guarantees that questions are returned from the stop detail and that the views scroll. A test counting decoded entries and folds is a sturdier proof than wall-clock time alone on a busy CI runner.

## statement

(1) Intent: change set 01M6RXK2 serves milestone-26 (01M6MZ09, five small Rex apps), whose Inbox app change set already promises 'answer a question or gate'; the web UI's speed is fixed by the same shared read, so no second intent is linked. (2) Cheap read: inboxGateQuestions stops clearing IndexPath. It opens the run worktree's own index (index.DefaultPath of the worktree root, the index the run driver already keeps caught up, as stepRunner.catchUpModel does), so no read decodes the whole disk log. It also inverts the loop: each planned change set's model is folded at most once per call and every open handle is resolved against it, never question × change set. If the worktree's model can't be read, it still returns one InboxGateQuestion per handle from the stop detail (statement empty), so the item never loses its questions or its action. (3) Proof: an internal/api budget test over the large generated log with an index, as the existing read-budget tests have, puts a plan run at its gate with several open questions spread over several planned change sets. It asserts that inbox returns every question with its statement, decodes no full log (the catch-up decodes 0 entries when fresh), folds each planned change set at most once, and answers within the read budget, well under 1 s. The daemon poll's gate-question fill (daemon.go) goes through the same path and is covered by the same test. (4) Inbox app: picking a gate ('open') item always opens the gate answer form: the full question text, with advisor recommendations, scrolls (arrow and page keys and the mouse wheel within the modal's height), with one TextArea per open question. Submitting runs question.answer <handle> <text> --run <run> for each question (a Decision each), then run.resume, then closes, as answerGate does today. A blank answer asks again. A failed answer reopens the form on the questions not yet answered with the error, as answerGate does today, so the owner can retry; only a failed resume shows a scrollable detail with the error. A question with an empty statement shows 'This question's text could not be read.' with its position among the gate's questions ('question 2 of 4'), never its handle (owner-surfaces-no-ids). The plain detail panel of every other item kind scrolls too.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-fa508624-1](../question/question-fa508624-1.md)
- derived-from [plan-inbox-gate-fast](../source/plan-inbox-gate-fast.md)

### In

- decided-by [inbox-gate-read-fast](../capability/rhiza/inbox-gate-read-fast.md)
- decided-by [inbox-gate-questions-kept](../capability/rhiza/inbox-gate-questions-kept.md)
- decided-by [inbox-app-gate-scroll](../capability/rhiza/inbox-app-gate-scroll.md)

## History

- `01M6RXK2SB0M12FXVEGXER022R` Fix The inbox takes 30-40 s whenever a run is stopped at a gate with open questions, and the Inbox app can't act on that item. Measured on the VM with milestone 29's plan stopped at its gate (4 open questions): rhiza inbox 31.7 s, the web UI's Overview 42.5 s and /fragments/inbox 26.9 s (it was ~0.2 s); the owner: "The app took maybe 10 seconds to load." A CPU profile of rhiza inbox puts 80% in inboxGateQuestions (internal/api/inbox.go): it opens the run's worktree with IndexPath cleared, so every read decodes the full disk log (11,400 entries, ReadLog 13.5 s), and it folds w.model(cs.ID) for every question handle times every planned change set. Make it cheap (index or cache, one fold per change set) so the inbox, the Inbox app and the web UI stay well under 1 s. Also, in the Inbox app this 'open' gate item only shows its long detail, which can't be scrolled, and offers no action; owner: "I opened the nebari inbox message and it was pretty long and I couldn't scroll it and couldn't really take any action on it." The app must let the owner scroll the detail and answer each open gate question (a Decision per question, then resume the run), as it does for single questions.
