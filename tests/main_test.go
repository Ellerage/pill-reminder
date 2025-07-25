package tests

import (
	"os"
	testsdb "pill-reminder/tests/testdb"
	"testing"
)

func TestMain(m *testing.M) {
	teardownMongo := testsdb.SetupMongo()
	teardownRedis := testsdb.SetupRedis()

	code := m.Run()

	teardownMongo()
	teardownRedis()

	os.Exit(code)
}
