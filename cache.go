package alfred

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	iCache   cache.Cacher
	filename string
	wf       *Workflow
	err      error
}

func (w *Workflow) getCacheDir() (dir string) {
	dir, ok := w.dirs[cacheDirKey]
	if ok {
		return
	}

	dir = tmpDir
	w.dirs[cacheDirKey] = dir
	return
}

// SetCacheDir overrides default cache directory
func (w *Workflow) SetCacheDir(dir string) (err error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(abs); err != nil {
		return err
	}

	w.dirs[cacheDirKey] = abs
	return
}

func (w *Workflow) getCacheSuffix() (suffix string) {
	suffix = w.cache.suffix
	if suffix != "" {
		return
	}

	defer func() {
		// Note set env or default value to workflow to reduce Getenv call
		w.cache.suffix = suffix
	}()

	suffix = os.Getenv("alfred_workflow_bundleid")
	if suffix != "" {
		return
	}

	// Note default is empty
	return ""
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
		w.logger.Warnln(err)
		cr = cache.NewNilCache()
	}

	c := &Cache{
		err:      err,
		wf:       w,
		iCache:   cr,
		filename: filename,
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
func (c *Cache) LoadItems(ttl time.Duration) *Cache {
	var err error
	defer func() {
		c.err = err
	}()

	var items Items
	if err = c.iCache.Load(&items); err != nil {
		return c
	}

	if c.iCache.Expired(ttl) {
		err = fmt.Errorf("%s ttl is expired: %w", c.filename, ErrCacheExpired)
		c.wf.logger.Infoln(err)
		return c
	}

	c.wf.std.items = items
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
	if c.wf.IsEmpty() {
		return c
	}

	items := c.wf.std.items
	if err = c.iCache.Store(&items); err != nil {
		c.wf.logger.Warnln(err)
	}
	return c
}

// ClearItems clear cache data
func (c *Cache) ClearItems() *Cache {
	var err error
	defer func() {
		c.err = err
	}()

	if err = c.iCache.Clear(); err != nil {
		c.wf.logger.Warnln(err)
	}
	return c
}
