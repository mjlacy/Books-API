package database

import (
	"bookAPI"
	"encoding/base64"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
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

func (repo Repository) GetBooks() (bookAPI.Books, error){
	session := repo.session.Clone()
	defer session.Close()
	var result []bookAPI.Book
	var results bookAPI.Books
	err := session.DB(repo.databaseName).C(repo.collectionName).Find(bson.M{}).All(&result)
	if err == nil {
		results.Books = result
	}
	return results, err
}

func (repo Repository) GetBookByBookId(id string) (b *bookAPI.Book, err error){
	session := repo.session.Clone()
	defer session.Close()
	idNum, _ := strconv.Atoi(id)
	err = session.DB(repo.databaseName).C(repo.collectionName).Find(bson.M{"bookId" : idNum}).One(&b)
	return
}

func (repo Repository) PostBook(book *bookAPI.Book) (err error){
	session := repo.session.Clone()
	defer session.Close()
	err = session.DB(repo.databaseName).C(repo.collectionName).Insert(book)
	return
}

func (repo Repository) PutBook(id string, book *bookAPI.Book) (err error){
	session := repo.session.Clone()
	defer session.Close()
	idNum, _ := strconv.Atoi(id)
	_, err = session.DB(repo.databaseName).C(repo.collectionName).Upsert(bson.M{"bookId" : idNum}, book)
	return
}

func (repo Repository) DeleteBook(id string) (err error){
	session := repo.session.Clone()
	defer session.Close()
	idNum, _ := strconv.Atoi(id)
	err = session.DB(repo.databaseName).C(repo.collectionName).Remove(bson.M{"bookId" : idNum})
	return
}
