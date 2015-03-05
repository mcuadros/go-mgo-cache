// Package mgocache provides an implementation of httpcache.Cache that stores and
// retrieves data using Mongo.
package mgocache

import (
	"log"
	"time"
	
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Cache objects store and retrieve data using Mongo.
type Cache struct {
	// mongo collection where the cache will be stored
	Collection *mgo.Collection
}

// New returns a new Cache
func New(collection *mgo.Collection) *Cache {
	return &Cache{
		Collection: collection,
	}
}

func (self *Cache) Get(key string) (resp []byte, ok bool) {
	result := record{}
	err := self.Collection.Find(bson.M{"key": key}).One(&result)
	if err != nil {
		return []byte{}, false
	}

	return result.Content, true
}

func (self *Cache) Set(key string, content []byte) {
	_, err := self.Collection.Upsert(bson.M{"key": key}, &record{
		Created: time.Now(),
		Updated: time.Now(),
		Key:     key,
		Content: content,
	})

	if err != nil {
		log.Printf("Can't insert record in mongo: %v\n", err)
		return
	}

	return
}

func (self *Cache) Delete(key string) {
	err := self.Collection.Remove(bson.M{"key": key})
	if err != nil {
		log.Printf("Can't remove record: %s", err)
	}
}

func (self *Cache) Indexes() {
	index := mgo.Index{
		Key:      []string{"key"},
		Unique:   true,
		DropDups: true,
	}

	err := self.Collection.EnsureIndex(index)
	if err != nil {
		log.Printf("Can't ensure index: %s", err)
	}
}

type record struct {
	Created time.Time
	Updated time.Time
	Key     string
	Content []byte
}
