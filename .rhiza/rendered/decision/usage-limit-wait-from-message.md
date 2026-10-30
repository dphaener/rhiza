# usage-limit-wait-from-message

decision · `01M65F4TAMG78HCQ0FFBTPBZ52`

## statement

A run recognises the usage limit from the limit message and waits until the reset time it gives; with no reset time in the message, it waits a fixed interval

## wording

Arbiter answer: the reset time in the message is the recommended wait

## Links

### Out

- answers [usage-limit-detect-and-wait](../question/usage-limit-detect-and-wait.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
