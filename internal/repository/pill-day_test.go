package repository

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/tests/seeds"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPillDay_GetByDateAndChatId(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	pillDay := generatePillDay()
	pillDayDate, parseDateErr := time.Parse("2006-01-02", pillDay.Date)
	assert.NoError(t, parseDateErr)

	db.Exec("INSERT INTO pillDays (date, timeOfTaking, chatId) VALUES (?, ?, ?)", pillDay.Date, pillDay.TimeOfTaking, pillDay.ChatId)

	repo := NewPillDayRepo(db)

	found, getErr := repo.GetByDateAndChatId(pillDay.ChatId, pillDayDate)
	assert.NoError(t, getErr)

	assert.Equal(t, pillDay, *found)
}

func TestPillDay_Create(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	repo := NewPillDayRepo(db)

	chatId := gofakeit.Int64()
	nowDate := utils.GetFormattedNowDate()

	timeOfTaking := gofakeit.Date()
	formattedTime := timeOfTaking.Format("15:04")

	expected := model.PillDay{
		ChatId:       chatId,
		TimeOfTaking: &formattedTime,
		Date:         nowDate,
	}

	err := repo.Create(chatId, &timeOfTaking)
	if err != nil {
		t.Fatal(err)
	}

	result, err := seeds.FindPillDayByChatId(t, db, chatId)

	if err != nil {
		t.Fatal(err)
	}

	require.NoError(t, err)
	assert.Equal(t, expected, *result)
}

func TestPillDayRepo_UpdateTimeByDate(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	fakePillDay := generatePillDay()

	now := utils.GetNowDateTime()
	time := gofakeit.Date().Format("15:04")

	repo := NewPillDayRepo(db)

	date := now.Format("2006-01-02")
	seeds.PillDaySeed(db, &seeds.PillDayParams{Date: &date, TimeOfTaking: &time, ChatId: &fakePillDay.ChatId})

	updateError := repo.UpdateTimeByDate(fakePillDay.ChatId, now)
	assert.NoError(t, updateError)

	formattedNowTime := now.Format("15:04")
	expected := model.PillDay{
		ChatId:       fakePillDay.ChatId,
		Date:         now.Format("2006-01-02"),
		TimeOfTaking: &formattedNowTime,
	}

	found, err := seeds.FindPillDayByChatId(t, db, fakePillDay.ChatId)
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)

	assert.Equal(t, expected.ChatId, found.ChatId)
	assert.Equal(t, expected.Date, found.Date)
	assert.Equal(t, expected.TimeOfTaking, found.TimeOfTaking)
}

// UTILS
func generatePillDay() model.PillDay {
	timeOfTaking := gofakeit.Date().Format("15:04")

	return model.PillDay{
		ChatId:       gofakeit.Int64(),
		Date:         gofakeit.Date().Format("2006-01-02"),
		TimeOfTaking: &timeOfTaking,
	}
}
