package api

import (
	"github.com/ronannnn/infra/repos"
	srv "github.com/ronannnn/infra/services/api"

	"github.com/ronannnn/infra/models"
)

func New() srv.Repo {
	return &repo{
		repos.NewDefaultCrudRepo[models.Api](),
	}
}

type repo struct {
	repos.DefaultCrudRepo[models.Api]
}
