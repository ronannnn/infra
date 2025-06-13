package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewDepartmentRepo() service.DepartmentRepo {
	return &departmentRepo{}
}

type departmentRepo struct {
	DefaultCrudRepo[model.Department]
}
