package alfred

import (
	"errors"
	"fmt"
	"time"

	"github.com/konoui/go-alfred/cache"
)

// ErrCacheExpired represent ttl is expired
var ErrCacheExpired = errors.New("cache expired")

// Cache wrapes cache.Cacher
// If cache load/store error occurs, workflow will continue to work
type Cache struct {
	iCache   cache.Cacher
	filename string
	wf       *Workflow
	err      error
	maxAge   time.Duration
}

type Cacher interface {
	Workflow() *Workflow
	MaxAge(time.Duration) CacheLoader
	Err() error
	StoreItems() Cacher
	ClearItems() Cacher
}

type CacheLoader interface {
	LoadItems() Cacher
}

func (w *Workflow) getCacheSuffix() (suffix string) {
	suffix = w.customEnvs.cacheSuffix
	if suffix != "" {
		return
	}

	// default value is empty
	return ""
}

// Cache creates singleton instance
// If key is empty, return Noop cache
func (w *Workflow) Cache(key string) Cacher {
	if key == "" {
		w.sLogger().Debugln("try to use nil cacher as cache key is empty")
		return newNil("", w, nil)
	}

	filename := key + w.getCacheSuffix()
	if v, ok := w.cache.Load(filename); ok {
		return v.(*Cache)
	}

	cr, err := cache.New(w.GetCacheDir(), filename)
	if err != nil {
		err = fmt.Errorf("failed to create a cache. try to use a nil cacher: %w", err)
		w.sLogger().Errorln(err)
		return newNil(filename, w, err)
	}

	c := &Cache{
		err:      err,
		wf:       w,
		iCache:   cr,
		filename: filename,
	}
	w.cache.Store(filename, c)
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

func (c *Cache) MaxAge(age time.Duration) CacheLoader {
	c.maxAge = age
	return c
}

// LoadItems reads data from cache file
func (c *Cache) LoadItems() Cacher {
	var err error
	defer func() {
		c.err = err
	}()

	var items Items
	if err = c.iCache.Load(&items); err != nil {
		return c
	}

	if c.iCache.Expired(c.maxAge) {
		err = fmt.Errorf("%s ttl is expired: %w", c.filename, ErrCacheExpired)
		c.wf.sLogger().Infoln(err.Error())
		return c
	}

	c.wf.items = items
	return c
}

// StoreItems saves data into cache file
func (c *Cache) StoreItems() Cacher {
	var err error
	defer func() {
		c.err = err
	}()

	// Note: If there is no item, we avoid to save data into cache.
	// We define it is no error case
	if c.wf.IsEmpty() {
		return c
	}

	items := c.wf.items
	if err = c.iCache.Store(&items); err != nil {
		c.wf.sLogger().Errorln(err)
	}
	return c
}

// ClearItems clear cache data
func (c *Cache) ClearItems() Cacher {
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
