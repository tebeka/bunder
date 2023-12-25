package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// "BenchmarkRand-12       	40832551	        29.10 ns/op	       0 B/op	       0 allocs/op"
	// ->
	//  BenchmarkRand-12, 29.10
	benchRe = regexp.MustCompile(`(Benchmark[_A-Z][^\s]+)\s+\d+\s+(\d+(\.\d+))? ns/op`)
)

func parseLine(line string) (string, float64) {
	matches := benchRe.FindStringSubmatch(line)
	if len(matches) == 0 {
		return "", 0
	}

	d, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		// The regexp should make sure it's a float, if we get here then what?
		panic(fmt.Sprintf("can't parse float: %q", matches[2]))
	}

	return matches[1], d
}

func parseFile(r io.Reader) (map[string][]float64, error) {
	s := bufio.NewScanner(r)
	ds := make(map[string][]float64) // name -> duration, duration ...
	for s.Scan() {
		name, d := parseLine(s.Text())
		if name == "" {
			continue
		}
		ds[name] = append(ds[name], d)
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("can't scan - %w", err)
	}

	return ds, nil
}

func parseConfig(r io.Reader) (map[string]time.Duration, error) {
	var conf struct {
		Version    string
		Thresholds []struct {
			Name      string
			Threshold time.Duration
		}
	}

	dec := yaml.NewDecoder(r)
	if err := dec.Decode(&conf); err != nil {
		return nil, fmt.Errorf("can't decode YAML - %w", err)
	}

	if conf.Version != "v1" {
		return nil, fmt.Errorf("unknown config version: %q", conf.Version)
	}

	ts := make(map[string]time.Duration)
	for _, l := range conf.Thresholds {
		ts[l.Name] = l.Threshold
	}

	return ts, nil
}

func loadConfig(fileName string) (map[string]time.Duration, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseConfig(file)
}

type Benchmark struct {
	Name      string
	Avg       float64
	Threshold time.Duration
}

func avg(values []float64) float64 {
	// values can't be empty
	if len(values) == 0 {
		panic("avg: empty values")
	}

	total := 0.0
	for _, v := range values {
		total += v
	}

	return total / float64(len(values))
}

func findOffending(benches map[string][]float64, thresholds map[string]time.Duration) []Benchmark {
	var bad []Benchmark

	for name, durations := range benches {
		t, ok := thresholds[name]
		if !ok {
			continue
		}

		ba := avg(durations)
		if time.Duration(ba) > t {
			b := Benchmark{
				Name:      name,
				Avg:       ba,
				Threshold: t,
			}
			bad = append(bad, b)
		}
	}

	return bad
}

var usage = `usage: %s [options] [FILE]
Checks that benchmark results are below thresholds.

`

func main() {
	cfgFile := ".bunder.yml"
	flag.StringVar(&cfgFile, "config", cfgFile, "config file name")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "error: wrong number of arguments\n")
		os.Exit(1)
	}

	thresholds, err := loadConfig(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %q: can't load config - %s\n", cfgFile, err)
		os.Exit(1)
	}

	var r io.Reader = os.Stdin
	fileName := "<stdin>"
	if flag.NArg() == 1 {
		file, err := os.Open(flag.Arg(0))
		fileName = flag.Arg(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
		r = file
	}

	benches, err := parseFile(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %q: can't parse benchmark output - %s\n", fileName, err)
		os.Exit(1)
	}

	bad := findOffending(benches, thresholds)
	if len(bad) == 0 {
		os.Exit(0)
	}

	for _, b := range bad {
		fmt.Printf("%s: avg = %.2f ns, threshold = %v\n", b.Name, b.Avg, b.Threshold)
	}
}