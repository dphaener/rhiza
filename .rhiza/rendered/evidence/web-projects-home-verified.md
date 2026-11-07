# web-projects-home-verified

evidence · `01M6STF5QFCC84K8896H9ET0RS` · actor [claude-code-f1a799a6](../actor/claude-code-f1a799a6.md)

## result

PASS

## statement

Verifier claude-code-f1a799a6 at e002685c (leaf commits e3cd3043, 87254fce, e002685c). Read the diffs. internal/web/projects_test.go TestProjectsHomeReadsTheRoot writes a real rhiza.toml with two initialised projects (alpha, beta) and one bare git checkout, serves /projects through api.WebRootHandler (so api.WebProjectsSource is exercised), and checks each row: alpha '1 running, 1 waiting on you · 2 in inbox', beta 'none live · 0 in inbox', each with a health pill equal to that project's own health op state, bare marked 'not initialised' with no counts or health, and every row linking /p/<name>/. It also checks that /, /runs, /roadmap, /inbox, /p/beta/ and /p/beta/runs link /projects in their nav. The attempt-2 flake (TempDir RemoveAll 'directory not empty') is fixed by settledDir, which retries removal until the dir stays gone for 1s. In a clean worktree at e002685c under load average about 7-8: go test -trimpath -count=1 ./internal/web/ passed 8 of 8 runs; -count=15 -run TestProjectsHome passed; no rhiza-web-projects-* dirs were left in /tmp. go vet ./internal/web/ was clean. In the run checkout, TestOverviewReadsDaemonFromRepo and api TestSessionRunRecordsLimitHit fail, but both pass at the same commit in a clean worktree (and the api one passes at base e3cd3043^), so they depend on the checkout's location, not on this leaf.

## Links

### Out

- verifies [web-projects-home](../capability/rhiza/web-projects-home.md)

## History

- `01M6RN3ZBGYDW6KPBSES9A5JXS` Web UI across projects: an all-projects home, a merged inbox tagged by project, and every page scoped to one project through a switcher
