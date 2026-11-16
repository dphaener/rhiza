# cli-child-env-qa-2

evidence · `01M7K0CW42MEHF027C28F7MSXG` · actor [claude-code-037cc825](../actor/claude-code-037cc825.md)

## result

PASS

## statement

QA cycle 2: scratch builds of HEAD 189c5029e and merge-base ce69fddc6 outside PATH; repo with a planted ./git, env -i, scratch HOME, PATH='.:/usr/bin:/bin' and ':/usr/bin:/bin': rhiza init, change new and status never ran the planted git on either build (HEAD now refuses via exec.ErrDot like the merge-base); PATH='/usr/bin:/bin' writes git config. go test envexec/userdir/run/index pass; scripts usage-live-record, project-add-url/path, project-ls-rm, root-init pass.

## Links

### Out

- verifies [cli-child-env](../capability/rhiza/integrity/self-check/cli-child-env.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
