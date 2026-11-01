# pr-159-no-agent-session-starts-runs

evidence · `01M6AZ8W1DYR5Q5H7NBC1D3QE3` · actor [darin-haener](../actor/darin-haener.md)

## head

0a6305f6fc6ca77f8233efa7ecd70d3d19c7420d

## pr

https://github.com/dphaener/rhiza/pull/159

## result

1 check passed: check SUCCESS

## statement

Pull request #159 (https://github.com/dphaener/rhiza/pull/159) passed CI at 0a6305f

## status

verified

## Links

### Out

- verifies [run-resume-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-resume/run-resume-refuses-agent-session.md)

## History

- `01M65G0WKX7Y74QCRB8EAN573P` A repo constraint states that no agent session starts a run, so every run-starting path (run start, run resume, a future daemon) is checked against it, not only rhiza run start
