package database

import (
	"bookAPI"
	"encoding/base64"
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

func (repo Repository) GetBooks() (b bookAPI.Books, err error){
	session := repo.session.Clone()
	defer session.Close()
	var result []bookAPI.Book
	err = session.DB(repo.databaseName).C(repo.collectionName).Find(bson.M{}).All(&result)
	if err == nil {
		b.Books = result
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

func (repo Repository) DeleteBook(_id string) (err error){
	session := repo.session.Clone()
	defer session.Close()
	err = session.DB(repo.databaseName).C(repo.collectionName).RemoveId(bson.ObjectIdHex(_id))
	return
}
