// nolint
package pipe

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestPipe_New(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *Pipe
		wantErr error
	}{
		{
			name:    "File not exists",
			path:    "xxx://somewhatpath/some.file",
			want:    nil,
			wantErr: ErrNoFile,
		},
		{
			name: "Ok",
			path: "c:/windows/system32/ping.exe",
			want: func() *Pipe {
				p, _ := New("c:/windows/system32/ping.exe")
				return p
			}(),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := New(tt.path)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Pipe.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pipe.New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPipe_Run(t *testing.T) {
	type args struct {
		ctx context.Context
		arg []string
	}
	tests := []struct {
		name    string
		p       *Pipe
		args    args
		wantErr error
	}{
		{
			name: "Ping",
			p: func() *Pipe {
				p, _ := New("c:/windows/system32/ping.exe")
				return p
			}(),
			args: args{
				ctx: context.Background(),
				arg: make([]string, 0),
			},
			wantErr: nil,
		},
		{
			name: "Ping",
			p: func() *Pipe {
				p, _ := New("c:/windows/system32/ping.exe")
				return p
			}(),
			args: args{
				ctx: context.Background(),
				arg: []string{"localhost", "8080"},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got1, got2, err := tt.p.Run(tt.args.ctx, tt.args.arg...)
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("Pipe.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 == nil {
				t.Error("Pipe.Run() first return value is nil")
			}
			if got2 == nil {
				t.Error("Pipe.Run() second return value is nil")
			}
		})
	}
}
