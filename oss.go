package infra

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/ronannnn/infra/cfg"
)

type Minioss interface {
	Save(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64) error
	Delete(ctx context.Context, bucketName, objectName string) error
	Get(ctx context.Context, bucketName, objectName string) (*minio.Object, error)
	GetDownloadUrl(ctx context.Context, bucketName, objectName string) (string, error)
	GetUploadUrl(ctx context.Context, bucketName, objectName string) (string, error)
	ListInfo(ctx context.Context, bucketName string) ([]minio.ObjectInfo, error)
}

func NewMinioss(
	cfg *cfg.Dfs,
) (Minioss, error) {
	if client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyId, cfg.AccessKeySecret, ""),
		Secure: cfg.Secure,
	}); err != nil {
		return nil, err
	} else {
		return &MiniossImpl{
			cfg:      cfg,
			minioCli: client,
		}, nil
	}
}

type MiniossImpl struct {
	cfg      *cfg.Dfs
	minioCli *minio.Client
}

func (srv *MiniossImpl) Save(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64) (err error) {
	// create bucket if not exists
	var exists bool
	if exists, err = srv.minioCli.BucketExists(ctx, bucketName); err != nil {
		return
	} else if !exists {
		if err = srv.minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: srv.cfg.Location}); err != nil {
			return
		}
	}
	_, err = srv.minioCli.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	return
}

func (srv *MiniossImpl) Delete(ctx context.Context, bucketName, objectName string) error {
	return srv.minioCli.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (srv *MiniossImpl) Get(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return srv.minioCli.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (srv *MiniossImpl) GetDownloadUrl(ctx context.Context, bucketName, objectName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))
	if presignedUrl, err := srv.minioCli.PresignedGetObject(ctx, bucketName, objectName, time.Second*time.Duration(srv.cfg.ExpiredInSec), reqParams); err != nil {
		return "", err
	} else {
		return presignedUrl.String(), nil
	}
}

func (srv *MiniossImpl) GetUploadUrl(ctx context.Context, bucketName, objectName string) (string, error) {
	if presignedUrl, err := srv.minioCli.PresignedPutObject(ctx, bucketName, objectName, time.Second*time.Duration(srv.cfg.ExpiredInSec)); err != nil {
		return "", err
	} else {
		return presignedUrl.String(), nil
	}
}

func (srv *MiniossImpl) ListInfo(ctx context.Context, bucketName string) ([]minio.ObjectInfo, error) {
	objectInfoList := make([]minio.ObjectInfo, 0)
	objectInfoChan := srv.minioCli.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for objectInfo := range objectInfoChan {
		if objectInfo.Err != nil {
			return nil, objectInfo.Err
		}
		objectInfoList = append(objectInfoList, objectInfo)
	}
	return objectInfoList, nil
}

func (srv *MiniossImpl) DeleteBucket(ctx context.Context, bucketName string) error {
	return srv.minioCli.RemoveBucket(ctx, bucketName)
}

type AliOss interface {
	Save(ctx context.Context, bucketName, objectName string, reader io.Reader) error
	Delete(ctx context.Context, bucketName, objectName string) error
	Get(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	GetDownloadUrl(ctx context.Context, bucketName, objectName string) (string, error)
	GetUploadUrl(ctx context.Context, bucketName, objectName string) (string, error)
}

func NewAliOss(
	cfg *cfg.Dfs,
) (AliOss, error) {
	if client, err := oss.New(cfg.Endpoint, cfg.AccessKeyId, cfg.AccessKeySecret); err != nil {
		return nil, err
	} else {
		return &AliOssImpl{
			cfg:       cfg,
			aliOssCli: client,
		}, nil
	}
}

type AliOssImpl struct {
	cfg       *cfg.Dfs
	aliOssCli *oss.Client
}

func (srv *AliOssImpl) Save(ctx context.Context, bucketName, objectName string, reader io.Reader) (err error) {
	// create bucket if not exists
	var exists bool
	if exists, err = srv.aliOssCli.IsBucketExist(bucketName); err != nil {
		return
	} else if !exists {
		if err = srv.aliOssCli.CreateBucket(bucketName); err != nil {
			return
		}
	}
	var bucket *oss.Bucket
	if bucket, err = srv.aliOssCli.Bucket(bucketName); err != nil {
		return
	}
	err = bucket.PutObject(objectName, reader)
	return
}

func (srv *AliOssImpl) Delete(ctx context.Context, bucketName, objectName string) (err error) {
	var bucket *oss.Bucket
	if bucket, err = srv.aliOssCli.Bucket(bucketName); err != nil {
		return
	}
	return bucket.DeleteObject(objectName)
}

func (srv *AliOssImpl) Get(ctx context.Context, bucketName, objectName string) (rc io.ReadCloser, err error) {
	var bucket *oss.Bucket
	if bucket, err = srv.aliOssCli.Bucket(bucketName); err != nil {
		return
	} else {
		return bucket.GetObject(objectName)
	}
}

func (srv *AliOssImpl) GetDownloadUrl(ctx context.Context, bucketName, objectName string) (url string, err error) {
	var bucket *oss.Bucket
	if bucket, err = srv.aliOssCli.Bucket(bucketName); err != nil {
		return
	} else {
		return bucket.SignURL(objectName, oss.HTTPGet, srv.cfg.ExpiredInSec)
	}
}

func (srv *AliOssImpl) GetUploadUrl(ctx context.Context, bucketName, objectName string) (url string, err error) {
	var bucket *oss.Bucket
	if bucket, err = srv.aliOssCli.Bucket(bucketName); err != nil {
		return
	} else {
		return bucket.SignURL(objectName, oss.HTTPPut, srv.cfg.ExpiredInSec)
	}
}
