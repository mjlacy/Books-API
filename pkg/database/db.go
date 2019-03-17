package database

import (
	"bookAPI"
	"context"
	"encoding/base64"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"

	//"github.com/globalsign/mgo"
	//"github.com/globalsign/mgo/bson"
	"github.com/mongodb/mongo-go-driver/mongo"

	//"github.com/mongodb/mongo-go-driver/mongo"
	//"github.com/mongodb/mongo-go-driver/bson"
	//"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/core/options"
)

type DatabaseConfig struct {
	DbURL          string `json:"DbURL"`
	DatabaseName   string `json:"DatabaseName"`
	CollectionName string `json:"CollectionName"`
}

type Repository struct{
	//session        *mgo.Session
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

	client, err := mongo.Connect(context.Background(), string(url), nil)
	if err != nil {
		return
	}
	database := client.Database(config.DatabaseName)

	r = &Repository{database: database, collection: database.Collection(config.CollectionName), databaseName: config.DatabaseName, collectionName: config.CollectionName}
	return

	//url, err := base64.StdEncoding.DecodeString(config.DbURL)
	//if err != nil {
	//	return
	//}
	//session, err := mgo.Dial(string(url))
	//if err != nil {
	//	return
	//}
	//session.SetMode(mgo.Monotonic, true)
	//
	//r = &Repository{session: session, databaseName: config.DatabaseName, collectionName: config.CollectionName}
	//return
}

//func (repo *Repository) Ping() error{
//	re := repo.session.Clone()
//	defer re.Close()
//	return re.Ping()
//}

func (repo Repository) GetBooks(s bookAPI.Book) (books []bookAPI.Book, err error){
	//conditions := bson.NewDocument()
	var originals []bookAPI.Book
	//var results []bookAPI.Book

	conditions := bson.NewArray()
	query := bson.NewDocument()

	if s.BookId != 0{
		conditions.Append(bson.VC.DocumentFromElements(bson.EC.Int32("bookId", s.BookId)))
	}

	if s.Title != ""{
		conditions.Append(bson.VC.DocumentFromElements(bson.EC.String("title", s.Title)))
	}

	if s.Author != ""{
		conditions.Append(bson.VC.DocumentFromElements(bson.EC.String("author", s.Author)))
	}

	if s.Year != 0{
		conditions.Append(bson.VC.DocumentFromElements(bson.EC.Int32("year", s.Year)))
	}

	if conditions.Len() != 0 {
		query = bson.NewDocument(bson.EC.Array("and", conditions))
	}

	cur, err := repo.collection.Find(context.Background(), query)
	if err != nil {
		return
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var elem bookAPI.Book

		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		originals = append(originals, elem)
	}
	if err:= cur.Err(); err != nil {
		 return nil, err
	}

	for _, value := range originals {
		books = append(books, value)
	}
	return

	//session := repo.session.Clone()
	//defer session.Close()
	//
	//conditions := []bson.M{}
	//query := bson.M{}
	//
	//var results []bookAPI.Book
	//
	//if s.BookId != 0{
	//	conditions = append(conditions, bson.M{"bookId": s.BookId})
	//}
	//
	//if s.Title != ""{
	//	conditions = append(conditions, bson.M{"title": s.Title})
	//}
	//
	//if s.Author != ""{
	//	conditions = append(conditions, bson.M{"author": s.Author})
	//}
	//
	//if s.Year != 0{
	//	conditions = append(conditions, bson.M{"year": s.Year})
	//}
	//
	//if len(conditions) != 0 {
	//	query = bson.M{"$and": conditions}
	//}
	//
	//err = session.DB(repo.databaseName).C(repo.collectionName).Find(query).All(&results)
	//if err == nil {
	//	b.Books = results
	//}
	//return
}

func (repo Repository) GetBookById(id string) (b *bookAPI.Book, err error){
	objectId, err := objectid.FromHex(id)
	if err != nil {
		return
	}

	filter := bson.NewDocument(bson.EC.ObjectID("_id", objectId))

	err = repo.collection.FindOne(context.Background(), filter).Decode(b)
	return

	//var oid bson.ObjectId
	//if bson.IsObjectIdHex(id){
	//	oid = bson.ObjectIdHex(id)
	//} else {
	//	return
	//}
	//session := repo.session.Clone()
	//defer session.Close()
	//err = session.DB(repo.databaseName).C(repo.collectionName).FindId(oid).One(&b)
	//return
}

func (repo Repository) PostBook(book *bookAPI.Book) (string, error){
	result, err := repo.collection.InsertOne(context.Background(), book)

	return result.InsertedID.(string), err

	//session := repo.session.Clone()
	//defer session.Close()
	//if book.Id.Hex() == ""{
	//	book.Id = bson.NewObjectId()
	//}
	//id = book.Id.Hex()
	//err = session.DB(repo.databaseName).C(repo.collectionName).Insert(book)
	//return
}

func (repo Repository) PutBook(id string, book *bookAPI.Book) (string, error){
	result, err := repo.collection.ReplaceOne(context.Background(), bson.NewDocument(bson.EC.String("_id", id)),
		book, options.OptUpsert(true))

	return result.UpsertedID.(string), err

	//if !bson.IsObjectIdHex(id){
	//	err = errors.New("Invalid id given")
	//	return
	//}
	//
	//if book.Id.Hex() != id{
	//	book.Id = bson.ObjectIdHex(id)
	//}
	//session := repo.session.Clone()
	//defer session.Close()
	//update, err := session.DB(repo.databaseName).C(repo.collectionName).UpsertId(bson.ObjectIdHex(id), book)
	//if update.UpsertedId != nil {
	//	updateId = update.UpsertedId.(bson.ObjectId).Hex()
	//}
	//return
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
	_, err = repo.collection.DeleteOne(context.Background(), bson.NewDocument(bson.EC.String("_id", id)))

	return

	//session := repo.session.Clone()
	//defer session.Close()
	//
	//err = session.DB(repo.databaseName).C(repo.collectionName).RemoveId(bson.ObjectIdHex(id)) //Will error if not found
	//return
}
