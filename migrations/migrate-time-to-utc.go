package main

import (
	"context"
	"fmt"
	"log/slog"
	"pill-reminder/configs"
	"pill-reminder/internal/db"
	"pill-reminder/internal/model"
	"time"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	cfg := configs.InitConfig()

	fmt.Println(cfg.MONGO_URL)

	mongo := db.ConnectMongo(db.ConnectMongoOptions{
		Uri:    cfg.MONGO_URL,
		DBName: cfg.MONGO_DB_NAME,
	})

	// User settings migration
	cursor, err := mongo.Collection("users").Find(context.TODO(), bson.M{})

	if err != nil {
		fmt.Println("Error", err)
	}

	users := make([]model.User, 0)

	if err := cursor.All(context.TODO(), &users); err != nil {
		fmt.Println("Error", err)
	}

	for _, user := range users {
		loc, errTimezone := time.LoadLocation(user.Timezone)

		if errTimezone != nil {
			fmt.Println(err)
		}

		now := time.Now().In(loc)
		parsed, err := time.Parse("15:04", user.TimeToNotify)

		if err != nil {
			fmt.Println("parse error:", err)
		}

		userTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			parsed.Hour(), parsed.Minute(), 0, 0,
			loc,
		)

		userTimeUTC := userTime.UTC()

		toUpdate := bson.M{"timeToNotify": userTimeUTC.Format("15:04")}

		mongo.Collection("users").UpdateOne(
			context.TODO(),
			bson.M{"chatId": user.ChatId},
			bson.M{
				"$set": toUpdate,
			},
		)
	}

	// Pill day migration
	cursorPillDay, err := mongo.Collection("pill-day").Find(context.TODO(), bson.M{})
	if err != nil {
		slog.Error(err.Error())
	}

	pillDays := make([]model.PillDay, 0)

	err = cursorPillDay.All(context.TODO(), &pillDays)
	if err != nil {
		slog.Error(err.Error())
	}

	slog.Info(fmt.Sprintf("Pill days length: %d", len(pillDays)))

	for _, pillDay := range pillDays {
		loc, errTimezone := time.LoadLocation("Asia/Tbilisi")

		if errTimezone != nil {
			fmt.Println(err)
		}

		parsedDate, err := time.Parse("2006-01-02", pillDay.Date)
		if err != nil {
			slog.Error(err.Error())
		}

		parsedTime, err := time.Parse("15:04", *pillDay.TimeOfTaking)
		if err != nil {
			fmt.Println("parse error:", err)
		}

		userTime := time.Date(
			parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
			parsedTime.Hour(), parsedTime.Minute(), 0, 0,
			loc,
		)

		userTimeUTC := userTime.UTC()

		toUpdate := bson.M{"timeOfTaking": userTimeUTC.Format("15:04")}

		mongo.Collection("pill-day").UpdateOne(
			context.TODO(),
			bson.M{"date": pillDay.Date},
			bson.M{
				"$set": toUpdate,
			},
		)

	}

}
