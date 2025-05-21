package repo_impl

import (
	"backend-github-trending/banana"
	"backend-github-trending/db"
	"backend-github-trending/log"
	"backend-github-trending/model"
	"backend-github-trending/repository"
	"context"
	"database/sql"
	"github.com/lib/pq"
	"time"
)

type GithubRepoImpl struct {
	sql *db.Sql // connect Database
}

func NewGithubRepo(sql *db.Sql) repository.GithubRepo {
	return &GithubRepoImpl{
		sql: sql,
	}
}

func (g GithubRepoImpl) SelectRepoByName(context context.Context, name string) (model.GithubRepo, error) {
	repo := model.GithubRepo{}
	err := g.sql.Db.GetContext(context, &repo, "select * from github_repo where name = $1", name)

	if err != nil {
		if err == sql.ErrNoRows {
			return repo, banana.RepoNotFound
		}
		log.Error(err.Error())
		return repo, err
	}

	return repo, nil
}

func (g GithubRepoImpl) SaveRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	statement := `insert into repos(name, description, url, color, lang, fork, stars,stars_today, build_by, created_at, updated_at)
				  value(:name, :description, :url,:color,:lang,:fork,:stars,:stars_today,:build_by,:created_at,:updated_at)`

	repo.CreatedAt = time.Now()
	repo.UpdatedAt = time.Now()

	_, err := g.sql.Db.NamedExecContext(context, statement, &repo)
	if err != nil {
		if err, ok := err.(*pq.Error); ok {
			if err.Code == "23505" {
				return repo, banana.RepoConflict
			}
		}
		log.Error(err.Error())
		return repo, banana.RepoInsertFalse
	}

	return repo, nil
}

func (g GithubRepoImpl) SelectRepos(context context.Context, userId string, limit int) ([]model.GithubRepo, error) {
	repos := []model.GithubRepo{}
	err := g.sql.Db.SelectContext(context, &repos,
		`
			SELECT 
				repos.name, repos.description, repos.url, repos.color, repos.lang, 
				repos.fork, repos.stars, repos.stars_today, repos.build_by, repos.updated_at, 
				COALESCE(repos.name = bookmarks.repo_name, FALSE) as bookmarked
			FROM repos
			FULL OUTER JOIN bookmarks 
			ON repos.name = bookmarks.repo_name AND 
			   bookmarks.user_id=$1  
			WHERE repos.name IS NOT NULL 
			ORDER BY updated_at ASC LIMIT $2
		`, userId, limit)

	if err != nil {
		log.Error(err.Error())
		return repos, err
	}
	return repos, nil
}

func (g GithubRepoImpl) UpdateRepo(context context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	// name, description, url, color,lang, fork, stars, stars_today, build_by, created_at, updated_at
	statement := `update repos 
				 set 
				     stars =  :stars,
				     fork =  :fork,
				     stars_today = :stars_today,
				     build_by = :build_by,
				     updated_at = :updated_at
				where name = :name`
	result, err := g.sql.Db.NamedExecContext(context, statement, &repo)
	if err != nil {
		log.Error(err.Error())
		return repo, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		log.Error(err.Error())
		return repo, banana.RepoNotUpdated
	}

	if count == 0 {
		return repo, banana.RepoNotUpdated
	}

	return repo, nil
}
