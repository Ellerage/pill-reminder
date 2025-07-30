package tests

import (
	"os"
	"pill-reminder/internal/i18n"
	testsdb "pill-reminder/tests/testdb"
	"pill-reminder/tests/utils"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	teardownMongo := testsdb.SetupMongo()
	teardownRedis := testsdb.SetupRedis()

	servicesTeardown := utils.Init()
	i18n.Init()

	code := m.Run()

	servicesTeardown()
	time.Sleep(2 * time.Second)

	teardownRedis()
	teardownMongo()

	os.Exit(code)
}
