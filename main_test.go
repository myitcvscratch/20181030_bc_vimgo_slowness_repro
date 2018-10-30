// +build !race

package main

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/rogpeppe/go-internal/goproxytest"
	"github.com/rogpeppe/go-internal/gotooltest"
	"github.com/rogpeppe/go-internal/testscript"
)

var (
	proxyURL string
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(gobinMain{m}, map[string]func() int{
		"gopherjs": main1,
	}))
}

type gobinMain struct {
	m *testing.M
}

func (m gobinMain) Run() int {
	// Start the Go proxy server running for all tests.
	srv, err := goproxytest.NewServer("testdata/mod", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot start proxy: %v", err)
		return 1
	}
	proxyURL = srv.URL

	return m.m.Run()
}

func TestScripts(t *testing.T) {
	p := testscript.Params{
		Dir: "testdata",
		Cmds: {
			"modified": func(ts *testscript.TestScript, neg bool, args []string) {

			},
		},
		Setup: func(e *testscript.Env) error {
			e.Vars = append(e.Vars,
				"NODE_PATH="+os.Getenv("NODE_PATH"),
				"GOPROXY="+proxyURL,
			)
			return nil
		},
	}
	if err := gotooltest.Setup(&p); err != nil {
		t.Fatal(err)
	}
	testscript.Run(t, p)
}

func buildModified() func(ts *testscript.TestScript, neg bool, args []string) {
	modCacheLock := new(sync.Mutex)
	modCache := make(map[string]func(ts *testscript.TestScript, neg bool, args []string))

	return func(ts *testscript.TestScript, neg bool, args []string) {
		modCacheLock.Lock()
		mod // put cursor after mod
		defer modCacheLock.Unlock()
	}
}
