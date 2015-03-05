package mgocache

import (
	"testing"
	
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type S struct{}

var _ = Suite(&S{})

func (s *S) Test(c *C) {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		c.Fatalf("Can't connect to mongo, go error %v\n", err)
	}

	collection := session.DB("test").C("foo")
	defer session.Close()

	cache := New(collection)
	cache.Indexes()

	key := "testKey"
	_, ok := cache.Get(key)

	c.Assert(ok, Equals, false)

	val := []byte("some bytes")
	cache.Set(key, val)

	retVal, ok := cache.Get(key)
	c.Assert(ok, Equals, true)
	c.Assert(string(retVal), Equals, string(val))

	val = []byte("some other bytes")
	cache.Set(key, val)

	retVal, ok = cache.Get(key)
	c.Assert(ok, Equals, true)
	c.Assert(string(retVal), Equals, string(val))

	cache.Delete(key)

	_, ok = cache.Get(key)
	c.Assert(ok, Equals, false)
}
