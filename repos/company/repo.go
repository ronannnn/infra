package company

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/repos"
	srv "github.com/ronannnn/infra/services/company"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
	repos.DefaultCrudRepo[*models.Company]
}
