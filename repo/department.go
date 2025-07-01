package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewDepartmentRepo() service.DepartmentRepo {
	return &departmentRepo{
		DefaultCrudRepo: NewDefaultCrudRepo[model.Department]("Company", "Leader", "Parent"),
	}
}

type departmentRepo struct {
	DefaultCrudRepo[model.Department]
}
