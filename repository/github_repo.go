package repository

import (
	"backend-github-trending/model"
	"golang.org/x/net/context"
)

type GithubRepo interface {
	SaveRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error)
	SelectRepos(context context.Context, userId string, limit int) ([]model.GithubRepo, error)
	SelectRepoByName(context context.Context, name string) (model.GithubRepo, error)
	UpdateRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error)
}
