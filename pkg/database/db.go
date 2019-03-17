package database

import (
	"bookAPI"
	"context"
	"encoding/base64"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConfig struct {
	DbURL          string `json:"DbURL"`
	DatabaseName   string `json:"DatabaseName"`
	CollectionName string `json:"CollectionName"`
}

type Repository struct{
	database       *mongo.Database
	collection     *mongo.Collection
	databaseName   string
	collectionName string
}

func InitializeMongoDatabase(config *DatabaseConfig) (r *Repository, err error) {
	url, err := base64.StdEncoding.DecodeString(config.DbURL)
	if err != nil {
		return
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(string(url)), nil)
	if err != nil {
		return
	}
	database := client.Database(config.DatabaseName)

	r = &Repository{database: database, collection: database.Collection(config.CollectionName), databaseName: config.DatabaseName, collectionName: config.CollectionName}
	return
}

//func (repo *Repository) Ping() error{
//	re := repo.session.Clone()
//	defer re.Close()
//	return re.Ping()
//}

func (repo Repository) GetBooks(s bookAPI.Book) (books bookAPI.Books, err error){
	var originals []bookAPI.Book

	conditions := []bson.E{}

	if s.BookId != 0{
		conditions = append(conditions, bson.E{"bookId", s.BookId})
	}

	if s.Title != ""{
		conditions = append(conditions, bson.E{"title", s.Title})
	}

	if s.Author != ""{
		conditions = append(conditions, bson.E{"author", s.Author})
	}

	if s.Year != 0{
		conditions = append(conditions, bson.E{"year", s.Year})
	}

	cur, err := repo.collection.Find(context.Background(), conditions)
	if err != nil {
		return
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var elem bookAPI.Book

		err := cur.Decode(&elem)
		if err != nil {
			return bookAPI.Books{}, err
		}
		originals = append(originals, elem)
	}
	if err:= cur.Err(); err != nil {
		 return bookAPI.Books{}, err
	}

	for _, value := range originals {
		books.Books = append(books.Books, value)
	}
	return
}

func (repo Repository) GetBookById(id string) (b *bookAPI.Book, err error){
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	err = repo.collection.FindOne(context.Background(), bson.D{{"_id", objectId}}).Decode(&b)

	if err != nil && err.Error() == "mongo: no documents in result" {
		return nil, nil
	}

	return
}

func (repo Repository) PostBook(book *bookAPI.Book) (string, error){
	if book.Id == primitive.NilObjectID {
		book.Id = primitive.NewObjectID()
	}

	_, err := repo.collection.InsertOne(context.Background(), book)

	return book.Id.Hex(), err
}

func (repo Repository) PutBook(id string, book *bookAPI.Book) (bool, *bookAPI.Book, error){
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		err = errors.New("invalid id given")
		return false, nil, err
	}

	book.Id, _ = primitive.ObjectIDFromHex(id)

	upsert := options.ReplaceOptions{}

	result, err := repo.collection.ReplaceOne(context.Background(), bson.D{{"_id", objectId}},
		book, upsert.SetUpsert(true))

	return result.ModifiedCount > 0, book, err
}

//func (repo Repository) PatchBook(id string, update bson.M) (err error){
//	var oid bson.ObjectId
//	if bson.IsObjectIdHex(id){
//		oid = bson.ObjectIdHex(id)
//	} else {
//		err = errors.New("Invalid id given")
//		return
//	}
//	session := repo.session.Clone()
//	defer session.Close()
//
//	for k, v := range update { //converts incoming float64's to int's
//		if v, ok := v.(float64); ok {
//			_, decimal := math.Modf(v)
//			if decimal == 0 {
//				update[k] = int(v)
//			}
//		}
//	}
//
//	err = session.DB(repo.databaseName).C(repo.collectionName).UpdateId(oid, bson.M{"$set": update}) //Will error if not found
//	return
//}

func (repo Repository) DeleteBook(id string) (err error){
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil{
		err = errors.New("invalid id given")
		return
	}

	_, err = repo.collection.DeleteOne(context.Background(), bson.D{{"_id", objectId}})
	return
}
