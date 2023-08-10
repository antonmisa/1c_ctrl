package v1

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/antonmisa/1cctl/internal/entity"
	ucm "github.com/antonmisa/1cctl/internal/usecase/mocks"
	lm "github.com/antonmisa/1cctl/pkg/logger/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClustersRoute(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		uri            string
		ctrlMockResult []entity.Cluster
		ctrlMockError  error
		code           int
		retVal         string
	}{
		{
			name:   "Erorr wo entrypoint",
			method: http.MethodGet,
			uri:    "/v1/cluster/list",
			ctrlMockResult: []entity.Cluster{
				{
					ID:   "123",
					Name: "test",
				},
			},
			code:   http.StatusBadRequest,
			retVal: "{\"error\":\"bad request\"}",
		},
		{
			name:   "Success",
			method: http.MethodGet,
			uri:    "/v1/cluster/list?entrypoint=1capp01:1545",
			ctrlMockResult: []entity.Cluster{
				{
					ID:   "123",
					Name: "test",
				},
			},
			code:   200,
			retVal: "{\"clusters\":[{\"id\":\"123\",\"host\":\"\",\"port\":\"\",\"name\":\"test\",\"exp\":0,\"lt\":0,\"mms\":0,\"mmts\":0,\"sl\":0,\"sftl\":0,\"lb\":\"\",\"errth\":0,\"kpp\":0}]}",
		},
		{
			name:   "Error entrypoint incorrect",
			method: http.MethodGet,
			uri:    "/v1/cluster/list?entrypoint=unknown:1545",
			ctrlMockResult: []entity.Cluster{
				{
					ID:   "123",
					Name: "test",
				},
			},
			ctrlMockError: errors.New("error entrypoint incorrect"),
			code:          500,
			retVal:        "{\"error\":\"internal problems\"}",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logMock := lm.NewInterface(t)

			logMock.On("Info",
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			logMock.On("Error",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			ctrlMock := ucm.NewCtrl(t)

			ctrlMock.On("Clusters",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.Anything).
				Return(tc.ctrlMockResult, tc.ctrlMockError).
				Maybe()

			handler := gin.New()
			NewRouter(handler, logMock, ctrlMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.uri, nil)
			handler.ServeHTTP(w, req)

			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.retVal, w.Body.String())
		})
	}
}

func TestInfobasesRoute(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		uri            string
		ctrlMockResult []entity.Infobase
		ctrlMockError  error
		code           int
		retVal         string
	}{
		{
			name:   "Erorr wo entrypoint",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/list",
			ctrlMockResult: []entity.Infobase{
				{
					ID:   "123",
					Name: "test",
				},
			},
			code:   http.StatusBadRequest,
			retVal: "{\"error\":\"bad request\"}",
		},
		{
			name:   "Success",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/list?entrypoint=1capp01:1545",
			ctrlMockResult: []entity.Infobase{
				{
					ID:   "123",
					Name: "test",
				},
				{
					ID:   "1234",
					Name: "test1",
				},
			},
			code:   200,
			retVal: "{\"infobases\":[{\"id\":\"123\",\"name\":\"test\",\"desc\":\"\"},{\"id\":\"1234\",\"name\":\"test1\",\"desc\":\"\"}]}",
		},
		{
			name:   "Error entrypoint incorrect",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/list?entrypoint=unknown:1545",
			ctrlMockResult: []entity.Infobase{
				{
					ID:   "123",
					Name: "test",
				},
			},
			ctrlMockError: errors.New("error entrypoint incorrect"),
			code:          500,
			retVal:        "{\"error\":\"internal problems\"}",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logMock := lm.NewInterface(t)

			logMock.On("Info",
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			logMock.On("Error",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			ctrlMock := ucm.NewCtrl(t)

			ctrlMock.On("Infobases",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Return(tc.ctrlMockResult, tc.ctrlMockError).
				Maybe()

			handler := gin.New()
			NewRouter(handler, logMock, ctrlMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.uri, nil)
			handler.ServeHTTP(w, req)

			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.retVal, w.Body.String())
		})
	}
}

func TestSessionsRoute(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		uri            string
		ctrlMockResult []entity.Session
		ctrlMockError  error
		code           int
		retVal         string
	}{
		{
			name:   "Erorr wo entrypoint",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/session/list",
			ctrlMockResult: []entity.Session{
				{
					ID:   "123",
					Host: "test",
				},
			},
			code:   http.StatusBadRequest,
			retVal: "{\"error\":\"bad request\"}",
		},
		{
			name:   "Success",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/session/list?entrypoint=1capp01:1545",
			ctrlMockResult: []entity.Session{
				{
					ID:   "123",
					Host: "test",
				},
				{
					ID:   "1234",
					Host: "test1",
				},
			},
			code:   200,
			retVal: "{\"sessions\":[{\"id\":\"123\",\"sid\":0,\"ib\":\"\",\"conn\":\"\",\"proc\":\"\",\"uname\":\"\",\"host\":\"test\",\"appid\":\"\",\"loc\":\"\",\"started\":\"0001-01-01T00:00:00Z\",\"active\":\"0001-01-01T00:00:00Z\",\"hib\":\"\",\"hibtm\":0,\"hibterm\":0,\"blockdb\":0,\"blockls\":0,\"bytes\":0,\"bytes5m\":0,\"calls\":0,\"calls5m\":0,\"bytesdb\":0,\"bytesdb5m\":0,\"dbproci\":\"\",\"dbproc\":0,\"dbprocat\":\"\",\"dur\":0,\"durdb\":0,\"durcur\":0,\"durcurdb\":0,\"dur5m\":0,\"durdb5m\":0,\"memcur\":0,\"mem5m\":0,\"mem\":0,\"readcur\":0,\"read5m\":0,\"read\":0,\"writecur\":0,\"write5m\":0,\"write\":0,\"dursvccur\":0,\"dursvc5m\":0,\"dursvc\":0,\"svc\":\"\",\"cpucur\":0,\"cpu5m\":0,\"cpu\":0,\"sep\":\"\"},{\"id\":\"1234\",\"sid\":0,\"ib\":\"\",\"conn\":\"\",\"proc\":\"\",\"uname\":\"\",\"host\":\"test1\",\"appid\":\"\",\"loc\":\"\",\"started\":\"0001-01-01T00:00:00Z\",\"active\":\"0001-01-01T00:00:00Z\",\"hib\":\"\",\"hibtm\":0,\"hibterm\":0,\"blockdb\":0,\"blockls\":0,\"bytes\":0,\"bytes5m\":0,\"calls\":0,\"calls5m\":0,\"bytesdb\":0,\"bytesdb5m\":0,\"dbproci\":\"\",\"dbproc\":0,\"dbprocat\":\"\",\"dur\":0,\"durdb\":0,\"durcur\":0,\"durcurdb\":0,\"dur5m\":0,\"durdb5m\":0,\"memcur\":0,\"mem5m\":0,\"mem\":0,\"readcur\":0,\"read5m\":0,\"read\":0,\"writecur\":0,\"write5m\":0,\"write\":0,\"dursvccur\":0,\"dursvc5m\":0,\"dursvc\":0,\"svc\":\"\",\"cpucur\":0,\"cpu5m\":0,\"cpu\":0,\"sep\":\"\"}]}",
		},
		{
			name:   "Error entrypoint incorrect",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/session/list?entrypoint=unknown:1545",
			ctrlMockResult: []entity.Session{
				{
					ID:   "123",
					Host: "test",
				},
			},
			ctrlMockError: errors.New("error entrypoint incorrect"),
			code:          500,
			retVal:        "{\"error\":\"internal problems\"}",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logMock := lm.NewInterface(t)

			logMock.On("Info",
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			logMock.On("Error",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			ctrlMock := ucm.NewCtrl(t)

			ctrlMock.On("Sessions",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Return(tc.ctrlMockResult, tc.ctrlMockError).
				Maybe()

			handler := gin.New()
			NewRouter(handler, logMock, ctrlMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.uri, nil)
			handler.ServeHTTP(w, req)

			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.retVal, w.Body.String())
		})
	}
}

func TestSessionsByInfobaseRoute(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		uri            string
		ctrlMockResult []entity.Session
		ctrlMockError  error
		code           int
		retVal         string
	}{
		{
			name:   "Erorr wo entrypoint",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/1234-5678/session/list",
			ctrlMockResult: []entity.Session{
				{
					ID:   "123",
					Host: "test",
				},
			},
			code:   http.StatusBadRequest,
			retVal: "{\"error\":\"bad request\"}",
		},
		{
			name:   "Success",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/1234-5678/session/list?entrypoint=1capp01:1545",
			ctrlMockResult: []entity.Session{
				{
					ID:   "123",
					Host: "test",
				},
				{
					ID:   "1234",
					Host: "test1",
				},
			},
			code:   200,
			retVal: "{\"sessions\":[{\"id\":\"123\",\"sid\":0,\"ib\":\"\",\"conn\":\"\",\"proc\":\"\",\"uname\":\"\",\"host\":\"test\",\"appid\":\"\",\"loc\":\"\",\"started\":\"0001-01-01T00:00:00Z\",\"active\":\"0001-01-01T00:00:00Z\",\"hib\":\"\",\"hibtm\":0,\"hibterm\":0,\"blockdb\":0,\"blockls\":0,\"bytes\":0,\"bytes5m\":0,\"calls\":0,\"calls5m\":0,\"bytesdb\":0,\"bytesdb5m\":0,\"dbproci\":\"\",\"dbproc\":0,\"dbprocat\":\"\",\"dur\":0,\"durdb\":0,\"durcur\":0,\"durcurdb\":0,\"dur5m\":0,\"durdb5m\":0,\"memcur\":0,\"mem5m\":0,\"mem\":0,\"readcur\":0,\"read5m\":0,\"read\":0,\"writecur\":0,\"write5m\":0,\"write\":0,\"dursvccur\":0,\"dursvc5m\":0,\"dursvc\":0,\"svc\":\"\",\"cpucur\":0,\"cpu5m\":0,\"cpu\":0,\"sep\":\"\"},{\"id\":\"1234\",\"sid\":0,\"ib\":\"\",\"conn\":\"\",\"proc\":\"\",\"uname\":\"\",\"host\":\"test1\",\"appid\":\"\",\"loc\":\"\",\"started\":\"0001-01-01T00:00:00Z\",\"active\":\"0001-01-01T00:00:00Z\",\"hib\":\"\",\"hibtm\":0,\"hibterm\":0,\"blockdb\":0,\"blockls\":0,\"bytes\":0,\"bytes5m\":0,\"calls\":0,\"calls5m\":0,\"bytesdb\":0,\"bytesdb5m\":0,\"dbproci\":\"\",\"dbproc\":0,\"dbprocat\":\"\",\"dur\":0,\"durdb\":0,\"durcur\":0,\"durcurdb\":0,\"dur5m\":0,\"durdb5m\":0,\"memcur\":0,\"mem5m\":0,\"mem\":0,\"readcur\":0,\"read5m\":0,\"read\":0,\"writecur\":0,\"write5m\":0,\"write\":0,\"dursvccur\":0,\"dursvc5m\":0,\"dursvc\":0,\"svc\":\"\",\"cpucur\":0,\"cpu5m\":0,\"cpu\":0,\"sep\":\"\"}]}",
		},
		{
			name:   "Error entrypoint incorrect",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/1234-5678/session/list?entrypoint=unknown:1545",
			ctrlMockResult: []entity.Session{
				{
					ID:   "123",
					Host: "test",
				},
			},
			ctrlMockError: errors.New("error entrypoint incorrect"),
			code:          500,
			retVal:        "{\"error\":\"internal problems\"}",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logMock := lm.NewInterface(t)

			logMock.On("Info",
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			logMock.On("Error",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			ctrlMock := ucm.NewCtrl(t)

			ctrlMock.On("Sessions",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Return(tc.ctrlMockResult, tc.ctrlMockError).
				Maybe()

			handler := gin.New()
			NewRouter(handler, logMock, ctrlMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.uri, nil)
			handler.ServeHTTP(w, req)

			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.retVal, w.Body.String())
		})
	}
}

func TestConnectionsRoute(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		uri            string
		ctrlMockResult []entity.Connection
		ctrlMockError  error
		code           int
		retVal         string
	}{
		{
			name:   "Erorr wo entrypoint",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/connection/list",
			ctrlMockResult: []entity.Connection{
				{
					ID:   "123",
					Host: "test",
				},
			},
			code:   http.StatusBadRequest,
			retVal: "{\"error\":\"bad request\"}",
		},
		{
			name:   "Success",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/connection/list?entrypoint=1capp01:1545",
			ctrlMockResult: []entity.Connection{
				{
					ID:   "123",
					Host: "test",
				},
				{
					ID:   "1234",
					Host: "test1",
				},
			},
			code:   200,
			retVal: "{\"connections\":[{\"id\":\"123\",\"cid\":0,\"ib\":\"\",\"proc\":\"\",\"host\":\"test\",\"appid\":\"\",\"connected\":\"0001-01-01T00:00:00Z\",\"sid\":0,\"blocked\":0},{\"id\":\"1234\",\"cid\":0,\"ib\":\"\",\"proc\":\"\",\"host\":\"test1\",\"appid\":\"\",\"connected\":\"0001-01-01T00:00:00Z\",\"sid\":0,\"blocked\":0}]}",
		},
		{
			name:   "Error entrypoint incorrect",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/connection/list?entrypoint=unknown:1545",
			ctrlMockResult: []entity.Connection{
				{
					ID:   "123",
					Host: "test",
				},
			},
			ctrlMockError: errors.New("error entrypoint incorrect"),
			code:          500,
			retVal:        "{\"error\":\"internal problems\"}",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logMock := lm.NewInterface(t)

			logMock.On("Info",
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			logMock.On("Error",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			ctrlMock := ucm.NewCtrl(t)

			ctrlMock.On("Connections",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Return(tc.ctrlMockResult, tc.ctrlMockError).
				Maybe()

			handler := gin.New()
			NewRouter(handler, logMock, ctrlMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.uri, nil)
			handler.ServeHTTP(w, req)

			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.retVal, w.Body.String())
		})
	}
}

func TestConnectionsByInfobaseRoute(t *testing.T) {
	cases := []struct {
		name           string
		method         string
		uri            string
		ctrlMockResult []entity.Connection
		ctrlMockError  error
		code           int
		retVal         string
	}{
		{
			name:   "Erorr wo entrypoint",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/1234-5678/connection/list",
			ctrlMockResult: []entity.Connection{
				{
					ID:   "123",
					Host: "test",
				},
			},
			code:   http.StatusBadRequest,
			retVal: "{\"error\":\"bad request\"}",
		},
		{
			name:   "Success",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/1234-5678/connection/list?entrypoint=1capp01:1545",
			ctrlMockResult: []entity.Connection{
				{
					ID:   "123",
					Host: "test",
				},
				{
					ID:   "1234",
					Host: "test1",
				},
			},
			code:   200,
			retVal: "{\"connections\":[{\"id\":\"123\",\"cid\":0,\"ib\":\"\",\"proc\":\"\",\"host\":\"test\",\"appid\":\"\",\"connected\":\"0001-01-01T00:00:00Z\",\"sid\":0,\"blocked\":0},{\"id\":\"1234\",\"cid\":0,\"ib\":\"\",\"proc\":\"\",\"host\":\"test1\",\"appid\":\"\",\"connected\":\"0001-01-01T00:00:00Z\",\"sid\":0,\"blocked\":0}]}",
		},
		{
			name:   "Error entrypoint incorrect",
			method: http.MethodGet,
			uri:    "/v1/cluster/1capp01:1541/infobase/1234-5678/connection/list?entrypoint=unknown:1545",
			ctrlMockResult: []entity.Connection{
				{
					ID:   "123",
					Host: "test",
				},
			},
			ctrlMockError: errors.New("error entrypoint incorrect"),
			code:          500,
			retVal:        "{\"error\":\"internal problems\"}",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logMock := lm.NewInterface(t)

			logMock.On("Info",
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			logMock.On("Error",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Maybe()

			ctrlMock := ucm.NewCtrl(t)

			ctrlMock.On("Connections",
				mock.MatchedBy(func(ctx context.Context) bool { return true }),
				mock.AnythingOfType("string"),
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything).
				Return(tc.ctrlMockResult, tc.ctrlMockError).
				Maybe()

			handler := gin.New()
			NewRouter(handler, logMock, ctrlMock)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.uri, nil)
			handler.ServeHTTP(w, req)

			require.Equal(t, tc.code, w.Code)
			require.Equal(t, tc.retVal, w.Body.String())
		})
	}
}
