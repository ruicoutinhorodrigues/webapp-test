package main

import (
	"os"
	"testing"

	"local/rui.rodrigues/webapp-test/pkg/repository/dbrepo"
)

var (
	app          application
	expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE2NzYxMjQwNzUsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkpvaG4gRG9lIiwic3ViIjoiMSJ9.2MBfW4JgumnphtxpfVud6-jLV3xC-zd19rGfJJDOVy0"
)

func TestMain(m *testing.M) {
	app.DB = &dbrepo.TestDBRepo{}
	app.Domain = "example.com"
	app.JWTSecret = "great_jwt_secret"

	os.Exit(m.Run())
}
