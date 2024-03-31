package delete_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/RomanLevBy/UrlShortener/internal/http-server/handlers/url/delete"
	"github.com/RomanLevBy/UrlShortener/internal/http-server/handlers/url/delete/mocks"
	"github.com/RomanLevBy/UrlShortener/internal/lib/api/response"
	"github.com/RomanLevBy/UrlShortener/internal/lib/logger/handlers/slogdiscard"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test-success-alias",
		},
		{
			name:      "Empty alias",
			alias:     "",
			respError: "Alias is empty",
		},
		{
			name:      "Delete Repo Error",
			alias:     "failed-repo-error",
			respError: "Fail to delete url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlRemoverMock := mocks.NewURLRemover(t)

			if tc.respError == "" || tc.mockError != nil {
				urlRemoverMock.On("DeleteURL", tc.alias, mock.AnythingOfType("string")).
					Return(tc.mockError).
					Once()
			}

			handler := delete.New(slogdiscard.NewDiscardLogger(), urlRemoverMock)
			req, err := http.NewRequest(http.MethodDelete, "/url/"+tc.alias, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.alias)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp response.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))



			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
