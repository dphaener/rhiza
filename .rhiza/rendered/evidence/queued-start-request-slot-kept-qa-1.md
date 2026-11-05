# queued-start-request-slot-kept-qa-1

evidence · `01M6MS6BS7NH93VFNAASP280EZ` · actor [claude-code-93557052](../actor/claude-code-93557052.md)

## result

PASS

## statement

Built rhiza from the branch head and drove it by hand in a scratch git repo (fake rex on PATH, detach recorded via RHIZA_TEST_DETACH_LOG): three change sets queued with rhiza run start --workflow plan. With no request, daemon --once --max-runs 3 started all three. With request 'Add exports' queued, daemon --once --max-runs 3 started two queued runs plus the request's plan run, leaving one start in starts.json; the next poll started it. With --max-runs 1 and a queued request, the queued start took the slot and the request stayed queued. Also ran TestScripts/daemon-queued-start-slot and daemon-request-slot: ok.

## Links

### Out

- verifies [queued-start-request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/queued-start-request-slot-kept.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
