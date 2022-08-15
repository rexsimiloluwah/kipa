package repository

import (
	"context"
	"errors"
	"fmt"
	"keeper/internal/config"
	"keeper/internal/models"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	userCollectionName = "users"
)

var userDetailsProjection = bson.D{
	primitive.E{Key: "_id", Value: 1},
	primitive.E{Key: "firstname", Value: 1},
	primitive.E{Key: "lastname", Value: 1},
	primitive.E{Key: "email", Value: 1},
	primitive.E{Key: "role", Value: 1},
	primitive.E{Key: "password", Value: 1},
	primitive.E{Key: "email_verified", Value: 1},
	primitive.E{Key: "created_at", Value: 1},
	primitive.E{Key: "updated_at", Value: 1},
}

type IUserRepository interface {
	CreateUser(user *models.User) error
	FindUserByEmail(email string) (*models.User, error)
	FindUserById(id string) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userId string) error
}

type UserRepository struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserRepository(cfg *config.Config, dbClient *mongo.Client) IUserRepository {
	userCollection := dbClient.Database(cfg.DbName).Collection(userCollectionName)
	return &UserRepository{
		collection: userCollection,
		ctx:        context.TODO(),
	}
}

// Create a new user
func (r *UserRepository) CreateUser(user *models.User) error {
	_, err := r.collection.InsertOne(r.ctx, user)
	if err != nil {
		logrus.Errorf("error creating user: %s", err.Error())
		return fmt.Errorf("error creating user: %s", err.Error())
	}
	return nil
}

// Find a user by email
func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	opts := options.FindOne().SetProjection(userDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	logrus.Info("found user: ", user)
	return user, nil
}

// Find a user by id
func (r *UserRepository) FindUserById(id string) (*models.User, error) {
	user := &models.User{}
	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, models.ErrInvalidObjectID
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	opts := options.FindOne().SetProjection(userDetailsProjection)
	if err := r.collection.FindOne(r.ctx, filter, opts).Decode(user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}
	logrus.Info("found user: ", user)
	return user, nil
}

// Returns all the users
func (r *UserRepository) FindAllUsers() ([]models.User, error) {
	users := []models.User{}
	filter := bson.D{}

	opts := options.Find().SetProjection(userDetailsProjection)
	cursor, err := r.collection.Find(r.ctx, filter, opts)
	if err != nil {
		logrus.WithError(err).Errorf("cannot find many")
		return nil, models.ErrUsersNotFound
	}
	if err = cursor.All(r.ctx, &users); err != nil {
		return nil, models.ErrUsersNotFound
	}
	logrus.Debug("users: ", users)
	return users, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	filter := bson.D{primitive.E{Key: "_id", Value: user.ID}}
	uByte, err := bson.Marshal(user)
	if err != nil {
		return errors.New("error marshalling user")
	}
	update := make(map[string]interface{}, 0)
	if err = bson.Unmarshal(uByte, update); err != nil {
		return errors.New("error unmarshalling update data")
	}

	opts := options.Update().SetUpsert(true)
	result, err := r.collection.UpdateOne(r.ctx, filter, bson.D{primitive.E{Key: "$set", Value: update}}, opts)
	if err != nil {
		logrus.WithError(err).Error("error updating user")
		return err
	}
	if result.UpsertedCount == 0 {
		return models.ErrUpdatingUser
	}
	return nil
}

func (r *UserRepository) DeleteUser(userId string) error {
	ID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: ID}}
	result, err := r.collection.DeleteOne(r.ctx, filter)
	if err != nil {
		logrus.WithError(err).Error("error deleting user")
		return models.ErrDeletingUser
	}
	if result.DeletedCount == 0 {
		return models.ErrDeletingUser
	}
	return nil
}
