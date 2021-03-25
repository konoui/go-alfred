package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type gitHubUpdater struct {
	client                *github.Client
	currentVersion        string
	owner                 string
	repo                  string
	vFormat               bool
	newerVersionAvailable bool
	checkInterval         time.Duration
	fetchURL              string
}

func NewGitHubSource(owner, repo string, opts ...Option) UpdaterSource {
	g := &gitHubUpdater{
		client:        github.NewClient(nil),
		owner:         owner,
		repo:          repo,
		checkInterval: 24 * 7 * 2 * time.Hour,
		vFormat:       false,
	}

	for _, opt := range opts {
		opt(g)
	}
	return g

}

func (g *gitHubUpdater) NewerVersionAvailable(currentVersion string) (bool, error) {
	g.currentVersion = currentVersion
	ok, _, err := g.newerVersionAvailableContext(context.Background(), currentVersion)
	return ok, err
}

func (g *gitHubUpdater) Update(ctx context.Context) error {
	ok, url, err := g.newerVersionAvailableContext(ctx, g.currentVersion)
	if err != nil {
		return err
	}
	if ok {
		return updateContext(ctx, url)
	}
	return nil
}

func (g *gitHubUpdater) setCheckInterval(interval time.Duration) {
	g.checkInterval = interval
}

func (g *gitHubUpdater) enableVFormat() {
	g.vFormat = true
}

func (g *gitHubUpdater) newerVersionAvailableContext(ctx context.Context, currentVersion string) (ok bool, url string, err error) {
	if g.newerVersionAvailable && g.fetchURL != "" {
		return true, g.fetchURL, nil
	}

	t, err := newTimer()
	if err != nil {
		return false, "", err
	}
	defer func() {
		// increase interval if error occurs
		if err != nil {
			_ = t.increase(1 * time.Hour)
		}
	}()

	if !t.passed(g.checkInterval) {
		return false, "", nil
	}

	tag, url, err := getLatestAssetInfo(ctx, g.client, g.owner, g.repo)
	if err != nil {
		return false, "", err
	}

	ok, err = compareVersions(tag, currentVersion, g.vFormat)
	if err != nil {
		return false, "", fmt.Errorf("version formats are invalid %s, %s: %w",
			currentVersion, tag, err)
	}
	if !ok {
		// if current is latest, check next time after interval
		return false, url, t.increase(g.checkInterval)
	}
	// newer version available
	d := func() {
		g.fetchURL = url
		g.newerVersionAvailable = true
	}
	d()
	return true, url, nil
}

func getLatestAssetInfo(ctx context.Context, client *github.Client, owner, repo string) (_, _ string, _ error) {
	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		return "", "", err
	}

	if latestRelease == nil || latestRelease.TagName == nil {
		return "", "", errors.New("found no release")
	}

	for _, asset := range latestRelease.Assets {
		if asset.Name == nil {
			continue
		}
		if hasAlfredWorkflowExt(*asset.Name) {
			return *latestRelease.TagName, asset.GetBrowserDownloadURL(), nil
		}
	}
	return *latestRelease.TagName, "", errors.New("found no alfredworkflow assets")
}