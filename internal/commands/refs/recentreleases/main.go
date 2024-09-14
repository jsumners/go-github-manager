package recentreleases

import (
	"fmt"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/google/go-github/v63/github"
	"github.com/jsumners/ghm/internal/app"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
)

var repoName string

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recent-releases",
		Short: "Find refs for most recent releases.",
		Long: heredoc.Doc(`
			Queries a repo, or all repos in an org, for the two most recent
			releases. The result is a list of repos and target tag refs
			in CSV format.
		`),
		RunE: func(*cobra.Command, []string) error {
			return runE(app, repoName)
		},
	}

	cmd.Flags().StringVarP(
		&app.Owner,
		"owner",
		"o",
		"",
		"The user/org containing the repo(s) to target.",
	)
	cmd.MarkFlagRequired("owner")

	cmd.Flags().StringVarP(
		&repoName,
		"repo-name",
		"r",
		"",
		"Specific repository to target.",
	)

	return cmd
}

func runE(app *app.App, repoName string) error {
	results := make([]string, 0)

	if repoName != "" {
		releases, err := getReleases(app, repoName)
		if err != nil {
			return err
		}
		commits := findTags(releases)
		results = append(results, fmt.Sprintf("%s, %s", repoName, commits))
	} else {
		repos, err := getRepos(app)
		if err != nil {
			return err
		}

		wg := sync.WaitGroup{}
		for _, repo := range repos {
			go func() {
				wg.Add(1)
				defer wg.Done()

				releases, err := getReleases(app, repo.GetName())
				if err != nil {
					io.WriteString(os.Stderr, fmt.Sprintf("%s: could not get releases: %s\n", repo.GetName(), err))
					return
				}
				if len(releases) == 0 {
					return
				}
				commits := findTags(releases)
				results = append(results, fmt.Sprintf("%s, %s", repo.GetName(), commits))
			}()
		}
		wg.Wait()

	}

	slices.Sort(results)

	fmt.Println("repo, current, previous")
	for _, result := range results {
		fmt.Println(result)
	}

	return nil
}

func findTags(releases []*github.RepositoryRelease) string {
	var currentMajor int
	var currentRelease *github.RepositoryRelease
	var previousRelease *github.RepositoryRelease
	for _, release := range releases {
		if currentMajor == 0 {
			currentRelease = release
			version := release.GetTagName()
			currentMajor = tagToMajor(version)
			continue
		}

		version := release.GetTagName()
		major := tagToMajor(version)
		if major < currentMajor {
			previousRelease = release
			break
		}
	}

	return fmt.Sprintf(
		"%s, %s",
		currentRelease.GetTagName(),
		previousRelease.GetTagName(),
	)
}

func tagToMajor(tag string) int {
	parts := strings.Split(tag, ".")
	parts[0] = strings.Replace(parts[0], "v", "", 1)
	return cast.ToInt(parts[0])
}

func getRepos(app *app.App) ([]*github.Repository, error) {
	listOpts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	repos, resp, err := app.Client.Repositories.ListByOrg(app.Context, app.Owner, listOpts)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("got invalid response: %w", err)
	}

	repos = slices.DeleteFunc(repos, func(r *github.Repository) bool {
		// We do not care about archived repos.
		return r.GetArchived() == true
	})

	return repos, nil
}

func getReleases(app *app.App, repoName string) ([]*github.RepositoryRelease, error) {
	listOpts := &github.ListOptions{PerPage: 100}

	releases := make([]*github.RepositoryRelease, 0)
	for {
		_releases, resp, err := app.Client.Repositories.ListReleases(app.Context, app.Owner, repoName, listOpts)
		if err != nil {
			return nil, err
		}
		releases = append(releases, _releases...)

		if resp.NextPage == 0 {
			break
		}
		listOpts.Page = resp.NextPage
	}

	releases = slices.DeleteFunc(releases, func(r *github.RepositoryRelease) bool {
		// Draft releases do not have `PublishedAt` fields.
		return r.GetDraft() == true
	})

	// Sort most recent to oldest, i.e. first element is the latest release.
	slices.SortFunc(releases, func(a, b *github.RepositoryRelease) int {
		if a.PublishedAt == nil {
			fmt.Println("wtf a:", a.GetURL())
			panic("a")
		}
		if b.PublishedAt == nil {
			fmt.Println("wtf b:", b.GetURL())
			panic("b")
		}
		if a.PublishedAt.Before(b.PublishedAt.Time) {
			return 1
		}
		if b.PublishedAt.Before(a.PublishedAt.Time) {
			return -1
		}
		return 0
	})

	return releases, nil
}
