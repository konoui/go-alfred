package alfred

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ErrCacheExpired represent ttl is expired
var ErrCacheExpired = errors.New("cache expired")

type Cache struct {
	icache               internalCacher
	wf                   *Workflow
	maxAge               time.Duration
	staleWhileRevalidate time.Duration
	fetcher              Fetcher
}

type Cacher interface {
	MaxAge(time.Duration) CacheControlerOrLoader
	Store() error
	Clear() error
}

type CacheLoader interface {
	Load() error
}

type CacheControlerOrLoader interface {
	CacheLoader
	StaleWhileRevalidate(Fetcher) CacheLoader
}

type Fetcher func() (any, error)

func (w *Workflow) Cache(key string) Cacher {
	if key == "" {
		return &Cache{
			icache: newNilCache(),
			wf:     w,
		}
	}

	return &Cache{
		icache: &cache{
			dir:  GetCacheDir(),
			file: key + ".json",
		},
		wf: w,
	}
}

func (c *Cache) MaxAge(age time.Duration) CacheControlerOrLoader {
	c.maxAge = age
	return c
}

func (c *Cache) StaleWhileRevalidate(fetcher Fetcher) CacheLoader {
	c.staleWhileRevalidate = c.maxAge * 2
	c.fetcher = fetcher
	return c
}

func (c *Cache) Load() error {
	age := c.staleWhileRevalidate
	if age == 0 && c.icache.expired(c.maxAge) {
		return ErrCacheExpired
	}

	if age > 0 && c.icache.expired(c.staleWhileRevalidate) {
		cmd := exec.Command(os.Args[0], os.Args...) //nolint
		cmd.Env = os.Environ()
		c.wf.Job(GetBundleID()).Logging().Start(cmd)
		items, err := c.fetcher()
		if err != nil {
			c.wf.sLogger().Errorf("failed to fetcher: %w", err)
			return err
		}
		if err := c.icache.store(&items); err != nil {
			c.wf.sLogger().Errorf("failed to store: %w", err)
			return err
		}
	}

	err := c.icache.load(&c.wf.items)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Store() error {
	items := &c.wf.items
	return c.icache.store(items)
}

func (c *Cache) Clear() error {
	return c.icache.clear()
}

type internalCacher interface {
	load(any) error
	store(any) error
	clear() error
	expired(time.Duration) bool
}

// cache is file level cache
type cache struct {
	dir  string
	file string
}

// load read data saved cache into v
func (c *cache) load(v any) error {
	p := c.path()
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(v); err != nil {
		return fmt.Errorf("failed to load data from cache (%s): %w", p, err)
	}

	return nil
}

// store save data into cache
func (c *cache) store(v any) (err error) {
	f, err := os.CreateTemp(GetCacheDir(), GetBundleID())
	if err != nil {
		return err
	}

	old := f.Name()
	defer func() {
		cerr := f.Close()
		if err == nil && cerr != nil {
			err = cerr
			return
		}
		err = os.Rename(old, c.path())
	}()

	err = json.NewEncoder(f).Encode(v)
	if err != nil {
		return fmt.Errorf("failed to save data into cache (%s): %w", old, err)
	}

	return nil
}

// clear remove cache file if exist
func (c *cache) clear() error {
	p := c.path()
	if PathExists(p) {
		return os.Remove(p)
	}

	return nil
}

// expired return true if cache is expired
func (c *cache) expired(maxAge time.Duration) bool {
	age, err := c.age()
	if err != nil {
		return true
	}

	return age > maxAge
}

// nilCache noop cache which does nothing useful
type nilCache struct{}

// NewNilCache creates a new noop cache Instance
func newNilCache() internalCacher {
	return nilCache{}
}

// Load return nil
func (c nilCache) load(_ any) error {
	return nil
}

// Store return nil
func (c nilCache) store(_ any) error {
	return nil
}

// Clear return nil
func (c nilCache) clear() error {
	return nil
}

// Expired return true that means cache is always expired
func (c nilCache) expired(_ time.Duration) bool {
	return true
}

// age return the time since the data is cached at
func (c *cache) age() (time.Duration, error) {
	p := c.path()
	fi, err := os.Stat(p)
	if err != nil {
		return 0, err
	}

	return time.Since(fi.ModTime()), nil
}

// path return the path of cache file
func (c *cache) path() string {
	return filepath.Join(c.dir, c.file)
}
