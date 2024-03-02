package httpx

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewCsvBodyParser(t *testing.T) {
	type myStruct struct {
		Idx    int     `json:"idx"`
		Name   string  `json:"name"`
		Height float64 `json:"height"`
	}

	csvStr := `0, john lee, 123.45
1, tom hank, 678.90
2, emily lee, 163.00
`
	csvStrWithHeader := `idx, name, height
0, john lee, 123.45
1, tom hank, 678.90
2, emily lee, 163.00
`
	csvStrWithQuote := `idx, name, height
0, "john lee", 123.45
1, "tom hank", 678.90
2, "emily lee", 163.00
`
	csvStrWithComment := `idx, name, height
# this is comment
0, "john lee", 123.45
1, "tom hank", 678.90
2, "emily lee", 163.00
`
	myStructs := []myStruct{
		{
			Idx:    0,
			Name:   "john lee",
			Height: 123.45,
		},
		{
			Idx:    1,
			Name:   "tom hank",
			Height: 678.90,
		},
		{
			Idx:    2,
			Name:   "emily lee",
			Height: 163.00,
		},
	}

	type args struct {
		hasHeader bool
		delim     rune
		comment   rune
	}
	tests := []struct {
		name string
		args args
		csv  string
		want []myStruct
	}{
		{
			name: "header no, comma",
			args: args{
				hasHeader: false,
				delim:     ',',
				comment:   '#',
			},
			csv:  csvStr,
			want: myStructs,
		},
		{
			name: "header yes, comma",
			args: args{
				hasHeader: true,
				delim:     ',',
				comment:   '#',
			},
			csv:  csvStrWithHeader,
			want: myStructs,
		},
		{
			name: "header yes, comma, quote",
			args: args{
				hasHeader: true,
				delim:     ',',
				comment:   '#',
			},
			csv:  csvStrWithQuote,
			want: myStructs,
		},
		{
			name: "header yes, comma, quote, comment",
			args: args{
				hasHeader: true,
				delim:     ',',
				comment:   '#',
			},
			csv:  csvStrWithComment,
			want: myStructs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewCsvBodyParser(tt.args.hasHeader, tt.args.delim, tt.args.comment)
			var got []myStruct
			if err := parser(strings.NewReader(tt.csv), &got); err != nil {
				t.Errorf("NewCsvBodyParser: parse error %v", err)
			} else if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("NewCsvBodyParser: want %v got %v\n", tt.want, got)
			}
		})
	}
}
