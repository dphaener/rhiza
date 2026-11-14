# review-01m7ah6rayrzax6z6r63q0m2qk-r6-f3

question · `01M7BZ7EERY8PQN0AAF643PNEB`

## statement

Plan review round 6 left finding f3 (contradiction on handback-target-survives-include): inlined.outer in internal/workflow/include.go on main already maps handback:<step> to handback:<workflow>/<step> (CutPrefix on HandbackPrefix, landed in #572 by change set 01M79KGR, leaf handback-route). An unknown included step is already refused by Validate as not a step of the workflow. The leaf's premise that outer prefixes the whole target into <workflow>/handback:<step> does not hold. The critic confirmed it: Verified: origin/main internal/workflow/include.go lines 226-227 already does strings.CutPrefix(t, HandbackPrefix) and returns HandbackPrefix + workflow + "/" + step, so the leaf's premise (outer prefixes the whole target into <workflow>/handback:<step>) is false on main. The leaf contradicts the code it would change. Major stands. How should it be settled?

## Links

### Out

- blocks [handback-target-survives-include](../capability/rhiza/workflows/built-in-workflows/handback-target-survives-include.md)

### In

- answers [handback-include-already-mapped](../decision/handback-include-already-mapped.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
