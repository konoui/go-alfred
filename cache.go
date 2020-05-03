package alfred

import (
	"fmt"
	"os"
	"time"

	"github.com/konoui/go-alfred/cache"
)

var cacheDefaultDir = os.TempDir()
var cacheDefaultSuffix = "-alfred.cache"

// Cache wrapes cache.Cacher
// If cache load/store error occurs, workflow will continue to work
type Cache struct {
	cache.Cacher
	filename string
	ttl      time.Duration
	wf       *Workflow
	err      error
}

type expiredError struct {
	err error
}

func (e *expiredError) Error() string {
	return e.err.Error()
}

// IsExpired return true if err equals cache expired error
func IsExpired(err error) bool {
	if _, ok := err.(*expiredError); ok {
		return true
	}
	return false
}

// SetCacheDir overrides default cache directory
func (w *Workflow) SetCacheDir(dir string) (err error) {
	if _, err = os.Stat(dir); err != nil {
		return
	}
	cacheDefaultDir = dir
	return
}

// SetCacheSuffix overrides suffix of default cache file
func (w *Workflow) SetCacheSuffix(suffix string) {
	cacheDefaultSuffix = suffix
}

// Cache creates singleton instance
func (w *Workflow) Cache(key string) *Cache {
	filename := key + cacheDefaultSuffix
	if v, ok := w.caches.Load(filename); ok {
		return v.(*Cache)
	}

	cr, err := cache.New(cacheDefaultDir, filename)
	if err != nil {
		err = fmt.Errorf("failed to create cache due to %v. try to work using nil cacher", err)
		cr = cache.NewNilCache()
	}

	c := &Cache{
		err:      err,
		wf:       w,
		Cacher:   cr,
		filename: filename,
		ttl:      0 * time.Second,
	}
	w.caches.Store(filename, c)

	return c
}

// Workflow returns workflow instance from cache one
func (c *Cache) Workflow() *Workflow {
	return c.wf
}

// Err returns cache operation err
func (c *Cache) Err() error {
	return c.err
}

// LoadItems reads data from cache file
func (c *Cache) LoadItems() *Cache {
	var items Items
	if err := c.Load(&items); err != nil {
		c.err = err
		return c
	}
	if c.Expired(c.ttl) {
		c.err = &expiredError{err: fmt.Errorf("cache expired")}
		return c
	}

	// update
	c.err = nil
	c.wf.std.Items = items
	return c
}

// StoreItems saves data into cache file
func (c *Cache) StoreItems() *Cache {
	// Note: If there is no item, we avoid to save data into cache.
	// We define it is no error case
	items := c.wf.std.Items
	if len(items) == 0 {
		c.err = nil
		return c
	}

	c.err = c.Store(items)
	return c
}

// MaxAge sets cache TTL
func (c *Cache) MaxAge(ttl time.Duration) *Cache {
	c.ttl = ttl
	return c
}

// Delete cache data
func (c *Cache) Delete() *Cache {
	c.err = c.Clear()
	return c
}
