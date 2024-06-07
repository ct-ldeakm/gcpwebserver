package gcpwebserv

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestPaths(t *testing.T) {

	type test struct {
		host   string
		path   string
		method string
		ids    []string
		values []string
	}

	tests := []test{
		{host: "http://localhost:8080", path: "/1", method: "GET", ids: []string{"id"}, values: []string{"1"}},
		{host: "http://localhost:8080", path: "/1/2", method: "GET", ids: []string{"id", "region"}, values: []string{"1", "2"}},
		{host: "http://localhost:8080", path: "/static/bob.txt", method: "GET", ids: []string{"id", "region"}, values: []string{"1", "2"}},
	}

	server := Setup(nil)

	ctx, endTest := context.WithCancel(context.Background())

	go func() {
		err := Run(ctx, server)
		if err != nil {
			t.Fatalf("Server Run Error:%s", err)
		}
	}()

	for _, test := range tests {

		routePath := fmt.Sprintf("%s /{%s}", test.method, strings.Join(test.ids, "}/{"))
		Route(routePath, func(w http.ResponseWriter, r *http.Request) {
			var rescol []string
			for _, id := range test.ids {
				rescol = append(rescol, r.PathValue(id))
			}
			fmt.Fprintf(w, strings.Join(rescol, ","))

		})
	}

	for _, test := range tests {
		url := fmt.Sprintf("%s%s", test.host, test.path)
		resp, err := http.Get(url)
		if err != nil {
			t.Logf("Http Request Error:%s", err)
		}
		bts, err := io.ReadAll(resp.Body)
		result := strings.Split(string(bts), ",")
		if !reflect.DeepEqual(result, test.values) {
			t.Fatalf("Have :%v want : %v", result, test.values)
		}

	}
	endTest()
}
