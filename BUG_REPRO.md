# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	lawindex/cmd/lawindex	[no test files]
ok  	lawindex/internal/api	0.014s
ok  	lawindex/internal/archive	0.020s
ok  	lawindex/internal/catalog	0.019s
ok  	lawindex/internal/model	0.002s
ok  	lawindex/internal/report	0.001s
ok  	lawindex/internal/review	0.016s
ok  	lawindex/internal/store	0.019s
--- FAIL: TestOptionalScoringDisabledMessage (0.01s)
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x16121c]

goroutine 19 [running]:
testing.tRunner.func1.2({0x2055e0, 0x40e8d0})
	/usr/local/go/src/testing/testing.go:1631 +0x1c4
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1634 +0x33c
panic({0x2055e0?, 0x40e8d0?})
	/usr/local/go/src/runtime/panic.go:770 +0x124
lawindex/internal/review.(*Service).RequestScore(0x4000021e60, {0x400017a000, 0x18})
	/app/internal/review/review.go:34 +0xec
lawindex/internal/workflow30.TestOptionalScoringDisabledMessage(0x400014c680)
	/app/internal/workflow30/optional_scoring_test.go:27 +0x22c
testing.tRunner(0x400014c680, 0x26c8d8)
	/usr/local/go/src/testing/testing.go:1689 +0xec
created by testing.(*T).Run in goroutine 1
	/usr/local/go/src/testing/testing.go:1742 +0x318
FAIL	lawindex/internal/workflow30	0.018s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/lawindex): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/lawindex): exit `0`
- Frontend build (web): exit `0`
