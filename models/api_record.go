package models

import (
	"fmt"
	"time"
)

// api record

type ApiRecord struct {
	Base
	Path       string
	Method     string
	Ip         string
	Latency    time.Duration
	Body       string
	StatusCode int
}

func (model *ApiRecord) String() string {
	return fmt.Sprintf("[%s %s %3d] from %s in %f sec [%s]",
		model.Method,
		model.Path,
		model.StatusCode,
		model.Ip,
		model.Latency.Seconds(),
		model.Body,
	)
}
