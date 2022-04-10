package maths_test

import (
	"github.com/itsnoproblem/pokt-lint/maths"
	"testing"
)

func TestMean(t *testing.T) {
	nums := []struct {
		Numbers  []float64
		Expected float64
	}{
		{
			Numbers:  []float64{10, 20, 30, 40, 50},
			Expected: 30,
		},
		{
			Numbers:  []float64{10, 20, 30, 40, 50, 60},
			Expected: 35,
		},
		{
			Numbers:  []float64{20, 20, 30, 70, 50, 50},
			Expected: 40,
		},
		{
			Numbers:  []float64{40, 20, 40, 0, 60, 50},
			Expected: 35,
		},
		{
			Numbers:  []float64{3, 1, 1, 50, 20},
			Expected: 15,
		},
		{
			Numbers:  []float64{0, 1},
			Expected: 0.5,
		},
		{
			Numbers:  []float64{11, 1},
			Expected: 6,
		},
		{
			Numbers:  []float64{},
			Expected: 0,
		},
	}
	for i, test := range nums {
		m := maths.Mean(test.Numbers)
		if m != test.Expected {
			t.Fatalf("%d of %d: expected [%f] got [%f]",
				i+1, len(nums), test.Expected, m)
		}
	}
}

func TestMedian(t *testing.T) {
	nums := []struct {
		Numbers  []float64
		Expected float64
	}{
		{
			Numbers:  []float64{10, 20, 30, 40, 50},
			Expected: 30,
		},
		{
			Numbers:  []float64{10, 20, 30, 40, 50, 60},
			Expected: 35,
		},
		{
			Numbers:  []float64{40, 20, 30, 10, 60, 50},
			Expected: 35,
		},
		{
			Numbers:  []float64{40, 20, 40, 0, 60, 50},
			Expected: 40,
		},
		{
			Numbers:  []float64{1, 1, 1, 1, 1, 50, 20},
			Expected: 1,
		},
		{
			Numbers:  []float64{},
			Expected: 0,
		},
	}
	for i, test := range nums {
		m := maths.Median(test.Numbers)
		if m != test.Expected {
			t.Fatalf("%d of %d: expected [%f] got [%f]",
				i+1, len(nums), test.Expected, m)
		}
	}
}

func TestMin(t *testing.T) {
	nums := []struct {
		Numbers  []float64
		Expected float64
	}{
		{
			Numbers:  []float64{2, 1, 3, 2, 31, 9},
			Expected: 1,
		},
		{
			Numbers:  []float64{10, 20, 30, 40, 50, 0},
			Expected: 0,
		},
		{
			Numbers:  []float64{},
			Expected: 0,
		},
		{
			Numbers:  []float64{99},
			Expected: 99,
		},
	}

	for i, tt := range nums {
		m := maths.Min(tt.Numbers)
		if m != tt.Expected {
			t.Fatalf("%d of %d: expected [%f] got [%f]", i+1, len(tt.Numbers), tt.Expected, m)
		}
	}
}
