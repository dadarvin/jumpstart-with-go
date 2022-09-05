package middleware

import (
	"entry_task/internal/util/httputil"
	"entry_task/internal/util/testutil"
	"entry_task/pkg/dto/base"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func testFunc() httprouter.Handle {
	return func(w http.ResponseWriter, _ *http.Request, param httprouter.Params) {
		httputil.HttpResponse(w, http.StatusOK, "Sukses", nil)
	}
}

func TestMiddleware_IsAuthorized(t *testing.T) {
	type args struct {
		secret string
		header map[string]string
	}
	tests := []struct {
		name string
		args args
		want func(t *testing.T, result *base.JsonResponse, status int)
	}{
		{
			name: "Token Not Valid",
			args: args{
				secret: "secret",
				header: map[string]string{
					"Token": "th1s154t0k3n",
				},
			},
			want: func(t *testing.T, result *base.JsonResponse, status int) {
				assert.Equal(t, http.StatusUnauthorized, status)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Middleware{
				secret: tt.args.secret,
			}

			rr := testutil.NewRequestRecorder(t,
				m.IsAuthorized(testFunc()), http.MethodPost,
				"/", testutil.WithRequestHeader(tt.args.header),
			)

			var resp base.JsonResponse

			testutil.ParseResponse(t, rr, &resp)
			tt.want(t, &resp, rr.Code)
		})
	}
}
