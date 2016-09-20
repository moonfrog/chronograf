package enterprise_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/influxdata/mrfusion"
	"github.com/influxdata/mrfusion/enterprise"
	"github.com/influxdata/mrfusion/mock"
	"github.com/influxdata/plutonium/meta/control"
)

func Test_Enterprise_FetchesDataNodes(t *testing.T) {
	t.Parallel()

	ctrl := &mock.ControlClient{
		Cluster: &control.Cluster{},
	}
	cl := &enterprise.Client{
		Ctrl: ctrl,
	}

	err := cl.Open()

	if err != nil {
		t.Fatal("Unexpected error while creating enterprise client. err:", err)
	}

	if ctrl.ShowClustersCalled != true {
		t.Fatal("Expected request to meta node but none was issued")
	}
}

func Test_Enterprise_IssuesQueries(t *testing.T) {
	t.Parallel()

	called := false
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		called = true
		if r.URL.Path != "/query" {
			t.Fatal("Expected request to '/query' but was", r.URL.Path)
		}
		rw.Write([]byte(`{}`))
	}))
	defer ts.Close()

	cl := &enterprise.Client{
		Ctrl: mock.NewMockControlClient(ts.URL),
	}

	err := cl.Open()
	if err != nil {
		t.Fatal("Unexpected error initializing client: err:", err)
	}

	_, err = cl.Query(context.Background(), mrfusion.Query{"show shards", "_internal", "autogen"})

	if err != nil {
		t.Fatal("Unexpected error while querying data node: err:", err)
	}

	if called == false {
		t.Fatal("Expected request to data node but none was received")
	}
}