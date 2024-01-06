package container

import (
	"cmp"
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestSortedMap_Equal(t *testing.T) {
	type args[K cmp.Ordered] struct {
		cmp func(a, b K) int
	}
	type testCase[K cmp.Ordered, V any] struct {
		name string
		c    SortedMap[K, V]
		args args[K]
		want []K
	}

	sortedMap := SortedMap[string, int]{}
	sortedMap["1"] = 1
	sortedMap["3"] = 3
	sortedMap["2"] = 2
	//sortedKeys := []string{"1", "2", "3"}

	t.Run("equal", func(t *testing.T) {
		newMap := SortedMap[string, int]{}
		newMap["3"] = 3
		newMap["2"] = 2
		newMap["1"] = 1

		want := sortedMap
		got := newMap
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Equal = %v, want %v", got, want)
		}
	})
}

func TestSortedMap_KeyOrderNotDefined(t *testing.T) {
	sortedMap := SortedMap[string, int]{}
	insertedOrder := []string{"1", "3", "2", "4", "9", "8", "7", "5", "6"}
	for _, k := range insertedOrder {
		if v, err := strconv.Atoi(k); err == nil {
			sortedMap[k] = v
		} else {
			panic(err)
		}
	}

	t.Run("no order", func(t *testing.T) {
		want := insertedOrder
		var got []string
		for key := range sortedMap {
			got = append(got, key)
		}

		fmt.Println("got", got)
		if reflect.DeepEqual(want, got) {
			t.Errorf("KeyOrderNotDefined = %v, want %v", got, want)
		}
	})
}

func TestSortedMap_LastManStanding(t *testing.T) {

	sortedMap := SortedMap[string, int]{}
	sortedMap["1"] = 3
	sortedMap["1"] = 13

	t.Run("LastMan", func(t *testing.T) {
		want := 13
		got := sortedMap["1"]
		if !reflect.DeepEqual(want, got) {
			t.Errorf("LastMan = %v, want %v", got, want)
		}
	})
}

func TestSortedMap_SortedKeys(t *testing.T) {
	type args[K cmp.Ordered] struct {
		cmp func(a, b K) int
	}

	type testCase[K cmp.Ordered, V any] struct {
		name string
		c    SortedMap[K, V]
		args args[K]
		want []K
	}

	sortedMap := SortedMap[string, int]{}
	sortedMap["1"] = 1
	sortedMap["3"] = 3
	sortedMap["2"] = 2
	sortedKeys := []string{"1", "2", "3"}
	sortedDescKeys := []string{"3", "2", "1"}

	tests := []testCase[string, int]{
		{
			name: "cmp nil - asc",
			c:    sortedMap,
			args: args[string]{
				cmp: nil,
			},
			want: sortedKeys,
		},

		{
			name: "cmp - asc",
			c:    sortedMap,
			args: args[string]{
				cmp: cmp.Compare[string],
			},
			want: sortedKeys,
		},

		{
			name: "cmp - desc",
			c:    sortedMap,
			args: args[string]{
				cmp: func(a, b string) int {
					return -cmp.Compare[string](a, b)
				},
			},
			want: sortedDescKeys,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.SortedKeys(tt.args.cmp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortedKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
