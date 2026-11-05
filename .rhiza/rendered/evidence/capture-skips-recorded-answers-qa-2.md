# capture-skips-recorded-answers-qa-2

evidence · `01M6MV6DAPN47NKFFFGT49TH30` · actor [claude-code-11d93090](../actor/claude-code-11d93090.md)

## result

PASS

## statement

QA cycle 2, hands-on with a binary built from the branch head, in a scratch repo with a fake rex and a three-question transcript (Which colour?/Which store?/Which port?). The change set held a Decision with question 'Which port?' verbatim and one with 'which colour?' (case differs). rhiza session done printed 'Proposed decision answer-<sid8>-1', '-2' and 'Skipped answer-<sid8>-3: already recorded'; only answer-1 and answer-2 were written. In a second change set, with 'Which store?' recorded and a retired 'Which colour?' Decision, session done --json returned decisions [answer-<sid8>-1, answer-<sid8>-3] and skipped [answer-<sid8>-2]; answer-1 carries question 'Which colour?' and the owner-pick wording. A repeat session done returned already:true and proposed nothing. go test TestScripts/session-interview-recorded and describe-golden pass.

## Links

### Out

- verifies [capture-skips-recorded-answers](../capability/rhiza/lineage/interview-capture/capture-skips-recorded-answers.md)

## History

- `01M6J6FDD9GXDW4AENMS6YNY63` Fix Plan interviews duplicate decisions: when the planner already wrote Decisions for the owner's answers, session done still auto-creates answer-<session>-N Decisions from the same AskUserQuestion picks (seen in plan run 01M6FY99J90TC3RGT2ZCDXFNR5: seven duplicates), with default owner wording even when the answers came from a proxy. session done should skip picks a Decision in the change set already answers (or link to it), and the planner guide should say which side records them. Also each rhiza apply write takes ~6 s on the VM (31 writes > 3 min); see the daemon CPU bug
