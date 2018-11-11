package cmd_test

import (
	// "github.com/stretchr/testify/assert"
	"github.com/tjimsk/todos/cmd"
	"os"
	"path"
	"testing"
)

func TestParseGoFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	fp := path.Join(wd, "/testdata/test.go")

	tags := cmd.ParseGoFile(fp)

	for _, _t := range tags {
		t.Log(_t)
	}
}
