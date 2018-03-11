package database

import (
	"gopkg.in/mgo.v2"
	"encoding/base64"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"BookAPI/pkg/models"
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

func InitializeMongoDatabase(config *DatabaseConfig) *Repository {
	url, err := base64.StdEncoding.DecodeString(config.DbURL)
	if err != nil {
		fmt.Println("Error base64 decoding connection string")
	}
	session, err := mgo.Dial(string(url))
	if err != nil {
		fmt.Println("Error connecting to database")
	}
	session.SetMode(mgo.Monotonic, true)

	return &Repository{session: session, databaseName: config.DatabaseName, collectionName: config.CollectionName}
}

func (repo Repository) GetBook() ([]models.Book, error){
	session := repo.session.Clone()
	defer session.Close()
	var result []models.Book
	err := session.DB(repo.databaseName).C(repo.collectionName).Find(bson.M{}).All(&result)
	return result, err
}

func (repo Repository) GetBookById(id string) (*models.Book, error){
	session := repo.session.Clone()
	defer session.Close()
	var result *models.Book
	err := session.DB(repo.databaseName).C(repo.collectionName).FindId(bson.ObjectIdHex(id)).One(&result)
	return result, err
}

func (repo Repository) PostBook(book *models.Book) (error){
	session := repo.session.Clone()
	defer session.Close()
	err := session.DB(repo.databaseName).C(repo.collectionName).Insert(book)
	return err
}

func (repo Repository) PutBook(_id string, book *models.Book) (error){
	session := repo.session.Clone()
	defer session.Close()
	_, err := session.DB(repo.databaseName).C(repo.collectionName).UpsertId(bson.ObjectIdHex(_id), book)
	return err
}

func (repo Repository) DeleteBook(_id string) (error){
	session := repo.session.Clone()
	defer session.Close()
	err := session.DB(repo.databaseName).C(repo.collectionName).RemoveId(bson.ObjectIdHex(_id))
	return err
}
