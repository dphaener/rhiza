# request-ref-qa-1

evidence · `01M6W1234PYBBJF17T8391AEWF` · actor [claude-code-2590959f](../actor/claude-code-2590959f.md)

## result

PASS

## statement

QA at 66e50695 with a built binary: request 'Fix login' --kind bug --ref '  https://github.com/o/r/issues/7  ' stored it trimmed; a ref with an OSC sequence, BEL and newline came back verbatim in request ls --json and Go-quoted in human request ls ('[ref "x\x1b]0;pwn\a\ny"]', no raw control bytes); 2048 bytes queued; 2049 bytes exited 2 with '--ref is 2049 bytes, at most 2048' and queued nothing; --ref '' queued a request with no ref key. TestScripts/request-ref and describe-golden pass.

## Links

### Out

- verifies [request-ref](../capability/rhiza/orchestration/daemon/plan-requests/request-ref.md)

## History

- `01M6VY7JMM5NPB3QVG1ECZH1YW` Requests carry an opaque external reference: rhiza request '<text>' --ref <string> (for example an issue URL) stores the ref on the request and on the change set its plan run plans; it is carried in the rhiza.eventCommand JSON payload for every move of that change set (proposed, accepted, landed and so on), shown by request ls and change show, and listed in the ship PR body. Rhiza never interprets, fetches or validates it beyond a length limit and treats it as untrusted text wherever it shows it, so an outside tool can link a request to its source (e.g. close a GitHub issue on landed) with no tracker-specific code in Rhiza. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
