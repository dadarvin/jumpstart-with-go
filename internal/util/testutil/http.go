package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type requestParam struct {
	path   map[string]string
	header map[string]string
	cookie map[string]string
	body   interface{}
}

type RequestOption func(*requestParam)

func NewRequestRecorder(t *testing.T, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params), method, routerPath string, opts ...RequestOption) *httptest.ResponseRecorder {
	var (
		path      = routerPath
		bodyBytes io.Reader
	)

	opt := &requestParam{}

	for _, o := range opts {
		o(opt)
	}

	for k, v := range opt.path {
		path = strings.Replace(path, fmt.Sprintf(":%s", k), v, 1)
	}

	if opt.body != nil {
		marshaled, err := json.Marshal(opt.body)
		if err != nil {
			t.Fatal(err)
		}

		bodyBytes = bytes.NewReader(marshaled)
	}

	req, err := http.NewRequest(method, path, bodyBytes)
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range opt.header {
		req.Header.Add(k, v)
	}

	for k, v := range opt.cookie {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	router := httprouter.New()
	router.Handle(method, routerPath, fnHandler)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func WithBody(param interface{}) RequestOption {
	return func(c *requestParam) {
		c.body = param
	}
}

func WithRequestHeader(param map[string]string) RequestOption {
	return func(c *requestParam) {
		c.header = param
	}
}

func ParseResponse(t *testing.T, rr *httptest.ResponseRecorder, response interface{}) {
	data := rr.Body.Bytes()
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err, rr.Body.String())
	}
}
