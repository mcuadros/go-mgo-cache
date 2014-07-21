package mgocache

import (
	"bytes"
	"testing"
)

import . "launchpad.net/gocheck"
import "labix.org/v2/mgo"

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

	key := "testKey"
	_, ok := cache.Get(key)

	c.Assert(ok, Equals, false)

	val := []byte("some bytes")
	cache.Set(key, val)

	retVal, ok := cache.Get(key)
	c.Assert(ok, Equals, true)
	c.Assert(bytes.Equal(retVal, val), Equals, true)

	cache.Delete(key)

	_, ok = cache.Get(key)
	c.Assert(ok, Equals, false)
}
