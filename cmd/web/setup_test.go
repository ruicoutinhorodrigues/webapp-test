package main

import (
	"os"
	"testing"

	"local/rui.rodrigues/webapp-test/pkg/repository/dbrepo"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"

	uploadPath = "./testdata/uploads"

	app.Session = getSession()

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())
}
