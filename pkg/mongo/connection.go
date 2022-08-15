package mongo

import (
	"context"
	"fmt"
	"keeper/internal/config"
	"time"

	_ "keeper/pkg/log"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Connection struct {
	Client *mongo.Client
	Ctx    context.Context
}

func NewConnection(cfg *config.Config) Connection {
	var MONGO_CONN_URI string
	if cfg.Env == "development" {
		MONGO_CONN_URI = fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort)
	} else {
		MONGO_CONN_URI = cfg.MongoDbConnUri
	}

	clientOpts := options.Client().ApplyURI(MONGO_CONN_URI)
	// context: to cancel the connection operation if it times out
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()
	//establish the connection
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		logrus.Panicf("Error connecting to mongo: %v", err)
	}

	// ping the primary to ensure a valid connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.Panicf("Error pinging mongo connection: %v", err)
	}

	logrus.Info("Successfully pinged the Mongo database.")

	return Connection{
		Client: client,
		Ctx:    ctx,
	}
}

func (c *Connection) Disconnect() {
	c.Client.Disconnect(c.Ctx)
}
