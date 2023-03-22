package external_services

import (
	"context"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"niromash-api/api"
	"niromash-api/utils/environment"
	"sync"
)

var _ api.GithubService = (*githubService)(nil)

type githubService struct {
	service *ExternalServicesService
	client  *github.Client
}

func newGithubService() *githubService {
	return &githubService{}
}

func (g *githubService) Init(service *ExternalServicesService) error {
	g.service = service
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: environment.GetGithubToken()},
	)
	g.client = github.NewClient(oauth2.NewClient(context.Background(), ts))

	return nil
}

func (g *githubService) ListRepositoriesConcurrent(callback api.ListRepositoriesCallback) (err error) {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var wg sync.WaitGroup
	for {
		wg.Add(1)
		var (
			repos []*github.Repository
			resp  *github.Response
		)

		repos, resp, err = g.client.Repositories.List(context.Background(), "", opt)
		if err != nil {
			return
		}
		go callback(repos, &wg)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	wg.Wait()
	return
}
