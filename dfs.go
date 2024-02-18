package infra

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type DfsCfg struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id"`
	SecretAccessKey string `mapstructure:"secret-access-key"`
	Location        string `mapstructure:"location"`
	RootBucket      string `mapstructure:"root-bucket"`
	Secure          bool   `mapstructure:"secure"`
}

type Dfs interface {
	Save(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64) error
	Delete(ctx context.Context, bucketName, objectName string) error
	Get(ctx context.Context, bucketName, objectName string) (*minio.Object, error)
	GetDownloadUrl(ctx context.Context, bucketName, objectName string) (string, error)
	ListInfo(ctx context.Context, bucketName string) ([]minio.ObjectInfo, error)
	DeleteBucket(ctx context.Context, bucketName string) error
}

func NewDfs(
	cfg DfsCfg,
) (Dfs, error) {
	if dfs, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.SecretAccessKey, ""),
		Secure: cfg.Secure,
	}); err != nil {
		return nil, err
	} else {
		return &DfsImpl{
			cfg:    cfg,
			dfsCli: dfs,
		}, nil
	}
}

type DfsImpl struct {
	cfg    DfsCfg
	dfsCli *minio.Client
}

func (srv *DfsImpl) Save(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64) (err error) {
	// create bucket if not exists
	var exists bool
	if exists, err = srv.dfsCli.BucketExists(ctx, bucketName); err != nil {
		return
	} else if !exists {
		if err = srv.dfsCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: srv.cfg.Location}); err != nil {
			return
		}
	}
	_, err = srv.dfsCli.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return
}

func (srv *DfsImpl) Delete(ctx context.Context, bucketName, objectName string) error {
	return srv.dfsCli.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (srv *DfsImpl) Get(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return srv.dfsCli.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (srv *DfsImpl) GetDownloadUrl(ctx context.Context, bucketName, objectName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))
	if presignedUrl, err := srv.dfsCli.PresignedGetObject(ctx, bucketName, objectName, time.Second*60*10, reqParams); err != nil {
		return "", err
	} else {
		return presignedUrl.String(), nil
	}
}

func (srv *DfsImpl) ListInfo(ctx context.Context, bucketName string) ([]minio.ObjectInfo, error) {
	objectInfoList := make([]minio.ObjectInfo, 0)
	objectInfoChan := srv.dfsCli.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for objectInfo := range objectInfoChan {
		if objectInfo.Err != nil {
			return nil, objectInfo.Err
		}
		objectInfoList = append(objectInfoList, objectInfo)
	}
	return objectInfoList, nil
}

func (srv *DfsImpl) DeleteBucket(ctx context.Context, bucketName string) error {
	return srv.dfsCli.RemoveBucket(ctx, bucketName)
}
