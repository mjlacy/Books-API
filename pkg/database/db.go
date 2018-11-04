package database

import (
	"bookAPI"
	"encoding/base64"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DatabaseConfig struct {
	DbURL          string `json:"DbURL"`
	DatabaseName   string `json:"DatabaseName"`
	CollectionName string `json:"CollectionName"`
}

type Repository struct{
	session        *mgo.Session
	databaseName   string
	collectionName string
}

func InitializeMongoDatabase(config *DatabaseConfig) (r *Repository, err error) {
	url, err := base64.StdEncoding.DecodeString(config.DbURL)
	if err != nil {
		return
	}
	session, err := mgo.Dial(string(url))
	if err != nil {
		return
	}
	session.SetMode(mgo.Monotonic, true)

	r = &Repository{session: session, databaseName: config.DatabaseName, collectionName: config.CollectionName}
	return
}

func (repo *Repository) Ping() error{
	re := repo.session.Clone()
	defer re.Close()
	return re.Ping()
}

func (repo Repository) GetBooks(s bookAPI.Book) (b bookAPI.Books, err error){
	session := repo.session.Clone()
	defer session.Close()

	conditions := []bson.M{}
	query := bson.M{}

	var results []bookAPI.Book

	if s.BookId != 0{
		conditions = append(conditions, bson.M{"bookId": s.BookId})
	}

	if s.Title != ""{
		conditions = append(conditions, bson.M{"title": s.Title})
	}

	if s.Author != ""{
		conditions = append(conditions, bson.M{"author": s.Author})
	}

	if s.Year != 0{
		conditions = append(conditions, bson.M{"year": s.Year})
	}

	if len(conditions) != 0 {
		query = bson.M{"$and": conditions}
	}

	err = session.DB(repo.databaseName).C(repo.collectionName).Find(query).All(&results)
	if err == nil {
		b.Books = results
	}
	return
}

func (repo Repository) GetBookById(id string) (b *bookAPI.Book, err error){
	var oid bson.ObjectId
	if bson.IsObjectIdHex(id){
		oid = bson.ObjectIdHex(id)
	} else {
		return
	}
	session := repo.session.Clone()
	defer session.Close()
	err = session.DB(repo.databaseName).C(repo.collectionName).FindId(oid).One(&b)
	return
}

func (repo Repository) PostBook(book *bookAPI.Book) (id string, err error){
	session := repo.session.Clone()
	defer session.Close()
	if book.Id.Hex() == ""{
		book.Id = bson.NewObjectId()
	}
	id = book.Id.Hex()
	err = session.DB(repo.databaseName).C(repo.collectionName).Insert(book)
	return
}

func (repo Repository) PutBook(id string, book *bookAPI.Book) (updateId string, err error){
	if !bson.IsObjectIdHex(id){
		return
	}

	if book.Id.Hex() != id{
		book.Id = bson.ObjectIdHex(id)
	}
	session := repo.session.Clone()
	defer session.Close()
	update, err := session.DB(repo.databaseName).C(repo.collectionName).UpsertId(bson.ObjectIdHex(id), book)
	if update.UpsertedId != nil {
		updateId = update.UpsertedId.(bson.ObjectId).Hex()
	}
	return
}

func (repo Repository) DeleteBook(id string) (err error){
	session := repo.session.Clone()
	defer session.Close()

	book, err := repo.GetBookById(id)
	if err != nil {
		return
	}

	if book == nil {
		return errors.New("not found")
	}

	err = session.DB(repo.databaseName).C(repo.collectionName).RemoveId(bson.ObjectIdHex(id))
	return
}
