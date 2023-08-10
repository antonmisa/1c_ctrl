package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetKeyValue(t *testing.T) {
	cases := []struct {
		name      string
		line      string
		delimeter rune
		k         string
		v         string
		respError string
	}{
		{
			name:      "Success",
			line:      " id: test",
			delimeter: ':',
			k:         "id",
			v:         "test",
		},
		{
			name:      "Success w space",
			line:      " id  :   test  ",
			delimeter: ':',
			k:         "id",
			v:         "test",
		},
		{
			name:      "Success w tab & space",
			line:      `			 id  :   test  `,
			delimeter: ':',
			k:         "id",
			v:         "test",
		},
		{
			name:      "Success w tab & space 1",
			line:      `			 id  :   test value  `,
			delimeter: ':',
			k:         "id",
			v:         "test value",
		},
		{
			name:      "Not found",
			line:      " id - test",
			delimeter: ':',
			k:         "id",
			v:         "test",
			respError: ErrNotFound.Error(),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotKey, gotValue, err := GetKeyValue(tc.line, tc.delimeter)

			if err == nil {
				require.NoError(t, err)

				require.Equal(t, gotKey, tc.k)
				require.Equal(t, gotValue, tc.v)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, gotKey)
				require.Empty(t, gotValue)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		lines []string
		v     any
	}
	cases := []struct {
		name string
		args args
		res  any
		err  error
	}{
		{
			name: "OK cluster",
			args: args{
				lines: []string{
					"cluster:    test",
					"host: test",
				},
				v: Cluster{},
			},
			res: Cluster{
				ID:   "test",
				Host: "test",
			},
		},
		{
			name: "Not found cluster",
			args: args{
				lines: []string{
					"test:    test",
					"ghost: test",
				},
				v: Cluster{},
			},
			res: Cluster{},
			err: ErrNotFound,
		},
		{
			name: "OK infobase",
			args: args{
				lines: []string{
					"infobase:    test",
					"name: test1",
					"descr: test2",
				},
				v: Infobase{},
			},
			res: Infobase{
				ID:   "test",
				Name: "test1",
				Desc: "test2",
			},
		},
		{
			name: "Not found infobase",
			args: args{
				lines: []string{
					"test:    test",
					"ghost: test",
				},
				v: Infobase{},
			},
			res: Infobase{},
			err: ErrNotFound,
		},
		{
			name: "OK session",
			args: args{
				lines: []string{
					"session:    test",
					"host: test1",
					"started-at: 2023-08-08T10:48:43",
				},
				v: Session{},
			},
			res: Session{
				ID:      "test",
				Host:    "test1",
				Started: time.Date(2023, time.August, 8, 10, 48, 43, 0, time.UTC),
			},
		},
		{
			name: "Not found session",
			args: args{
				lines: []string{
					"test:    test",
					"ghost: test",
				},
				v: Session{},
			},
			res: Session{},
			err: ErrNotFound,
		},
		{
			name: "OK connection",
			args: args{
				lines: []string{
					"connection:    test",
					"host: test1",
				},
				v: Connection{},
			},
			res: Connection{
				ID:   "test",
				Host: "test1",
			},
		},
		{
			name: "Not found connection",
			args: args{
				lines: []string{
					"test:    test",
					"ghost: test",
				},
				v: Connection{},
			},
			res: Connection{},
			err: ErrNotFound,
		},
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var err error

			switch tc.args.v.(type) {
			case Cluster:
				var v Cluster
				err = Unmarshal(tc.args.lines, &v)

				require.Equal(t, v, tc.res)
			case Infobase:
				var v Infobase
				err = Unmarshal(tc.args.lines, &v)

				require.Equal(t, v, tc.res)
			case Session:
				var v Session
				err = Unmarshal(tc.args.lines, &v)

				require.Equal(t, v, tc.res)
			case Connection:
				var v Connection
				err = Unmarshal(tc.args.lines, &v)

				require.Equal(t, v, tc.res)
			default:
			}

			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, tc.err, err)
			}
		})
	}
}
