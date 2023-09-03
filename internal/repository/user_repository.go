package repository

import (
	"context"
	"gin-jwt-boilerplate/internal/domain"
	"gin-jwt-boilerplate/internal/mongodb"
	"gin-jwt-boilerplate/internal/utils"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"os"
)

// A UserMongoRepository that implements UserRepository
type UserMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	userRepository *UserMongoRepository
)

// NewUserMongoRepository creates a new UserMongoRepository instance
// which implements UserRepository
func newUserMongoRepository() *UserMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	container, err := mongodb.StartContainer(ctx)
	if err != nil {
		return nil
	}

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		utils.LogFatal("Error getting endpoint: ", err)
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		utils.LogFatal("Error creating mongo client: ", err)
	}

	err = mongoClient.Connect(ctx)
	if err != nil {
		utils.LogFatal("Error connecting to mongo client: ", err)
	}

	err = godotenv.Load()
	if err != nil {
		utils.LogFatal("Error loading .env file: ", err)
	}

	collection := mongoClient.Database(os.Getenv("LYDIA_DB_NAME")).Collection("users")

	return &UserMongoRepository{
		client:     mongoClient,
		collection: collection,
	}
}

// GetUserRepository returns a UserRepository instance if it is not initialized yet or
// returns the existing one
func GetUserRepository() *UserMongoRepository {
	if userRepository == nil {
		userRepository = newUserMongoRepository()
	}
	return userRepository
}

// SaveUser saves a user
func (r *UserMongoRepository) SaveUser(user domain.UserModel) (domain.UserModel, error) {
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

// GetUser gets a user by id
func (r *UserMongoRepository) GetUser(id bson.ObjectId) (domain.UserModel, error) {
	var user domain.UserModel
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}

// ExistsUser checks if a user exists
func (r *UserMongoRepository) ExistsUser(id bson.ObjectId) (bool, error) {
	var user domain.UserModel
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteUser deletes a user by id
func (r *UserMongoRepository) DeleteUser(id bson.ObjectId) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *UserMongoRepository) ExistsByUsername(username string) bool {
	count, err := r.collection.CountDocuments(context.Background(), bson.M{"username": username})
	if err != nil {
		return false
	}
	return count > 0
}

func (r *UserMongoRepository) GetUserByUsername(username string) (domain.UserModel, error) {
	var user domain.UserModel
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return domain.UserModel{}, err
	}
	return user, nil
}
