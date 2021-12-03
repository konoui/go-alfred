package alfred

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/konoui/go-alfred/cache"
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

func (w *Workflow) getCacheSuffix() (suffix string) {
	suffix = w.cache.suffix
	if suffix != "" {
		return
	}

	suffix = w.GetBundleID()
	w.cache.suffix = suffix

	// Note default is bundle id
	return
}

// Cache creates singleton instance
// If key is empty, return Noop cache
func (w *Workflow) Cache(key string) *Cache {
	if key == "" {
		w.sLogger().Debugln("try to use nil cacher as cache key is empty")
		return newNil("", w, nil)
	}

	filename := key + w.getCacheSuffix()
	if v, ok := w.cache.caches.Load(filename); ok {
		return v.(*Cache)
	}

	cr, err := cache.New(w.GetCacheDir(), filename)
	if err != nil {
		err = fmt.Errorf("failed to create cache. try to use nil cacher: %w", err)
		w.sLogger().Errorln(err)
		return newNil(filename, w, err)
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
		c.wf.sLogger().Infoln(err.Error())
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
		c.wf.sLogger().Errorln(err)
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
		c.wf.sLogger().Errorln(err)
	}
	return c
}

func newNil(filename string, wf *Workflow, err error) *Cache {
	return &Cache{
		err:      err,
		iCache:   cache.NewNilCache(),
		filename: filename,
		wf:       wf,
	}
}
