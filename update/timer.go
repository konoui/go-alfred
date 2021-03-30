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

func (t *timer) increase(hour time.Duration) error {
	updated := t.now.Add(hour)
	if err := os.Chtimes(t.path, updated, updated); err != nil {
		return err
	}
	return nil
}

func (t *timer) passed(ttl time.Duration) bool {
	return time.Since(t.modTime) > ttl
}
