package model

import (
	"fmt"
	"time"
)

type ApiRecord struct {
	Base
	Path       string
	Method     string
	Ip         string
	Latency    time.Duration
	Body       string
	StatusCode int
}

func (ApiRecord) TableName() string {
	return "api_records"
}

func (model *ApiRecord) String() string {
	return fmt.Sprintf("[%s %s %3d] from %s by %d in %f sec %s",
		model.Method,
		model.Path,
		model.StatusCode,
		model.Ip,
		model.CreatedBy,
		model.Latency.Seconds(),
		model.Body,
	)
}
