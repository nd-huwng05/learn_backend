package banana

import "errors"

var (
	RepoNotFound    = errors.New("Repo Not Found")
	RepoConflict    = errors.New("Repo Conflict")
	RepoInsertFalse = errors.New("Repo Insert False")
	RepoNotUpdated  = errors.New("Repo Not Updated")
)
