package department

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/repos"
	srv "github.com/ronannnn/infra/services/department"
)

func New() srv.Repo {
	return &repo{
		repos.NewDefaultCrudRepo[*models.Department](),
	}
}

type repo struct {
	repos.DefaultCrudRepo[*models.Department]
}
