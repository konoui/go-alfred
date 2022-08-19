package update

import (
	"os"
	"time"
)

type timer struct {
	path    string
	modTime time.Time
	now     time.Time
}

func newTimer() (*timer, error) {
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(self)
	if err != nil {
		return nil, err
	}

	return &timer{
		path:    self,
		modTime: info.ModTime(),
		now:     time.Now(),
	}, nil
}

func (t *timer) checkout() error {
	updated := t.now
	return os.Chtimes(t.path, updated, updated)
}

func (t *timer) increase(hour time.Duration) error {
	updated := t.modTime.Add(hour)
	// Note: must not be future time
	if updated.After(t.now) {
		updated = t.now
	}
	return os.Chtimes(t.path, updated, updated)
}

func (t *timer) passed(ttl time.Duration) bool {
	return time.Since(t.modTime) > ttl
}
