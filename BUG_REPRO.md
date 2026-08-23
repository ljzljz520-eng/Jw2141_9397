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
?   	example.com/xiangzhenfarm/cmd/farmregistry	[no test files]
?   	example.com/xiangzhenfarm/internal/quality	[no test files]
ok  	example.com/xiangzhenfarm/internal/cli	0.011s
ok  	example.com/xiangzhenfarm/internal/domain	0.001s
ok  	example.com/xiangzhenfarm/internal/importer	0.001s
ok  	example.com/xiangzhenfarm/internal/report	0.010s
--- FAIL: TestBusiness003Regression (0.01s)
    business_regression_test.go:35: invalid report was published: domain.VisitReport{ID:"report-C-103", FarmerID:"C-103", VisitDate:"2026-06-01", Officer:"Officer", Findings:"", MissingField:"findings", Status:"published", Sequence:25}
FAIL
FAIL	example.com/xiangzhenfarm/internal/service	0.019s
ok  	example.com/xiangzhenfarm/internal/store	0.011s
ok  	example.com/xiangzhenfarm/internal/validate	0.002s
ok  	example.com/xiangzhenfarm/internal/workflow	0.020s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/farmregistry): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/farmregistry): exit `0`
