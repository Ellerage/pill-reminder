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
	teardownRedis := testsdb.SetupRedis()
	_, teardownSqlLite := testsdb.SetupSQLite()
	servicesTeardown := utils.Init()
	i18n.Init()

	code := m.Run()

	servicesTeardown()
	time.Sleep(2 * time.Second)

	teardownRedis()
	time.Sleep(2 * time.Second)
	teardownSqlLite()

	os.Exit(code)
}
