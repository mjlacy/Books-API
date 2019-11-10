package mongo

import (
	"BookAPI/internal"

	"context"
	"encoding/base64"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	FiveSeconds = 5 * time.Second
	TenSeconds  = 10 * time.Second
)

type DatabaseConfig struct {
	DbURL          string `json:"DbURL"`
	DatabaseName   string `json:"DatabaseName"`
	CollectionName string `json:"CollectionName"`
}

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func InitializeMongoDatabase(config *DatabaseConfig) (r *Repository, err error) {
	url, err := base64.StdEncoding.DecodeString(config.DbURL)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), TenSeconds)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(string(url)), nil)
	if err != nil {
		return
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	database := client.Database(config.DatabaseName)

	r = &Repository{client: client, collection: database.Collection(config.CollectionName)}
	return
}

func (repo *Repository) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()
	return repo.client.Ping(ctx, readpref.Primary())
}

func (repo Repository) GetBooks(book internal.Book) (books internal.Books, err error) {
	books.Books = []internal.Book{}

	conditions := bson.M{}

	if book.BookId != 0 {
		conditions["bookId"] = bson.M{"$eq": book.BookId}
	}

	if book.Title != "" {
		conditions["title"] = bson.M{"$eq": book.Title}
	}

	if book.Author != "" {
		conditions["author"] = bson.M{"$eq": book.Author}
	}

	if book.Year != 0 {
		conditions["year"] = bson.M{"$eq": book.Year}
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	cur, err := repo.collection.Find(ctx, conditions)
	if err != nil {
		return
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result internal.Book
		err = cur.Decode(&result)
		if err != nil {
			return
		}
		books.Books = append(books.Books, result)
	}
	err = cur.Err()
	return
}

func (repo Repository) GetBookById(id string) (book *internal.Book, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, internal.ErrNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	err = repo.collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&book)
	if err == mongo.ErrNoDocuments {
		err = internal.ErrNotFound
	}

	return
}

func (repo Repository) PostBook(book *internal.Book) (id string, returnedBook *internal.Book, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	result, err := repo.collection.InsertOne(ctx, book)
	if err != nil {
		return
	}

	book.Id = result.InsertedID.(primitive.ObjectID)

	return book.Id.Hex(), book, err
}

func (repo Repository) PutBook(id string, book *internal.Book) (isCreated bool, returnedBook *internal.Book, err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, nil, internal.ErrInvalidId
	}

	upsert := options.ReplaceOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	result, err := repo.collection.ReplaceOne(ctx, bson.M{"_id": objectId}, book, upsert.SetUpsert(true))
	if err != nil {
		return false, nil, err
	}

	returnedBook = book
	returnedBook.Id = objectId

	return result.UpsertedCount > 0, returnedBook, err
}

func (repo Repository) PatchBook(id string, update internal.Book) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return internal.ErrNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	result, err := repo.collection.UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": buildMap(update)})
	if err == nil && result.ModifiedCount == 0 {
		err = internal.ErrNotFound
	}
	return
}

func (repo Repository) DeleteBook(id string) (err error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return internal.ErrNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), FiveSeconds)
	defer cancel()

	result, err := repo.collection.DeleteOne(ctx, bson.M{"_id": objectId})
	if err == nil && result.DeletedCount == 0 {
		err = internal.ErrNotFound
	}

	return
}

func buildMap(book internal.Book) bson.M {
	update := make(bson.M)

	if book.BookId != 0 {
		update["bookId"] = book.BookId
	}

	if book.Title != "" {
		update["title"] = book.Title
	}

	if book.Author != "" {
		update["author"] = book.Author
	}

	if book.Year != 0 {
		update["year"] = book.Year
	}

	return update
}
