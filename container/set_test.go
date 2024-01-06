package container

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"testing"
)

type testSetStruct struct {
	Name string
	Age  int
}

func TestNewSet(t *testing.T) {
	type args[T comparable] struct {
		s Set[T]
	}

	testSet := NewSet[string]()

	tests := []struct {
		name string
		args args[string]
		want Set[string]
	}{
		{
			name: "empty",
			args: args[string]{
				s: testSet,
			},
			want: testSet,
		},

		{
			name: "insert",
			args: args[string]{
				s: func() Set[string] {
					set := NewSet[string]()
					set.Insert("1")
					set.Insert("2")
					set.Insert("3")
					return set
				}(),
			},
			want: func() Set[string] {
				set := NewSet[string]()
				set.Insert("1")
				set.Insert("2")
				set.Insert("3")
				return set
			}(),
		},

		{
			name: "remove",
			args: args[string]{
				s: func() Set[string] {
					set := NewSet[string]()
					set.Insert("1")
					set.Insert("2")
					set.Insert("3")
					set.Remove("2")
					return set
				}(),
			},
			want: func() Set[string] {
				set := NewSet[string]()
				set.Insert("1")
				//set.Insert("2")
				set.Insert("3")
				return set
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Equals(t *testing.T) {
	type args[T comparable] struct {
		other Set[T]
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want bool
	}
	tests := []testCase[testSetStruct]{
		{
			name: "empty",
			s:    NewSet[testSetStruct](),
			args: args[testSetStruct]{
				other: func() Set[testSetStruct] {
					return NewSet[testSetStruct]()
				}(),
			},
			want: true,
		},

		{
			name: "equals",
			s: func() Set[testSetStruct] {
				set := NewSet[testSetStruct]()
				set.Insert(testSetStruct{
					Name: "John",
					Age:  23,
				})
				set.Insert(testSetStruct{
					Name: "Jane",
					Age:  21,
				})
				return set
			}(),
			args: args[testSetStruct]{
				other: func() Set[testSetStruct] {
					set := NewSet[testSetStruct]()
					set.Insert(testSetStruct{
						Name: "John",
						Age:  23,
					})
					set.Insert(testSetStruct{
						Name: "Jane",
						Age:  21,
					})
					return set
				}(),
			},
			want: true,
		},

		{
			name: "equals - false",
			s: func() Set[testSetStruct] {
				set := NewSet[testSetStruct]()
				set.Insert(testSetStruct{
					Name: "John",
					Age:  23,
				})
				return set
			}(),
			args: args[testSetStruct]{
				other: func() Set[testSetStruct] {
					set := NewSet[testSetStruct]()
					set.Insert(testSetStruct{
						Name: "John",
						Age:  24,
					})
					return set
				}(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Equals(tt.args.other); got != tt.want {
				t.Errorf("Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Len(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		want int
	}

	limit := 1024

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			want: 0,
		},
		{
			name: "1",
			s: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				return s
			}(),
			want: 1,
		},
		{
			name: "limit",
			s: func() Set[string] {
				s := NewSet[string]()
				for idx := 0; idx < limit; idx++ {
					s.Insert(fmt.Sprintf("%d", idx))
				}
				return s
			}(),
			want: limit,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Has(t *testing.T) {
	type args[T comparable] struct {
		element T
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want bool
	}
	tests := []testCase[testSetStruct]{
		{
			name: "empty - false",
			s:    NewSet[testSetStruct](),
			args: args[testSetStruct]{
				element: testSetStruct{
					Name: "John",
					Age:  23,
				},
			},
			want: false,
		},
		{
			name: "1",
			s: func() Set[testSetStruct] {
				s := NewSet[testSetStruct]()
				s.Insert(testSetStruct{
					Name: "John",
					Age:  23,
				})
				return s
			}(),
			args: args[testSetStruct]{
				element: testSetStruct{
					Name: "John",
					Age:  23,
				},
			},
			want: true,
		},
		{
			name: "last",
			s: func() Set[testSetStruct] {
				s := NewSet[testSetStruct]()
				s.Insert(testSetStruct{
					Name: "Sam",
					Age:  48,
				})
				s.Insert(testSetStruct{
					Name: "Jane",
					Age:  21,
				})
				s.Insert(testSetStruct{
					Name: "John",
					Age:  23,
				})
				return s
			}(),
			args: args[testSetStruct]{
				element: testSetStruct{
					Name: "John",
					Age:  23,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Has(tt.args.element); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Insert(t *testing.T) {
	type args[T comparable] struct {
		elements int
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want int
	}

	limit := 1000000

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{},
			want: 0,
		},
		{
			name: "1",
			s:    NewSet[string](),
			args: args[string]{
				elements: 1,
			},
			want: 1,
		},
		{
			name: "limit",
			s:    NewSet[string](),
			args: args[string]{
				elements: limit,
			},
			want: limit,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for idx := 0; idx < tt.args.elements; idx++ {
				tt.s.Insert(fmt.Sprintf("%d", idx))
			}

			if tt.s.Len() != tt.want {
				t.Errorf("Insert() want %v got %v", tt.want, tt.s.Len())
			}
		})
	}
}

func TestSet_InsertSame(t *testing.T) {
	type args[T comparable] struct {
		elements int
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want int
	}

	limit := 1000000

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{},
			want: 0,
		},
		{
			name: "1",
			s:    NewSet[string](),
			args: args[string]{
				elements: 1,
			},
			want: 1,
		},
		{
			name: "limit",
			s:    NewSet[string](),
			args: args[string]{
				elements: limit,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for idx := 0; idx < tt.args.elements; idx++ {
				tt.s.Insert("1")
			}

			if tt.s.Len() != tt.want {
				t.Errorf("InsertSame() want %v got %v", tt.want, tt.s.Len())
			}
		})
	}
}

func TestSet_InsertRand(t *testing.T) {
	type args[T comparable] struct {
		elements int
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want int
	}

	limit := 1000000

	tests := []testCase[string]{
		{
			name: "limit",
			s:    NewSet[string](),
			args: args[string]{
				elements: limit,
			},
			want: limit,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for idx := 0; idx < tt.args.elements; idx++ {
				r := rand.Intn(tt.args.elements)
				tt.s.Insert(fmt.Sprintf("%d", r))
			}

			if tt.s.Len() > tt.want {
				t.Errorf("InsertRand() want %v got %v", tt.want, tt.s.Len())
			}
		})
	}
}

func TestSet_Iterate(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			want: nil,
		},
		{
			name: "add 3 in order",
			s: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")

				return s
			}(),
			want: []string{"1", "2", "3"},
		},
		{
			name: "add 3 disorder",
			s: func() Set[string] {
				s := NewSet[string]()
				s.Insert("2")
				s.Insert("1")
				s.Insert("3")

				return s
			}(),
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Iterate()
			slices.Sort(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Remove(t *testing.T) {
	type args[T comparable] struct {
		elements []T
		removes  []T
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want Set[T]
	}

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{
				elements: nil,
				removes:  nil,
			},
			want: NewSet[string](),
		},
		{
			name: "add 1 - remove 0",
			s:    NewSet[string](),
			args: args[string]{
				elements: []string{"1"},
				removes:  nil,
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				return s
			}(),
		},
		{
			name: "add 1 - remove 1 - empty",
			s:    NewSet[string](),
			args: args[string]{
				elements: []string{"1"},
				removes:  []string{"1"},
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
		{
			name: "add 3 in order - remove 3 - empty",
			s:    NewSet[string](),
			args: args[string]{
				elements: []string{"1", "2", "3"},
				removes:  []string{"1", "2", "3"},
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
		{
			name: "add 3 - remove 3 - empty",
			s:    NewSet[string](),
			args: args[string]{
				elements: []string{"2", "1", "3"},
				removes:  []string{"1", "2", "3"},
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
		{
			name: "add 3 same - remove 1 - empty",
			s:    NewSet[string](),
			args: args[string]{
				elements: []string{"1", "1", "1"},
				removes:  []string{"1"},
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, ele := range tt.args.elements {
				tt.s.Insert(ele)
			}
			for _, ele := range tt.args.removes {
				tt.s.Remove(ele)
			}
			if !reflect.DeepEqual(tt.want, tt.s) {
				t.Errorf("Remove want %v, got %v", tt.want, tt.s)
			}
		})
	}
}

func TestSet_Intersection(t *testing.T) {
	type args[T comparable] struct {
		other Set[T]
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want Set[T]
	}

	setA := NewSet[string]()
	setA.Insert("1")
	setA.Insert("2")
	setA.Insert("3")
	setA.Insert("4")

	setB := NewSet[string]()
	setB.Insert("2")
	setB.Insert("3")

	setC := NewSet[string]()
	setC.Insert("100")
	setC.Insert("101")
	setC.Insert("102")
	setC.Insert("103")
	setC.Insert("104")

	setD := NewSet[string]()
	setD.Insert("2")
	setD.Insert("103")

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{
				other: NewSet[string](),
			},
			want: NewSet[string](),
		},
		{
			name: "A,A - A",
			s:    setA,
			args: args[string]{
				other: setA,
			},
			want: setA,
		},
		{
			name: "A,C - none",
			s:    setA,
			args: args[string]{
				other: setC,
			},
			want: NewSet[string](),
		},
		{
			name: "A,B - contains b",
			s:    setA,
			args: args[string]{
				other: setB,
			},
			want: setB,
		},
		{
			name: "A,D - 1",
			s:    setA,
			args: args[string]{
				other: setD,
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("2")
				return s
			}(),
		},
		{
			name: "C,A - none",
			s:    setC,
			args: args[string]{
				other: setA,
			},
			want: NewSet[string](),
		},
		{
			name: "B,A - contains b",
			s:    setB,
			args: args[string]{
				other: setA,
			},
			want: setB,
		},
		{
			name: "D,A - 1",
			s:    setD,
			args: args[string]{
				other: setA,
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("2")
				return s
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Intersection(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Union(t *testing.T) {
	type args[T comparable] struct {
		other Set[T]
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want Set[T]
	}

	setA := NewSet[string]()
	setA.Insert("1")
	setA.Insert("2")
	setA.Insert("3")
	setA.Insert("4")

	setB := NewSet[string]()
	setB.Insert("2")
	setB.Insert("3")

	setC := NewSet[string]()
	setC.Insert("100")
	setC.Insert("101")
	setC.Insert("102")
	setC.Insert("103")
	setC.Insert("104")

	setD := NewSet[string]()
	setD.Insert("2")
	setD.Insert("103")

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{
				other: NewSet[string](),
			},
			want: NewSet[string](),
		},
		{
			name: "A,A - A",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setA.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")
				s.Insert("4")
				return s
			}(),
		},
		{
			name: "A,B - contains B - A",
			s:    setA.Copy(),
			args: args[string]{
				other: setB.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")
				s.Insert("4")
				return s
			}(),
		},
		{
			name: "A,C - A+C",
			s:    setA.Copy(),
			args: args[string]{
				other: setC.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()

				// A
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")
				s.Insert("4")

				// C
				s.Insert("100")
				s.Insert("101")
				s.Insert("102")
				s.Insert("103")
				s.Insert("104")
				return s
			}(),
		},
		{
			name: "A,D - A+D",
			s:    setA.Copy(),
			args: args[string]{
				other: setD.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()

				// A
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")
				s.Insert("4")

				// D
				s.Insert("2")
				s.Insert("103")
				return s
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			if got := tt.s.Union(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Difference(t *testing.T) {
	type args[T comparable] struct {
		other Set[T]
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want Set[T]
	}

	setA := NewSet[string]()
	setA.Insert("1")
	setA.Insert("2")
	setA.Insert("3")
	setA.Insert("4")

	setB := NewSet[string]()
	setB.Insert("2")
	setB.Insert("3")

	setC := NewSet[string]()
	setC.Insert("100")
	setC.Insert("101")
	setC.Insert("102")
	setC.Insert("103")
	setC.Insert("104")

	setD := NewSet[string]()
	setD.Insert("2")
	setD.Insert("103")

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{
				other: NewSet[string](),
			},
			want: NewSet[string](),
		},
		{
			name: "A,A - empty",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setA.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
		{
			name: "A - B",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setB.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("4")
				return s
			}(),
		},
		{
			name: "A - C = A",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setC.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")
				s.Insert("4")
				return s
			}(),
		},
		{
			name: "B - A = empty",
			s: func() Set[string] {
				s := setB.Copy()
				return s
			}(),
			args: args[string]{
				other: setA.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
		{
			name: "C - A",
			s: func() Set[string] {
				s := setC.Copy()
				return s
			}(),
			args: args[string]{
				other: setA.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("100")
				s.Insert("101")
				s.Insert("102")
				s.Insert("103")
				s.Insert("104")
				return s
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Difference(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_SymmetricDifference(t *testing.T) {
	type args[T comparable] struct {
		other Set[T]
	}
	type testCase[T comparable] struct {
		name string
		s    Set[T]
		args args[T]
		want Set[T]
	}

	setA := NewSet[string]()
	setA.Insert("1")
	setA.Insert("2")
	setA.Insert("3")
	setA.Insert("4")

	setB := NewSet[string]()
	setB.Insert("2")
	setB.Insert("3")

	setC := NewSet[string]()
	setC.Insert("100")
	setC.Insert("101")
	setC.Insert("102")
	setC.Insert("103")
	setC.Insert("104")

	setD := NewSet[string]()
	setD.Insert("2")
	setD.Insert("103")

	tests := []testCase[string]{
		{
			name: "empty",
			s:    NewSet[string](),
			args: args[string]{
				other: NewSet[string](),
			},
			want: NewSet[string](),
		},
		{
			name: "A,A - empty",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setA.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				return s
			}(),
		},
		{
			name: "A - B",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setB.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("4")
				return s
			}(),
		},
		{
			name: "A - C",
			s: func() Set[string] {
				s := setA.Copy()
				return s
			}(),
			args: args[string]{
				other: setC.Copy(),
			},
			want: func() Set[string] {
				s := NewSet[string]()
				s.Insert("1")
				s.Insert("2")
				s.Insert("3")
				s.Insert("4")

				s.Insert("100")
				s.Insert("101")
				s.Insert("102")
				s.Insert("103")
				s.Insert("104")
				return s
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.SymmetricDifference(tt.args.other); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SymmetricDifference() = %v, want %v", got, tt.want)
			}
		})
	}
}
