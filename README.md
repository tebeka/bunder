# bunder - Check that Go Benchmark(s) Under Given Threshold

## Usage

```
usage: bunder [options] [FILE]
Checks that benchmark results are below thresholds.

  -config string
    	config file name (default ".bunder.yml")
  -version
    	show version and exit
```

## Example

```
$ go test -run NONE -bench . -benchmem -count 5 | tee wrand.txt
goos: linux
goarch: amd64
pkg: github.com/tebeka/wrand
cpu: 12th Gen Intel(R) Core(TM) i7-1255U
BenchmarkRand-12       	38350470	        29.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkRand-12       	40832551	        29.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkRand-12       	41504181	        29.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkRand-12       	45680728	        29.20 ns/op	       0 B/op	       0 allocs/op
BenchmarkRand-12       	45664462	        26.21 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandBig-12    	 8707904	       128.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandBig-12    	 9069816	       130.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandBig-12    	 8862433	       128.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandBig-12    	 9341898	       128.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkRandBig-12    	 9340988	       128.8 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	github.com/tebeka/wrand	13.793s

$ cat wrand.yml 
version: v1
thresholds:
  - name: BenchmarkRand-12
    threshold: 30ns
  - name: BenchmarkRandBig-12
    threshold: 0.112us

$ bunder -config wrand.yml wrand.txt 
BenchmarkRandBig-12: avg = 128.94 ns, threshold = 112ns
```

## Configuration

bunder will read configuration with benchmark thresholds. It should be a YAML file in the following format:

```yaml
version: v1
thresholds:
  - name: BenchmarkRand-12
    threshold: 30ns
  - name: BenchmarkRandBig-12
    threshold: 0.112us
```

`threshold` can be anything that [time.ParseDuration](https://pkg.go.dev/time#ParseDuration) can handle.

## Installing

`go install github.com/tebeka/bunder@latest`

The installed file will be in `$(go env GOPATH)/bin`.
