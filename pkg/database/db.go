package database

import (
	"bookAPI"

	"encoding/base64"
	"math"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type DatabaseConfig struct {
	DbURL          string `json:"DbURL"`
	DatabaseName   string `json:"DatabaseName"`
	CollectionName string `json:"CollectionName"`
}

type Repository struct {
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

func (repo *Repository) Ping() error {
	re := repo.session.Clone()
	defer re.Close()
	return re.Ping()
}

func (repo Repository) GetBooks(s bookAPI.Book) (b bookAPI.Books, err error) {
	session := repo.session.Clone()
	defer session.Close()

	conditions := []bson.M{}
	query := bson.M{}

	var results []bookAPI.Book

	if s.BookId != 0 {
		conditions = append(conditions, bson.M{"bookId": s.BookId})
	}

	if s.Title != "" {
		conditions = append(conditions, bson.M{"title": s.Title})
	}

	if s.Author != "" {
		conditions = append(conditions, bson.M{"author": s.Author})
	}

	if s.Year != 0 {
		conditions = append(conditions, bson.M{"year": s.Year})
	}

	if len(conditions) != 0 {
		query = bson.M{"$and": conditions}
	}

	err = session.DB(repo.databaseName).C(repo.collectionName).Find(query).All(&results)
	if err == nil {
		for _, book := range results {
			b.Books = append(b.Books, book)
		}
	}
	return
}

func (repo Repository) GetBookById(id string) (b *bookAPI.Book, err error) {
	var oid bson.ObjectId
	if bson.IsObjectIdHex(id){
		oid = bson.ObjectIdHex(id)
	} else {
		return
	}
	session := repo.session.Clone()
	defer session.Close()

	err = session.DB(repo.databaseName).C(repo.collectionName).FindId(oid).One(&b)
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	return
}

func (repo Repository) PostBook(b *bookAPI.Book) (id string, err error) {
	session := repo.session.Clone()
	defer session.Close()

	err = session.DB(repo.databaseName).C(repo.collectionName).Insert(b)
	return
}

func (repo Repository) PutBook(id string, b *bookAPI.Book) (updateId string, err error) {
	if !bson.IsObjectIdHex(id) {
		id = bson.NewObjectId().Hex()
	}

	session := repo.session.Clone()
	defer session.Close()
	update, err := session.DB(repo.databaseName).C(repo.collectionName).UpsertId(bson.ObjectIdHex(id), b)
	if update != nil && update.UpsertedId != nil {
		updateId = update.UpsertedId.(bson.ObjectId).Hex()
	}
	return
}

func (repo Repository) PatchBook(id string, update bson.M) (err error) {
	var oid bson.ObjectId
	if bson.IsObjectIdHex(id){
		oid = bson.ObjectIdHex(id)
	} else {
		err = bookAPI.ErrNotFound
		return
	}
	session := repo.session.Clone()
	defer session.Close()

	for k, v := range update { //converts incoming float64's to int's
		if v, ok := v.(float64); ok {
			_, decimal := math.Modf(v)
			if decimal == 0 {
				update[k] = int(v)
			}
		}
	}

	err = session.DB(repo.databaseName).C(repo.collectionName).UpdateId(oid, bson.M{"$set": update})
	if err != nil && err.Error() =="not found" {
		err = bookAPI.ErrNotFound
	}
	return
}

func (repo Repository) DeleteBook(id string) (err error) {
	var oid bson.ObjectId
	if bson.IsObjectIdHex(id){
		oid = bson.ObjectIdHex(id)
	} else {
		err = bookAPI.ErrNotFound
		return
	}

	session := repo.session.Clone()
	defer session.Close()

	err = session.DB(repo.databaseName).C(repo.collectionName).RemoveId(oid)
	if err != nil && err.Error() =="not found" {
		err = bookAPI.ErrNotFound
	}
	return
}
