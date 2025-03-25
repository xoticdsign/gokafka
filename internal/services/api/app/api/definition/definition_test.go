package definition

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	apihttp "gokafka/internal/services/api/http/api"
)

func TestPost_Functional(t *testing.T) {
	cases := []struct {
		name string

		inMethod string
		inTarget string
		inBody   apihttp.PostRequest

		wantBody   apihttp.PostResponse
		wantStatus int
	}{
		{
			name: "Happy Test",

			inMethod: http.MethodPost,
			inTarget: "/post",
			inBody: apihttp.PostRequest{
				Value: "test",
			},

			wantBody: apihttp.PostResponse{
				Value: "post created",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Incorrect Method",

			inMethod: http.MethodGet,
			inTarget: "/post",
			inBody: apihttp.PostRequest{
				Value: "test",
			},

			wantBody: apihttp.PostResponse{
				Value: "method not allowed",
			},
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			api := API{
				Service: &apihttp.UnimplementedService{},
			}

			inBody, _ := json.Marshal(cs.inBody)

			rec := httptest.NewRecorder()
			r := httptest.NewRequest(cs.inMethod, cs.inTarget, bytes.NewBuffer(inBody))

			api.Post(rec, r)

			resp := rec.Result()
			respBody, err := io.ReadAll(resp.Body)
			if assert.NoError(t, err) {
				assert.Equal(t, nil, err)
			}

			var gotBody apihttp.PostResponse

			err = json.Unmarshal(respBody, &gotBody)
			if assert.NoError(t, err) {
				assert.Equal(t, nil, err)
			}

			assert.Equal(t, cs.wantBody, gotBody)
			assert.Equal(t, cs.wantStatus, resp.StatusCode)
		})
	}
}
