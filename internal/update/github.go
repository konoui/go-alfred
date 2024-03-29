//go:generate mockgen -source=$GOFILE -destination=mock_github/$GOFILE -package=github_$GOPACKAGE
package update

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type gitHubUpdater struct {
	client              RepositoriesService
	currentVersion      string
	owner               string
	repo                string
	newVersionAvailable bool
	checkInterval       time.Duration
	fetchURL            string
}

type RepositoriesService interface {
	GetLatestRelease(context.Context, string, string) (*github.RepositoryRelease, *github.Response, error)
}

func NewGitHubSource(owner, repo, currentVersion string, opts ...Option) UpdaterSource {
	g := &gitHubUpdater{
		client:         github.NewClient(nil).Repositories,
		owner:          owner,
		repo:           repo,
		currentVersion: currentVersion,
		checkInterval:  24 * 7 * 2 * time.Hour,
	}

	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *gitHubUpdater) IsNewVersionAvailable(ctx context.Context) (bool, error) {
	ok, _, err := g.newVersionAvailableContext(ctx, g.currentVersion)
	return ok, err
}

func (g *gitHubUpdater) IfNewVersionAvailable() Updater {
	return g
}

func (g *gitHubUpdater) Update(ctx context.Context) error {
	ok, url, err := g.newVersionAvailableContext(ctx, g.currentVersion)
	if err != nil {
		return err
	}
	if ok {
		return updateContext(ctx, url)
	}
	return nil
}

func (g *gitHubUpdater) SetCheckInterval(interval time.Duration) {
	g.checkInterval = interval
}

func (g *gitHubUpdater) newVersionAvailableContext(ctx context.Context, currentVersion string) (ok bool, url string, err error) {
	if g.newVersionAvailable && g.fetchURL != "" {
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

	ok, err = compareVersions(tag, currentVersion)
	if err != nil {
		return false, "", fmt.Errorf("version formats are invalid %s, %s: %w",
			currentVersion, tag, err)
	}
	if !ok {
		// if current is latest, check next time after interval
		return false, "", t.checkout()
	}
	// new version available
	d := func() {
		g.fetchURL = url
		g.newVersionAvailable = true
	}
	d()
	return true, url, nil
}

func getLatestAssetInfo(ctx context.Context, client RepositoriesService, owner, repo string) (_, _ string, _ error) {
	latestRelease, _, err := client.GetLatestRelease(ctx, owner, repo)
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
	return "", "", errors.New("found no alfredworkflow assets")
}
