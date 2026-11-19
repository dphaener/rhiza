# false-hand-commit-plan-research

source · `01M7RJSVDV1TJADY14ZGGWJ0SR`

## quotes

```json
[]
```

## raw_sha256

c989c0fcfb6a1c591b67b52f46177c1fa50b520c976c948058bd97109fe5406a

## session

072d7607-1796-4a40-ae4a-01a1f12a770b

## summary

Autonomous plan, no owner interview. Read: run 01M7RHR8NEJTECVZ5FYFDJ0H18's status journals one hand-commit intervention for 351cf62f6 'Land and render the run-status flake fix', the conductor's own commit after rhiza land and rhiza render; its message has Change: in a body paragraph and only Co-Authored-By in the trailer block, so no Actor:, Leaf: or Review: trailer. handCommits (internal/api/run_handcommit.go, Decision hand-commit-by-trailer) already treats a commit whose Actor: trailer names the actor of one of the run's sessions as the run's, and the conductor's session record carries its actor (claude-code-9239a99d); so the commit only needs an Actor: trailer. The conduct guide (internal/guide/text/conduct.md, conduct-handback.md) says 'rhiza land, then rhiza render, and commit what they change' with no trailer rule; the commit is written by hand. Ship's own land (internal/api/ship_close.go) already commits .rhiza with 'Change: <cs>' and 'Actor: <actor handle>' trailers, the actor resolved as every write resolves it (shipActorHandle). Plan: rhiza land gains --commit, which lands, renders and commits .rhiza with those trailers, so the conductor's land commit carries its session actor without relying on the agent; the conduct guide tells the conductor to use it and to trailer any other commit; a test shows such a commit is not journalled as an intervention. Leaving handCommits' rule unchanged keeps a person's hand fix with only a Change: trailer counted as an intervention.

## title

False hand-commit plan research (autonomous plan of 01M7RJM8)

## Links

### In

- derived-from [conductor-land-commit-trailers](../decision/conductor-land-commit-trailers.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
