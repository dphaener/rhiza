# plan-request-plans-serve-milestone

source · `01M6JDZZPKS97TCJ3ZYA7H0QZP`

## quotes

```json
[]
```

## raw_sha256

c8008504b120c292f7b842bc038c4d4e7342cc8dd61d6b261ed00f26e0651675

## session

6b593167-fb44-4178-8462-01a1a4da96c7

## summary

Autonomous plan run on change set 01M6JDM. No owner present; the planner raised four questions as one (serves intent, how a request's plan reaches a milestone, the guard, the requests TUI field) and the arbiter answered with Decision request-plans-serve-a-milestone: the change set serves milestone-25; rhiza request gains an optional --milestone flag naming a milestone intent, stored on the request and handed to the planner; with none named the planner serves the open milestone it fits, never an umbrella or area intent under no milestone; the plan step's completion check fails the step when a planned change set serves no milestone directly or through an intent under one; the requests TUI app gets an optional milestone field defaulting to auto. The planner then wrote five leaves under plan-requests. Nothing left open.

## title

Plan of fix: request plans serve a milestone (autonomous, arbiter)

## Links

### In

- derived-from [request-plans-serve-a-milestone](../decision/request-plans-serve-a-milestone.md)

## History

- `01M6JDMJTSMFWS7FC804Z9FT76` Fix Plan runs started from the request queue produce change sets the daemon never builds: the planner links them serves to intents only, and the daemon skips any change set that serves no milestone (paging 01M6H3A65RDP and usage governor 01M6HHMTEDFK sat skipped as "serves no milestone" after being planned and accepted). A request's plan should link its change set to a milestone: the open milestone it fits, or a backlog milestone the request names, so plan=auto with accept=auto ships requests end to end
