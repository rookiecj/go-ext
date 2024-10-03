package logger

import (
	"context"
	"testing"
)

func TestLog(t *testing.T) {
	type args struct {
		ctx   context.Context
		level LogLevel
		msg   string
		args  []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Log - DebugLevel",
			args: args{
				ctx:   context.Background(),
				level: DebugLevel,
				msg:   "level: %s\n",
				args:  []any{"DebugLevel"},
			},
		},

		{
			name: "Log - InfoLevel",
			args: args{
				ctx:   context.Background(),
				level: InfoLevel,
				msg:   "level: %s\n",
				args:  []any{"InfoLevel"},
			},
		},

		{
			name: "Log - WarnLevel",
			args: args{
				ctx:   context.Background(),
				level: WarnLevel,
				msg:   "level: %s\n",
				args:  []any{"WarnLevel"},
			},
		},

		{
			name: "Log - ErrorLevel",
			args: args{
				ctx:   context.Background(),
				level: ErrorLevel,
				msg:   "level: %s\n",
				args:  []any{"ErrorLevel"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.args.ctx, tt.args.level, tt.args.msg, tt.args.args...)
		})
	}
}
