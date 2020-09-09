package database

import (
	"context"
	"errors"

	"github.com/igorhalfeld/madalena-backend/lib/model"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/igorhalfeld/madalena-backend/lib/logger"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	errMongoConnection = errors.New("database: não foi possível realizar a conexão com o banco")
	errMongoCreate     = errors.New("database: não foi possível criar o documento")
	errMongoCreateMany = errors.New("database: não foi possível criar os documentos")
	errMongoReadOne    = errors.New("database: não foi possível ler o documento")
	errMongoReadAll    = errors.New("database: não foi possível buscar os documentos")
)

// IMongo interface of Database
type IMongo interface {
	Create(ctx context.Context, collection string, document interface{}) (interface{}, error)
	CreateMany(ctx context.Context, collection string, documents []interface{}) (interface{}, error)
	Update(ctx context.Context, collection string, id string, document interface{}) (interface{}, error)
	ReadOne(ctx context.Context, collection string, id string, model interface{}) error
	ReadAll(ctx context.Context, collection string, search *model.SearchRequest, documents interface{}) error
}

// MongoConfig model of configuration
type MongoConfig struct {
	URL  string
	Name string
}

const mongoDefaultURL = "mongodb://localhost:27017"

type mongoImpl struct {
	conn *mongo.Database
}

// NewMongoConnection function return a new instance
func NewMongoConnection(config MongoConfig) IMongo {
	url := config.URL
	name := config.Name

	if url == "" {
		url = mongoDefaultURL
	}

	opts := options.Client().ApplyURI(url)
	conn, err := mongo.NewClient(opts)
	if err != nil {
		log.Error(errMongoConnection, err.Error())
		panic(errMongoConnection)
	}

	if err := conn.Connect(context.Background()); err != nil {
		log.Error(errMongoConnection)
		panic(errMongoConnection)
	}

	return &mongoImpl{
		conn: conn.Database(name),
	}
}

func (d mongoImpl) Create(ctx context.Context, collection string, document interface{}) (interface{}, error) {
	log := logger.FromContext(ctx)
	res, err := d.conn.Collection(collection).InsertOne(ctx, document)
	if err != nil {
		log.Error(errMongoCreate, err.Error())
		return nil, errMongoCreate
	}
	return res.InsertedID, nil
}

func (d mongoImpl) CreateMany(ctx context.Context, collection string, documents []interface{}) (interface{}, error) {
	log := logger.FromContext(ctx)
	res, err := d.conn.Collection(collection).InsertMany(ctx, documents)
	if err != nil {
		log.Error(errMongoCreateMany, err.Error())
		return nil, errMongoCreateMany
	}
	return res.InsertedIDs, nil
}

func (d mongoImpl) Update(ctx context.Context, collection string, id string, document interface{}) (interface{}, error) {
	log := logger.FromContext(ctx)
	res, err := d.conn.Collection(collection).UpdateOne(ctx, bson.M{"_id": id}, document)
	if err != nil {
		log.Error(errMongoCreate, err.Error())
		return nil, errMongoCreate
	}
	return res.UpsertedID, nil
}

func (d mongoImpl) ReadOne(ctx context.Context, collection string, id string, model interface{}) error {
	log := logger.FromContext(ctx)
	res := d.conn.Collection(collection).FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		log.Error(errMongoReadOne, res.Err().Error())
		return errMongoReadOne
	}

	err := res.Decode(model)
	if err != nil {
		log.Error(errMongoReadOne, err.Error())
		return errMongoReadOne
	}

	return nil
}

func (d mongoImpl) ReadAll(ctx context.Context, collection string, search *model.SearchRequest, documents interface{}) error {
	log := logger.FromContext(ctx)
	skip := int64(search.Offset)
	limit := int64(search.Limit)
	cur, err := d.conn.Collection(collection).Find(ctx, search.GetFilter(), &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		log.Error(errMongoReadAll, err.Error())
		return errMongoReadAll
	}
	defer cur.Close(ctx)

	err = cur.All(ctx, documents)
	if err != nil {
		log.Error(errMongoReadAll, err.Error())
		return err
	}

	return nil
}
