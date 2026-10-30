# failed-fixer-arbiter-ruling

source · `01M6506B3NZPEJG7ZB72C713FA`

## quotes

```json
[]
```

## raw_sha256

865c60929a82fbdc7535cf08b3debff84cf3883fa75608d6baf9b482a12334c5

## session

33462c8e-2065-4b43-a930-01a18a012b7a

## summary

Autonomous plan run with no owner present, so no owner words were captured. The question raised by the plan run, how many times a failed review fixer may be retried inside its cycle, was answered by the arbiter: a failed review fixer is retried within its cycle up to MaxFixCycles more times, and only when those retries are exhausted is the cycle spent. The checks of a spent cycle then fail with the fix-failed detail. Open: none for this change set.

## title

Failed fixer keeps its cycle: arbiter ruling

## History

- `01M61GW7TJ6BH6YX1T7SQDWK1J` A failed fixer does not use up a review fix cycle: the fixer is retried within its cycle, and the cycle's checks re-run only once its fixers pass.
