package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"path/filepath"
	"sber/pkg/settings"
	"sber/types"
)

func NewMongoDB(ctx context.Context) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/", settings.Config.DBUsername, os.Getenv("DB_PASSWORD"), settings.Config.DBHost, settings.Config.DBPort))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	// todo clear table users
	loadDbData(ctx, client)

	return client, nil
}

func loadDbData(ctx context.Context, client *mongo.Client) {
	_, err := client.Database(DbEmployees).Collection(DbTableUsers).DeleteMany(ctx, bson.M{})
	if err != nil {
		logrus.Fatal(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	byteValues, err := ioutil.ReadFile(filepath.Join(pwd, "scripts/db employees data.json"))
	if err != nil {
		logrus.Fatal(err)
	}

	var users []types.DBUser
	err = json.Unmarshal(byteValues, &users)
	if err != nil {
		logrus.Fatal(err)
	}

	for i := range users {
		user := users[i]
		_, err = client.Database(DbEmployees).Collection(DbTableUsers).InsertOne(ctx, user)
		if err != nil {
			logrus.Fatal(err)
		}
	}
	return
}
