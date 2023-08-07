// nolint
package pipe

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/antonmisa/1cctl/internal/entity"
	"github.com/antonmisa/1cctl/pkg/pipe/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type FakeReadCloser struct {
	body []byte
	pos  int

	enable bool
}

func NewFakeConnection3() *FakeReadCloser {
	text := `connection : 1111-3434-5656 
				infobase: 3333-4444
				process: 1-1-1-1
				host: test-ic
				application: 1cv8

				connection : 2222-3434-5656 
				infobase: 3333-4444
				process: 1-1-1-2
				host: test-ic-1
				application: 1cv8
				
				connection : 3333-3434-5656 
				infobase: 1111-4444
				process: 1-2-1-2
				host: test-ic-2
				application: 1cv8`

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeConnection2() *FakeReadCloser {
	text := `connection : 1111-3434-5656 
				infobase: 3333-4444
				process: 1-1-1-1
				host: test-ic
				application: 1cv8

				connection : 2222-3434-5656 
				infobase: 3333-4444
				process: 1-1-1-2
				host: test-ic-1
				application: 1cv8`

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeSession3() *FakeReadCloser {
	text := `session : 1111-3434-5656 
				infobase: 3333-4444
				connection: 3-4-5-6
				process: 1-1-1-1
				user-name: тестовый пользователь
				host: test-ic
				app-id: 1cv8

				session : 2222-3434-5656 
				infobase: 3333-4444
				connection: 3-4-5-7
				process: 1-1-1-2
				user-name: тестовый пользователь 1
				host: test-ic-1
				app-id: 1cv8
				
				session : 3333-3434-5656 
				infobase: 1111-4444
				connection: 1-4-5-7
				process: 1-2-1-2
				user-name: неизвестный пользователь
				host: test-ic-2
				app-id: 1cv8`

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeSession2() *FakeReadCloser {
	text := `session : 1111-3434-5656 
				infobase: 3333-4444
				connection: 3-4-5-6
				process: 1-1-1-1
				user-name: тестовый пользователь
				host: test-ic
				app-id: 1cv8

				session : 2222-3434-5656 
				infobase: 3333-4444
				connection: 3-4-5-7
				process: 1-1-1-2
				user-name: тестовый пользователь 1
				host: test-ic-1
				app-id: 1cv8`

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeSession0() *FakeReadCloser {
	text := ``

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeInfobase() *FakeReadCloser {
	text := `infobase : 1212-3434-5656 
			 name: test
			 descr:

			 infobase : 1111-2222-3333 
			  name: test_ib
			 descr: "test desc" `

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func NewFakeCluster() *FakeReadCloser {
	text := `cluster : 1212-3434-5656 
			 host: localhost 
			 port: 1234 
			 name: "test"

			 cluster : 1111-2222-3333 
			  host: localhost.tnx.ru    
			 port: 1545 
			  name: "test cluster" `

	return &FakeReadCloser{
		body: []byte(text),
	}
}

func (t *FakeReadCloser) SetEnable(v bool) {
	t.enable = v
}

func (t *FakeReadCloser) Read(p []byte) (n int, err error) {
	if !t.enable {
		return 0, io.EOF
	}

	if t.pos >= len(t.body) {
		return 0, io.EOF
	}

	t.pos = copy(p, t.body[t.pos:])

	if t.pos >= len(t.body) {
		return t.pos, io.EOF
	} else {
		return t.pos, nil
	}
}

func (t *FakeReadCloser) Close() error {
	t.enable = false

	return nil
}

var _ io.ReadCloser = (*FakeReadCloser)(nil)

func TestGetClusters(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		stdout            *FakeReadCloser
		cls               []entity.Cluster
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name:              "Error no command",
			ctx:               context.Background(),
			cs:                "localhost:1545",
			stdout:            NewFakeCluster(),
			cls:               make([]entity.Cluster, 0),
			respError:         ": no command",
			pipeMockError:     errors.New("no command"),
			comMockStartError: errors.New("start error"),
		},
		{
			name:              "Error start",
			ctx:               context.Background(),
			cs:                "localhost:1545",
			stdout:            NewFakeCluster(),
			cls:               make([]entity.Cluster, 0),
			respError:         ": start error",
			comMockStartError: errors.New("start error"),
		},
		{
			name:             "Error wait",
			ctx:              context.Background(),
			cs:               "localhost:1545",
			stdout:           NewFakeCluster(),
			cls:              make([]entity.Cluster, 0),
			respError:        ": wait error",
			comMockWaitError: errors.New("wait error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError)

			ctrl := New(pipeMock)

			cls, err := ctrl.GetClusters(tc.ctx, tc.cs)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, cls)
				require.ElementsMatch(t, cls, tc.cls)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, cls)
			}
		})
	}
}

func TestGetInfobases(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		cred              entity.Credentials
		stdout            *FakeReadCloser
		ibs               []entity.Infobase
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Success wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			cred:   entity.Credentials{},
			stdout: NewFakeInfobase(),
			ibs: []entity.Infobase{
				{
					ID:   "1212-3434-5656",
					Name: "test",
					Desc: "",
				},
				{
					ID:   "1111-2222-3333",
					Name: "test_ib",
					Desc: "\"test desc\"",
				},
			},
		},
		{
			name: "Success w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeInfobase(),
			ibs: []entity.Infobase{
				{
					ID:   "1212-3434-5656",
					Name: "test",
					Desc: "",
				},
				{
					ID:   "1111-2222-3333",
					Name: "test_ib",
					Desc: "\"test desc\"",
				},
			},
		},
		{
			name:          "Error no command",
			ctx:           context.Background(),
			cs:            "localhost:1545",
			stdout:        NewFakeInfobase(),
			ibs:           make([]entity.Infobase, 0),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			ibs, err := ctrl.GetInfobases(tc.ctx, tc.cs, tc.cl, tc.cred)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, ibs)
				require.ElementsMatch(t, ibs, tc.ibs)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, ibs)
			}
		})
	}
}

func TestGetSessions(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		ib                entity.Infobase
		cred              entity.Credentials
		stdout            *FakeReadCloser
		res               []entity.Session
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Success w empty ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:     entity.Infobase{},
			cred:   entity.Credentials{},
			stdout: NewFakeSession3(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
				{
					ID:           "3333-3434-5656",
					InfobaseID:   "1111-4444",
					ConnectionID: "1-4-5-7",
					ProcessID:    "1-2-1-2",
					UserName:     "неизвестный пользователь",
					Host:         "test-ic-2",
					AppID:        "1cv8",
				},
			},
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			cred:   entity.Credentials{},
			stdout: NewFakeSession2(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
			},
		},
		{
			name: "Success w empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeSession3(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
				{
					ID:           "3333-3434-5656",
					InfobaseID:   "1111-4444",
					ConnectionID: "1-4-5-7",
					ProcessID:    "1-2-1-2",
					UserName:     "неизвестный пользователь",
					Host:         "test-ic-2",
					AppID:        "1cv8",
				},
			},
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeSession2(),
			res: []entity.Session{
				{
					ID:           "1111-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-6",
					ProcessID:    "1-1-1-1",
					UserName:     "тестовый пользователь",
					Host:         "test-ic",
					AppID:        "1cv8",
				},
				{
					ID:           "2222-3434-5656",
					InfobaseID:   "3333-4444",
					ConnectionID: "3-4-5-7",
					ProcessID:    "1-1-1-2",
					UserName:     "тестовый пользователь 1",
					Host:         "test-ic-1",
					AppID:        "1cv8",
				},
			},
		},
		{
			name:          "Error no command",
			ctx:           context.Background(),
			cs:            "localhost:1545",
			cl:            entity.Cluster{},
			ib:            entity.Infobase{},
			cred:          entity.Credentials{},
			stdout:        NewFakeSession3(),
			res:           make([]entity.Session, 0),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Once()

			ctrl := New(pipeMock)

			res, err := ctrl.GetSessions(tc.ctx, tc.cs, tc.cl, tc.ib, tc.cred)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, res)
				require.ElementsMatch(t, res, tc.res)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, res)
			}
		})
	}
}

func TestDisableSessions(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		ib                entity.Infobase
		clCred            entity.Credentials
		ibCred            entity.Credentials
		code              string
		stdout            *FakeReadCloser
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Error infobase empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:            entity.Infobase{},
			clCred:        entity.Credentials{},
			ibCred:        entity.Credentials{},
			code:          "12345",
			stdout:        NewFakeSession3(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			ibCred: entity.Credentials{},
			code:   "12345",
			stdout: NewFakeSession2(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code:          "12345",
			stdout:        NewFakeSession3(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code:   "12345",
			stdout: NewFakeSession2(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ib: entity.Infobase{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			ibCred:        entity.Credentials{},
			code:          "12345",
			stdout:        NewFakeSession0(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			err := ctrl.DisableSessions(tc.ctx, tc.cs, tc.cl, tc.ib, tc.clCred, tc.ibCred, tc.code)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestEnableSessions(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		ib                entity.Infobase
		clCred            entity.Credentials
		ibCred            entity.Credentials
		code              string
		stdout            *FakeReadCloser
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Error infobase empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:            entity.Infobase{},
			clCred:        entity.Credentials{},
			ibCred:        entity.Credentials{},
			code:          "12345",
			stdout:        NewFakeSession3(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			ibCred: entity.Credentials{},
			code:   "12345",
			stdout: NewFakeSession3(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code:          "12345",
			stdout:        NewFakeSession3(),
			respError:     ": infobase is empty",
			pipeMockError: ErrInfobaseIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			ibCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			code:   "12345",
			stdout: NewFakeSession2(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ib: entity.Infobase{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			ibCred:        entity.Credentials{},
			code:          "12345",
			stdout:        NewFakeSession0(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			err := ctrl.EnableSessions(tc.ctx, tc.cs, tc.cl, tc.ib, tc.clCred, tc.ibCred, tc.code)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		s                 entity.Session
		clCred            entity.Credentials
		stdout            *FakeReadCloser
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Error session is empty empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s:             entity.Session{},
			clCred:        entity.Credentials{},
			stdout:        NewFakeSession3(),
			respError:     ": session is empty",
			pipeMockError: ErrSessionIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Session{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			stdout: NewFakeSession3(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Session{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout:        NewFakeSession3(),
			respError:     ": session is empty",
			pipeMockError: ErrSessionIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Session{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeSession3(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			s: entity.Session{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			stdout:        NewFakeSession3(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			err := ctrl.DeleteSession(tc.ctx, tc.cs, tc.cl, tc.s, tc.clCred)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestDeleteSessions(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		ss                []entity.Session
		clCred            entity.Credentials
		stdout            *FakeReadCloser
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Success session is empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss:     make([]entity.Session, 0),
			clCred: entity.Credentials{},
			stdout: NewFakeSession3(),
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Session{
				{
					ID: "3333-4444",
				},
				{
					ID: "5555-4444",
				},
			},
			clCred: entity.Credentials{},
			stdout: NewFakeSession3(),
		},
		{
			name: "Success empty sessions w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Session{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeSession3(),
		},
		{
			name: "Success w sessions w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Session{
				{
					ID: "3333-4444",
				},
				{
					ID: "5555-4444",
				},
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeSession3(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ss: []entity.Session{
				{
					ID: "12",
				},
			},
			clCred:        entity.Credentials{},
			stdout:        NewFakeSession3(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			err := ctrl.DeleteSessions(tc.ctx, tc.cs, tc.cl, tc.ss, tc.clCred)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestGetConnections(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		ib                entity.Infobase
		cred              entity.Credentials
		stdout            *FakeReadCloser
		res               []entity.Connection
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Success w empty ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib:     entity.Infobase{},
			cred:   entity.Credentials{},
			stdout: NewFakeConnection3(),
			res: []entity.Connection{
				{
					ID:         "1111-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-1",
					Host:       "test-ic",
					AppID:      "1cv8",
				},
				{
					ID:         "2222-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-2",
					Host:       "test-ic-1",
					AppID:      "1cv8",
				},
				{
					ID:         "3333-3434-5656",
					InfobaseID: "1111-4444",
					ProcessID:  "1-2-1-2",
					Host:       "test-ic-2",
					AppID:      "1cv8",
				},
			},
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			cred:   entity.Credentials{},
			stdout: NewFakeConnection2(),
			res: []entity.Connection{
				{
					ID:         "1111-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-1",
					Host:       "test-ic",
					AppID:      "1cv8",
				},
				{
					ID:         "2222-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-2",
					Host:       "test-ic-1",
					AppID:      "1cv8",
				},
			},
		},
		{
			name: "Success w empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeConnection3(),
			res: []entity.Connection{
				{
					ID:         "1111-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-1",
					Host:       "test-ic",
					AppID:      "1cv8",
				},
				{
					ID:         "2222-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-2",
					Host:       "test-ic-1",
					AppID:      "1cv8",
				},
				{
					ID:         "3333-3434-5656",
					InfobaseID: "1111-4444",
					ProcessID:  "1-2-1-2",
					Host:       "test-ic-2",
					AppID:      "1cv8",
				},
			},
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ib: entity.Infobase{
				ID: "3333-4444",
			},
			cred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeConnection2(),
			res: []entity.Connection{
				{
					ID:         "1111-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-1",
					Host:       "test-ic",
					AppID:      "1cv8",
				},
				{
					ID:         "2222-3434-5656",
					InfobaseID: "3333-4444",
					ProcessID:  "1-1-1-2",
					Host:       "test-ic-1",
					AppID:      "1cv8",
				},
			},
		},
		{
			name:          "Error no command",
			ctx:           context.Background(),
			cs:            "localhost:1545",
			cl:            entity.Cluster{},
			ib:            entity.Infobase{},
			cred:          entity.Credentials{},
			stdout:        NewFakeConnection3(),
			res:           make([]entity.Connection, 0),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Once()

			ctrl := New(pipeMock)

			res, err := ctrl.GetConnections(tc.ctx, tc.cs, tc.cl, tc.ib, tc.cred)

			if err == nil {
				require.NoError(t, err)

				require.NotEmpty(t, res)
				require.ElementsMatch(t, res, tc.res)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, res)
			}
		})
	}
}

func TestDeleteConnection(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		s                 entity.Connection
		clCred            entity.Credentials
		stdout            *FakeReadCloser
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Error connection is empty empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s:             entity.Connection{},
			clCred:        entity.Credentials{},
			stdout:        NewFakeConnection3(),
			respError:     ": connection is empty",
			pipeMockError: ErrConnectionIsEmpty,
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Connection{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{},
			stdout: NewFakeConnection3(),
		},
		{
			name: "Error empty ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Connection{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout:        NewFakeConnection3(),
			respError:     ": connection is empty",
			pipeMockError: ErrConnectionIsEmpty,
		},
		{
			name: "Success w ib w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			s: entity.Connection{
				ID: "3333-4444",
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeConnection3(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			s: entity.Connection{
				ID: "12",
			},
			clCred:        entity.Credentials{},
			stdout:        NewFakeConnection3(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			err := ctrl.DeleteConnection(tc.ctx, tc.cs, tc.cl, tc.s, tc.clCred)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}

func TestDeleteConnections(t *testing.T) {
	cases := []struct {
		name              string
		ctx               context.Context
		cs                string
		cl                entity.Cluster
		ss                []entity.Connection
		clCred            entity.Credentials
		stdout            *FakeReadCloser
		respError         string
		pipeMockError     error
		comMockStartError error
		comMockWaitError  error
	}{
		{
			name: "Success connection is empty wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss:     make([]entity.Connection, 0),
			clCred: entity.Credentials{},
			stdout: NewFakeConnection3(),
		},
		{
			name: "Success w ib wo cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Connection{
				{
					ID: "3333-4444",
				},
				{
					ID: "5555-4444",
				},
			},
			clCred: entity.Credentials{},
			stdout: NewFakeConnection3(),
		},
		{
			name: "Success empty sessions w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Connection{},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeConnection3(),
		},
		{
			name: "Success w sessions w cred",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl: entity.Cluster{
				ID: "1212-3434-5656",
			},
			ss: []entity.Connection{
				{
					ID: "3333-4444",
				},
				{
					ID: "5555-4444",
				},
			},
			clCred: entity.Credentials{
				Name: "test",
				Pwd:  "pwd",
			},
			stdout: NewFakeConnection3(),
		},
		{
			name: "Error no command",
			ctx:  context.Background(),
			cs:   "localhost:1545",
			cl:   entity.Cluster{},
			ss: []entity.Connection{
				{
					ID: "12",
				},
			},
			clCred:        entity.Credentials{},
			stdout:        NewFakeConnection3(),
			respError:     ": no command",
			pipeMockError: errors.New("no command"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			comMock := mocks.NewCommander(t)

			comMock.On("Start").
				Return(tc.comMockStartError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(true) }).
				Maybe()

			comMock.On("Wait").
				Return(tc.comMockWaitError).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			comMock.On("Cancel").
				Return(nil).
				Run(func(args mock.Arguments) { tc.stdout.SetEnable(false) }).
				Maybe()

			pipeMock := mocks.NewPiper(t)

			pipeMock.On("Run",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string")).
				Return(comMock, tc.stdout, tc.pipeMockError).
				Maybe()

			ctrl := New(pipeMock)

			err := ctrl.DeleteConnections(tc.ctx, tc.cs, tc.cl, tc.ss, tc.clCred)

			if err == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)
			}
		})
	}
}
