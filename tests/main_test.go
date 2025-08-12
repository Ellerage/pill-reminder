package tests

import (
	"os"
	"pill-reminder/internal/i18n"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	i18n.Init()

	code := m.Run()

	time.Sleep(2 * time.Second)

	os.Exit(code)
}
