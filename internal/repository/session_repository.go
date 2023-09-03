package repository

import (
	"context"
	"gin-jwt-boilerplate/internal/domain"
	"gin-jwt-boilerplate/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type SessionMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var (
	sessionRepository *SessionMongoRepository
)

func newSessionMongoRepository() *SessionMongoRepository {
	ctx := context.Background()
	// FIXME: Burada ileride uzaktaki bir mongodb instance'ına bağlanmak gerekecek
	// FIXME: Ortaklaştırılacak
	container, err := mongodb.StartContainer(ctx)
	if err != nil {
		return nil
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil
	}

	port, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+host+":"+port.Port()))
	if err != nil {
		return nil
	}

	collection := client.Database("lydia").Collection("sessions")

	return &SessionMongoRepository{
		client:     client,
		collection: collection,
	}
}

func GetSessionRepository() *SessionMongoRepository {
	if sessionRepository == nil {
		sessionRepository = newSessionMongoRepository()
	}
	return sessionRepository
}

// SaveSession is a function that creates a session
func (s SessionMongoRepository) SaveSession(model domain.SessionInfoModel) (domain.SessionInfoModel, error) {
	_, err := s.collection.InsertOne(context.Background(), model)
	if err != nil {
		return domain.SessionInfoModel{}, err
	}
	return model, nil
}

// GetUserSession is a function that gets a user session
func (s SessionMongoRepository) GetUserSession(id bson.ObjectId) (domain.SessionInfoModel, error) {
	var session domain.SessionInfoModel
	err := s.collection.FindOne(context.Background(), bson.M{"userId": id}).Decode(&session)
	if err != nil {
		return domain.SessionInfoModel{}, err
	}
	return session, nil
}

// DeleteSessionByUserId is a function that deletes all sessions of a user
func (s SessionMongoRepository) DeleteSessionByUserId(userId bson.ObjectId) error {
	_, err := s.collection.DeleteMany(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return err
	}
	return nil
}

// DeleteSessionById is a function that deletes a session by id
func (s SessionMongoRepository) DeleteSessionById(sessionId bson.ObjectId) error {
	_, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": sessionId})
	if err != nil {
		return err
	}
	return nil
}
