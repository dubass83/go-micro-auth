package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/dubass83/go-micro-auth/data/mock"
	data "github.com/dubass83/go-micro-auth/data/sqlc"
	"github.com/dubass83/go-micro-auth/util"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type RoundTripFunc func(rec *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestAuthenticate(t *testing.T) {
	user, password := randomUser()
	clientOk := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"error":false,"massage":"some dummy message"}`)),
			Header:     make(http.Header),
		}
	})
	clientBad := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	testCases := []struct {
		name          string
		body          map[string]interface{}
		buildStubs    func(store *mockdb.MockStore)
		client        *http.Client
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: map[string]interface{}{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
			},
			client: clientOk,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusAccepted, recorder.Code)
			},
		},
		{
			name: "BadParams",
			body: map[string]interface{}{
				"name": 1,
				"pass": 2,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(data.User{}, errors.New("some error"))
			},
			client: clientOk,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FailedToSentLogs",
			body: map[string]interface{}{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetUserByEmail(gomock.Any(), user.Email).
					Times(1).
					Return(user, nil)
			},
			client: clientBad,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {

		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			// start test server and send request
			server := NewTestServer(t, store)
			server.Client = tc.client
			recorder := httptest.NewRecorder()

			url := "/authenticate"
			// convert map to json
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			reqest, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, reqest)
			// check return status code
			tc.checkResponse(t, recorder)
		})

	}
}

func randomUser() (data.User, string) {
	password := util.RandomString(6)
	hash, err := util.HashPassword(password)
	if err != nil {
		return data.User{}, ""
	}
	user := data.User{
		ID:         1,
		Email:      util.RandomEmail(),
		Password:   hash,
		FirstName:  pgtype.Text{String: util.RandomString(6), Valid: true},
		LastName:   pgtype.Text{String: util.RandomString(6), Valid: true},
		UserActive: 1,
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
	return user, password
}
