package repository

import (
	"pill-reminder/internal/model"
	"pill-reminder/internal/utils"
	"pill-reminder/internal/utils/enums"
	"pill-reminder/tests/seeds"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_GetAll(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	user1 := seeds.UserSeed(db, nil)
	user2 := seeds.UserSeed(db, nil)

	expectedMap := map[int64]model.User{
		user1.ChatId: user1,
		user2.ChatId: user2,
	}

	users := []model.User{
		user1,
		user2,
	}

	repo := NewUserRepo(db)

	foundUsers, err := repo.GetAll()

	require.NoError(t, err)
	require.Len(t, foundUsers, len(users))

	for _, actual := range foundUsers {
		expected, ok := expectedMap[actual.ChatId]

		assert.True(t, ok)

		require.Equal(t, expected.ChatId, actual.ChatId)
		require.Equal(t, expected.Timezone, actual.Timezone)
		require.Equal(t, expected.TimeToNotify, actual.TimeToNotify)
		require.Equal(t, expected.Status, actual.Status)
	}
}

func TestUserRepo_GetAll_Empty(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()
	repo := NewUserRepo(db)

	users, err := repo.GetAll()
	require.NoError(t, err)
	require.Empty(t, users)
}

func TestUserRepo_GetByChatId(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	fakeUser := seeds.UserSeed(db, nil)

	repo := NewUserRepo(db)

	foundUser, err := repo.GetByChatId(fakeUser.ChatId)

	assert.NoError(t, err)

	require.Equal(t, fakeUser.ChatId, foundUser.ChatId)
	require.Equal(t, fakeUser.Timezone, foundUser.Timezone)
	require.Equal(t, fakeUser.TimeToNotify, foundUser.TimeToNotify)
	require.Equal(t, fakeUser.Status, foundUser.Status)
}

func TestUserRepo_GetByChatId_Not_Found(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	repo := NewUserRepo(db)

	foundUser, err := repo.GetByChatId(gofakeit.Int64())

	assert.Error(t, err)
	require.Equal(t, utils.ErrNotFound, err)
	assert.Nil(t, foundUser)
}

func TestUserRepo_Create(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	repo := NewUserRepo(db)

	toCreate := generateFakeUser()
	err := repo.Create(toCreate)
	assert.NoError(t, err)

	found, err := seeds.GetUserByChatId(t, db, toCreate.ChatId)
	require.NoError(t, err)
	assert.Equal(t, toCreate, *found)
}

func TestUserRepo_Update(t *testing.T) {
	db, teardown := SetupSQLite(t)
	defer teardown()

	fakeUser := seeds.UserSeed(db, nil)

	repo := NewUserRepo(db)

	timezone := gofakeit.TimeZone()
	timeToNotify := gofakeit.Date().Format("15:04")
	status := string(enums.UserStatusIdle)

	toUpdate := model.UserUpdate{
		Timezone:     &timezone,
		TimeToNotify: &timeToNotify,
		Status:       &status,
	}

	updateError := repo.Update(fakeUser.ChatId, toUpdate)

	assert.NoError(t, updateError)

	found, err := seeds.GetUserByChatId(t, db, fakeUser.ChatId)
	assert.NoError(t, err)

	require.Equal(t, *toUpdate.Timezone, found.Timezone)
	require.Equal(t, *toUpdate.TimeToNotify, found.TimeToNotify)
	require.Equal(t, *toUpdate.Status, found.Status)
}

// UTILS
func generateFakeUser() model.User {
	return model.User{
		ChatId:       gofakeit.Int64(),
		Timezone:     gofakeit.TimeZone(),
		TimeToNotify: gofakeit.Date().Format("15:04"),
		Status:       string(enums.UserStatusIdle),
	}
}
