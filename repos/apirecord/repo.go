package apirecord

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/repos"
	srv "github.com/ronannnn/infra/services/apirecord"
)

func ProvideStore() srv.Repo {
	return &repo{
		repos.NewDefaultCrudRepo[models.ApiRecord](),
	}
}

type repo struct {
	repos.DefaultCrudRepo[models.ApiRecord]
}
