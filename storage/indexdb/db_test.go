// Copyright 2015 Eleme Inc. All rights reserved.

package indexdb

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
	"github.com/eleme/banshee/util/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	fileName := "db-testing"
	db, err := Open(fileName)
	assert.Ok(t, err == nil)
	assert.Ok(t, util.IsFileExist(fileName))
	defer db.Close()
	defer os.RemoveAll(fileName)
}

func TestLoad(t *testing.T) {
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	idx := &models.Index{Name: "foo", Stamp: 1450430839, Score: 0.7, Average: 78.5}
	// Add one
	db.Put(idx)
	assert.Ok(t, db.m.Has(idx.Name))
	// Clear cache
	db.m.Clear()
	assert.Ok(t, db.m.Len() == 0)
	assert.Ok(t, !db.m.Has(idx.Name))
	// Reload
	db.load()
	// Must not empty and idx in cache
	assert.Ok(t, db.m.Len() == 1)
	assert.Ok(t, db.m.Has(idx.Name))
}

func TestPut(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Test.
	idx := &models.Index{Name: "foo", Stamp: 1450430837, Score: 1.2, Average: 109.5}
	err := db.Put(idx)
	assert.Ok(t, err == nil)
	// Must in cache
	assert.Ok(t, db.m.Has(idx.Name))
	// Must in db file
	v, err := db.db.Get([]byte(idx.Name), nil)
	assert.Ok(t, err == nil)
	idx1 := &models.Index{}
	decode(v, idx1)
	assert.Ok(t, idx1.Stamp == idx.Stamp)
	assert.Ok(t, idx1.Score == idx.Score)
	assert.Ok(t, idx1.Average == idx.Average)
}

func TestGet(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Not found.
	assert.Ok(t, ErrNotFound == db.Get(&models.Index{Name: "Not-Exist"}))
	// Put one.
	idx := &models.Index{Name: "foo", Stamp: 1450430837, Score: 0.3, Average: 100}
	db.Put(idx)
	// Get it from cache.
	i := &models.Index{Name: idx.Name}
	assert.Ok(t, nil == db.Get(i))
	assert.Ok(t, i.Equal(idx))
}

func TestDelete(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add one.
	idx := &models.Index{Name: "foo", Stamp: 1450430837, Score: 0.3, Average: 100}
	db.Put(idx)
	// Must in cache.
	assert.Ok(t, db.m.Has(idx.Name))
	// Delete it.
	err := db.Delete(idx.Name)
	assert.Ok(t, err == nil)
	// Must not exist in cache
	assert.Ok(t, !db.m.Has(idx.Name))
	// Must not in db.
	_, err = db.db.Get([]byte(idx.Name), nil)
	assert.Ok(t, err == leveldb.ErrNotFound)
	// Cant get again.
	assert.Ok(t, ErrNotFound == db.Get(idx))
}

func TestFilter(t *testing.T) {
	// Open db.
	fileName := "db-testing"
	db, _ := Open(fileName)
	defer db.Close()
	defer os.RemoveAll(fileName)
	// Add indexes.
	excludeName := "abfxyz"
	db.Put(&models.Index{Name: "abcefg"})
	db.Put(&models.Index{Name: "abcxyz"})
	db.Put(&models.Index{Name: excludeName})
	// Filter
	l := db.Filter("abc*")
	assert.Ok(t, len(l) == 2)
	assert.Ok(t, l[0].Name != excludeName && l[1].Name != excludeName)
}