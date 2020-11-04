package alfred

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/konoui/go-alfred/cache"
)

const (
	cacheDirKey = "cache-dir"
)

// ErrCacheExpired represent ttl is expired
var ErrCacheExpired = errors.New("cache expired")

type caches struct {
	suffix string
	caches sync.Map
}

// Cache wrapes cache.Cacher
// If cache load/store error occurs, workflow will continue to work
type Cache struct {
	cache.Cacher
	filename string
	ttl      time.Duration
	wf       *Workflow
	err      error
}

func (w *Workflow) getCacheDir() string {
	dir, ok := w.dirs[cacheDirKey]
	if ok {
		return dir
	}
	return os.TempDir()
}

// SetCacheDir overrides default cache directory
func (w *Workflow) SetCacheDir(dir string) (err error) {
	if _, err = os.Stat(dir); err != nil {
		return
	}
	w.dirs[cacheDirKey] = dir
	return
}

func (w *Workflow) getCacheSuffix() string {
	if w.cache.suffix != "" {
		return w.cache.suffix
	}

	bundleID := os.Getenv("alfred_workflow_bundleid")
	if bundleID != "" {
		return bundleID
	}

	return "-alfred.cache"
}

// SetCacheSuffix overrides suffix of default cache file
func (w *Workflow) SetCacheSuffix(suffix string) {
	w.cache.suffix = suffix
}

// Cache creates singleton instance
func (w *Workflow) Cache(key string) *Cache {
	filename := key + w.getCacheSuffix()
	if v, ok := w.cache.caches.Load(filename); ok {
		return v.(*Cache)
	}

	cr, err := cache.New(w.getCacheDir(), filename)
	if err != nil {
		err = fmt.Errorf("failed to create cache. try to use nil cacher: %w", err)
		w.logger.Println(err)
		cr = cache.NewNilCache()
	}

	c := &Cache{
		err:      err,
		wf:       w,
		Cacher:   cr,
		filename: filename,
		ttl:      0 * time.Second,
	}
	w.cache.caches.Store(filename, c)
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
	var err error
	defer func() {
		c.err = err
	}()

	var items Items
	if err = c.Load(&items); err != nil {
		return c
	}
	if c.Expired(c.ttl) {
		err = fmt.Errorf("%s ttl is expired: %w", c.filename, ErrCacheExpired)
		c.wf.logger.Println(err)
		return c
	}

	c.wf.std.Items = items
	return c
}

// StoreItems saves data into cache file
func (c *Cache) StoreItems() *Cache {
	var err error
	defer func() {
		c.err = err
	}()

	// Note: If there is no item, we avoid to save data into cache.
	// We define it is no error case
	items := c.wf.std.Items
	if len(items) == 0 {
		return c
	}

	if err = c.Store(items); err != nil {
		c.wf.logger.Println(err)
	}
	return c
}

// MaxAge sets cache TTL
func (c *Cache) MaxAge(ttl time.Duration) *Cache {
	defer func() {
		c.err = nil
	}()
	c.ttl = ttl
	return c
}

// ClearItems clear cache data
func (c *Cache) ClearItems() *Cache {
	var err error
	defer func() {
		c.err = err
	}()

	if err = c.Clear(); err != nil {
		c.wf.logger.Println(err)
	}
	return c
}
