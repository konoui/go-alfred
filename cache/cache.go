package cache

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Cacher implements a simple store/load API, saving data to specified directory.
type Cacher interface {
	Load(interface{}) error
	Store(interface{}) error
	Clear() error
	Expired(time.Duration) bool
}

// Cache is directory level cache
type Cache struct {
	Dir  string
	File string
	sync.Mutex
}

// New create a new cache Instance
func New(dir, file string) (Cacher, error) {
	if !pathExists(dir) {
		return &Cache{}, fmt.Errorf("%s directory does not exist", dir)
	}

	return &Cache{
		Dir:  dir,
		File: file,
	}, nil
}

// Load read data saved cache into v
func (c *Cache) Load(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	p := c.path()
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = gob.NewDecoder(f).Decode(v); err != nil {
		return errors.Wrapf(err, "failed to load data from cache (%s).", p)
	}

	return nil
}

// Store save data into cache
func (c *Cache) Store(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	p := c.path()
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = gob.NewEncoder(f).Encode(v); err != nil {
		return errors.Wrapf(err, "failed to save data into cache (%s)", p)
	}

	return nil
}

// Clear remove cache file if exist
func (c *Cache) Clear() error {
	c.Lock()
	defer c.Unlock()
	p := c.path()
	if pathExists(p) {
		return os.Remove(p)
	}

	return nil
}

// Expired return true if cache is expired
func (c *Cache) Expired(ttl time.Duration) bool {
	age, err := c.age()
	if err != nil {
		return true
	}

	return age > ttl
}

// age return the time since the data is cached at
func (c *Cache) age() (time.Duration, error) {
	p := c.path()
	fi, err := os.Stat(p)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Since(fi.ModTime()), nil
}

// path return the path of cache file
func (c *Cache) path() string {
	return filepath.Join(c.Dir, c.File)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
