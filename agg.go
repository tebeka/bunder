package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// AggFn assumes values is not empty
type AggFn func(values []float64) float64

func avg(values []float64) float64 {
	total := 0.0
	for _, v := range values {
		total += v
	}

	return total / float64(len(values))
}

func percentile(p int, values []float64) float64 {
	vs := make([]float64, len(values))
	copy(vs, values)
	sort.Float64s(vs)

	v := float64(p) / 100
	i := int(math.Floor(float64(len(vs)) * v))

	if len(vs)%2 == 1 {
		return vs[i]
	}

	m := (vs[i-1] + vs[i]) / 2
	return m
}

func minAgg(values []float64) float64 {
	m := values[0]
	for _, v := range values[1:] {
		if v < m {
			m = v
		}
	}

	return m
}

func maxAgg(values []float64) float64 {
	m := values[0]
	for _, v := range values[1:] {
		if v > m {
			m = v
		}
	}

	return m
}

func aggByName(name string) (AggFn, error) {
	if name == "" {
		return nil, fmt.Errorf("empty agg name")
	}

	switch name {
	case "mean", "avg":
		return avg, nil
	case "min":
		return minAgg, nil
	case "max":
		return maxAgg, nil
	}

	if name[0] != 'p' {
		return nil, fmt.Errorf("%q - unknown agg", name)
	}

	p, err := strconv.ParseInt(name[1:], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%q - bad int (%w)", name, err)
	}

	fn := func(values []float64) float64 {
		return percentile(int(p), values)
	}
	return AggFn(fn), nil
}
