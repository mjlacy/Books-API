package database

import (
	//"gopkg.in/mgo.v2"
	"encoding/base64"
	"fmt"
	//"gopkg.in/mgo.v2/bson"
	"BookAPI/pkg/models"
	//"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"log"
	"strconv"
	"github.com/mongodb/mongo-go-driver/bson"
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

func InitializeMongoDatabase(config *DatabaseConfig) *Repository {
	url, err := base64.StdEncoding.DecodeString(config.DbURL)
	if err != nil {
		fmt.Println("Error base64 decoding connection string")
	}
	client, err := mongo.Connect(context.Background(), "mongodb://" + string(url), nil)
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	database := client.Database(config.DatabaseName)
	collection := database.Collection(config.CollectionName)

	return &Repository{database: database, collection: collection, databaseName: config.DatabaseName, collectionName: config.CollectionName}
}

func (repo Repository) GetBook() (models.Books, error){
	var result []models.Book
	var results models.Books
	cur, err := repo.collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var elem models.Book
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	if err == nil {
		results.Books = result
	}
	return results, err
}

func (repo Repository) GetBookById(id string) (*models.Book, error){
	result := bson.NewDocument()
	idNum, _ := strconv.Atoi(id)
	filter := bson.NewDocument(bson.EC.Int32("BookId", int32(idNum)))
	err := repo.collection.FindOne(context.Background(), filter).Decode(result)
	if err != nil { log.Fatal(err) }

	book := models.Book{result.ElementAt(0).Value().ObjectID().String(), result.ElementAt(1).Value().Int32(),
		result.ElementAt(2).Value().StringValue(), result.ElementAt(3).Value().StringValue(),
		result.ElementAt(4).Value().Int32()}

	return &book, err
}

func (repo Repository) PostBook(book *models.Book) (error){
	_, err := repo.collection.InsertOne(context.Background(), book)
	return err
}

func (repo Repository) PutBook(id string, book *models.Book) (error){
	idNum, _ := strconv.Atoi(id)
	_, err := repo.collection.UpdateOne(context.Background(), bson.NewDocument(bson.EC.Int32("BookId", int32(idNum))), bson.NewDocument(
		bson.EC.SubDocumentFromElements("$set",
			bson.EC.Int32("BookId", book.BookId),
			bson.EC.String("Name", book.Name),
			bson.EC.String("Author", book.Author),
			bson.EC.Int32("Year", book.Year),
		),
	),)
	return err
}

func (repo Repository) DeleteBook(id string) (error){
	idNum, _ := strconv.Atoi(id)
	_, err := repo.collection.DeleteOne(context.Background(), bson.NewDocument(bson.EC.Int32("BookId", int32(idNum))))
	return err
}
