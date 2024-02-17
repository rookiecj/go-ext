package langx

import (
	"reflect"
	"testing"
)

type myStruct struct {
	key int
	s   string
}

func TestCopy(t *testing.T) {
	type args[T any] struct {
		src     T
		options []func(ptr *T)
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}

	baseMyStruct := myStruct{
		key: 13,
		s:   "13",
	}

	tests := []testCase[myStruct]{
		{
			name: "no option",
			args: args[myStruct]{
				src: baseMyStruct,
			},
			want: baseMyStruct,
		},
		{
			name: "with option",
			args: args[myStruct]{
				src: baseMyStruct,
				options: []func(*myStruct){
					func(c *myStruct) {
						c.key = 14
					},
				},
			},
			want: myStruct{
				key: 14,
				s:   "13",
			},
		},
		{
			name: "last option win",
			args: args[myStruct]{
				src: baseMyStruct,
				options: []func(*myStruct){
					func(c *myStruct) {
						c.key = 14
					},
					func(c *myStruct) {
						c.key = 1
					},
				},
			},
			want: myStruct{
				key: 1,
				s:   "13",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Copy(tt.args.src, tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Copy() = %v, want %v", got, tt.want)
			}
		})
	}
}
