# dispatch-kind-field-request-kind

decision · `01M6VZW00RN5VJZSHTTRN1Q8A9`

## reasoning

The finding is right: changeset.Kind is intention|record (internal/changeset/body.go:66-77), change new --kind takes intention|record, and change show --json already prints "kind" from c.Open.Kind, so reusing the name would collide, and a milestone request's change set must be kind record and milestone at once. The two are orthogonal (intention vs record is how a change set lands; feature/bug/milestone is what kind of request it came from), so they need two fields; overloading or replacing kind would break every existing log entry and the --kind flag. request_kind is the clearest name because the values are exactly the request kinds already fixed in internal/api/request.go (enum feature,bug,milestone, daemon.KindFeature/Bug/Milestone), it says where the value comes from, and it is the example name the plan reviewer and critic both proposed. Keeping the default feature for change sets no request opened matches @workflow-dispatch-rules-file and request.go's 'default feature', so the decision's rules are unchanged; only the field name is pinned.

## statement

The dispatch kind a change set records is a new field named request_kind (feature, bug or milestone), separate from the existing kind (intention or record), which keeps its meaning and values unchanged. The open entry carries request_kind next to the optional workflow; rhiza change show --json prints "request_kind" and "workflow" beside the existing "kind", and the text output labels it 'request kind'. A change set no request opened has request_kind feature, so a milestone request's change set reads kind record with request_kind milestone. @{01M6VYTAJBDW1YH8D7Y1XGCK76}'s statement and acceptance and @{01M6VYQWKXT7DC11ET2AAAER7G}'s wording 'records its kind' / 'kind bug' / 'kind feature' mean request_kind, and dispatch rules key on request_kind.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6vy7e-dispatch-r3-f1](../question/review-01m6vy7e-dispatch-r3-f1.md)

## History

- `01M6VY7G6HPJC3V0Q18E8EF9ZA` Workflow dispatch: every run the daemon starts goes through repo-defined selection rules (in rhiza.toml) that map a trigger (request queued, change set accepted and ready) and the request's kind (feature, bug, milestone) to a workflow, and rhiza request --workflow <name> overrides the choice for one request; the change set records the kind and any workflow its request chose so the ship pick follows it. Today daemonStartPlan always starts planWorkflow 'plan' (internal/api/run.go:763) and a ready change set always gets defaultWorkflow 'ship-change' (run.go:1257); workflow-overrides only shadow those names. Built-in default rules reproduce today's behaviour, so a repo adds e.g. a bugfix workflow in .rhiza/workflows/ and maps bug to it without changing Rhiza. Relates to the workflow-authoring change set 01M6HN24. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
