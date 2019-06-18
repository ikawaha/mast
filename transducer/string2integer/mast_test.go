package string2integer

import (
	"os"
	"reflect"
	"testing"
)

func TestMast_Run(t *testing.T) {
	t.Run("calendar", func(t *testing.T) {
		input := PairList{
			{"apr", 30},
			{"aug", 31},
			{"dec", 31},
			{"feb", 28},
			{"feb", 29},
		}
		m := BuildMast(input)

		cases := []struct {
			In       string
			Expected []int
		}{
			{In: "apr", Expected: []int{30}},
			{In: "aug", Expected: []int{31}},
			{In: "dec", Expected: []int{31}},
			{In: "feb", Expected: []int{28, 29}},
		}
		for i, v := range cases {
			if outs, ok := m.Run(v.In); !ok {
				t.Fatalf("%v: expected accepting '%v', but not accept", i, v.In)
			} else if expected := v.Expected; !reflect.DeepEqual(expected, outs) {
				t.Fatalf("%v: expected %+v, got %v", i, expected, outs)
			}
		}
	})
	t.Run("lucene example", func(t *testing.T) {
		input := PairList{
			{"mop", 0},
			{"moth", 1},
			{"pop", 2},
			{"star", 3},
			{"stop", 4},
			{"top", 5},
		}
		m := BuildMast(input)
		cases := []struct {
			In       string
			Expected []int
		}{
			{In: "mop", Expected: []int{0}},
			{In: "moth", Expected: []int{1}},
			{In: "pop", Expected: []int{2}},
			{In: "star", Expected: []int{3}},
			{In: "stop", Expected: []int{4}},
			{In: "top", Expected: []int{5}},
		}
		for i, v := range cases {
			if outs, ok := m.Run(v.In); !ok {
				t.Fatalf("%v: expected accepting '%v', but not accept", i, v.In)
			} else if expected := v.Expected; !reflect.DeepEqual(expected, outs) {
				t.Fatalf("%v: expected %+v, got %v", i, expected, outs)
			}
		}
	})
	t.Run("common prefix", func(t *testing.T) {
		input := PairList{
			{"he", 0},
			{"hell", 1},
			{"hello", 2},
		}
		m := BuildMast(input)
		cases := []struct {
			In       string
			Expected []int
		}{
			{In: "he", Expected: []int{0}},
			{In: "hell", Expected: []int{1}},
			{In: "hello", Expected: []int{2}},
		}
		for i, v := range cases {
			if outs, ok := m.Run(v.In); !ok {
				t.Fatalf("%v: expected accepting '%v', but not accept", i, v.In)
			} else if expected := v.Expected; !reflect.DeepEqual(expected, outs) {
				t.Fatalf("%v: expected %+v, got %v", i, expected, outs)
			}
		}
	})
}

func TestMast_Dot(t *testing.T) {
	t.Run("calendar", func(t *testing.T) {
		input := PairList{
			{"apr", 30},
			{"aug", 31},
			{"dec", 31},
			{"feb", 28},
			{"feb", 29},
		}
		m := BuildMast(input)
		if err := m.Dot(os.Stdout); err != nil {
			t.Fatal("unexpected error", err)
		}
	})
	t.Run("lucene example", func(t *testing.T) {
		input := PairList{
			{"mop", 0},
			{"moth", 1},
			{"pop", 2},
			{"star", 3},
			{"stop", 4},
			{"top", 5},
		}
		m := BuildMast(input)
		if err := m.Dot(os.Stdout); err != nil {
			t.Fatal("unexpected error", err)
		}
	})
	t.Run("common prefix", func(t *testing.T) {
		input := PairList{
			{"he", 0},
			{"hell", 1},
			{"hello", 2},
		}
		m := BuildMast(input)
		if err := m.Dot(os.Stdout); err != nil {
			t.Fatal("unexpected error", err)
		}
	})
}
