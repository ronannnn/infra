package loginrecord

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/repos"
	srv "github.com/ronannnn/infra/services/loginrecord"
)

func New() srv.Repo {
	return &repo{
		repos.NewDefaultCrudRepo[models.LoginRecord](),
	}
}

type repo struct {
	repos.DefaultCrudRepo[models.LoginRecord]
}
