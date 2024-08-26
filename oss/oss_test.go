package oss_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/oss"
	"github.com/stretchr/testify/require"
)

type testCfg struct {
	oss oss.Cfg
}

func TestAliOss(t *testing.T) {
	var err error
	testCfg := testCfg{}
	err = cfg.ReadFromFile("configs/config.aliosstest.toml", &testCfg)
	require.NoError(t, err)
	// init service
	var aliOss oss.AliOss
	aliOss, err = oss.NewAliOss(&testCfg.oss)
	require.NoError(t, err)
	ctx := context.Background()
	testFilename := "test.txt"
	t.Run("Test Save File", func(t *testing.T) {
		testReader := strings.NewReader("Hello, World!")
		err = aliOss.Save(ctx, testCfg.oss.RootBucket, testFilename, testReader)
		require.NoError(t, err)
	})
	t.Run("Test Get File", func(t *testing.T) {
		var rc io.ReadCloser
		rc, err = aliOss.Get(ctx, testCfg.oss.RootBucket, testFilename)
		require.NoError(t, err)
		// convert io.ReadCloser to string
		buf := new(strings.Builder)
		_, err = io.Copy(buf, rc)
		require.NoError(t, err)
		require.Equal(t, "Hello, World!", buf.String())
	})
	t.Run("Test Delete File", func(t *testing.T) {
		err = aliOss.Delete(ctx, testCfg.oss.RootBucket, testFilename)
		require.NoError(t, err)
		_, err = aliOss.Get(ctx, testCfg.oss.RootBucket, testFilename)
		require.Error(t, err)
	})
	t.Run("Test Get File Upload url", func(t *testing.T) {
		var uploadUrl string
		uploadUrl, err = aliOss.GetUploadUrl(ctx, testCfg.oss.RootBucket, testFilename)
		println(uploadUrl)
		require.NotEmpty(t, uploadUrl)
		require.NoError(t, err)
	})
	t.Run("Test Get File Download url", func(t *testing.T) {
		var downloadUrl string
		downloadUrl, err = aliOss.GetDownloadUrl(ctx, testCfg.oss.RootBucket, testFilename)
		println(downloadUrl)
		require.NotEmpty(t, downloadUrl)
		require.NoError(t, err)
	})
}
