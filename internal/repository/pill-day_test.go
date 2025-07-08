package repository

import (
	"context"
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestPillDay_GetByDateAndChatId(t *testing.T) {
	db, teardown := SetupMongo(t)
	defer teardown()

	pillDay := generatePillDay()
	pillDayDate, parseDateErr := time.Parse("2006-01-02", pillDay.Date)
	assert.NoError(t, parseDateErr)

	pillDayColl := db.Collection("pill-day")

	_, err := pillDayColl.InsertOne(context.Background(), pillDay)
	assert.NoError(t, err)

	repo := NewPillDayRepo(db)

	found, getErr := repo.GetByDateAndChatId(pillDay.ChatId, pillDayDate)
	assert.NoError(t, getErr)

	assert.Equal(t, pillDay, *found)
}

func TestPillDay_Create(t *testing.T) {
	db, teardown := SetupMongo(t)
	defer teardown()

	pillDayColl := db.Collection("pill-day")
	repo := NewPillDayRepo(db)

	// toCreate := generatePillDay()

	chatId := gofakeit.Int64()
	timeOfTaking := gofakeit.Date()
	nowDate := utils.GetFormattedNowDate()

	formattedTime := timeOfTaking.Format("15:04")

	expected := model.PillDay{
		ChatId:       chatId,
		TimeOfTaking: &formattedTime,
		Date:         nowDate,
	}

	err := repo.Create(chatId, &timeOfTaking)
	assert.NoError(t, err)

	var found model.PillDay
	err = pillDayColl.FindOne(context.Background(), bson.M{"chatId": chatId, "date": nowDate}).Decode(&found)
	require.NoError(t, err)
	assert.Equal(t, expected, found)
}

func TestPillDayRepo_UpdateTimeByDate(t *testing.T) {
	db, teardown := SetupMongo(t)
	defer teardown()

	pillDayColl := db.Collection("pill-day")

	fakePillDay := generatePillDay()

	now := utils.GetNowDateTime()
	time := gofakeit.Date().Format("15:04")

	_, err := pillDayColl.InsertOne(context.Background(), model.PillDay{
		ChatId:       fakePillDay.ChatId,
		Date:         now.Format("2006-01-02"),
		TimeOfTaking: &time,
	})

	assert.NoError(t, err)

	repo := NewPillDayRepo(db)

	updateError := repo.UpdateTimeByDate(fakePillDay.ChatId, now)
	assert.NoError(t, updateError)

	formattedNowTime := now.Format("15:04")
	expected := model.PillDay{
		ChatId:       fakePillDay.ChatId,
		Date:         now.Format("2006-01-02"),
		TimeOfTaking: &formattedNowTime,
	}

	var found model.PillDay
	err = pillDayColl.FindOne(context.Background(), bson.M{"chatId": fakePillDay.ChatId, "date": now.Format("2006-01-02")}).Decode(&found)
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
