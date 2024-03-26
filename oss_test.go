package infra_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/ronannnn/infra"
	"github.com/ronannnn/infra/cfg"
	"github.com/stretchr/testify/require"
)

func TestAliOss(t *testing.T) {
	var err error
	testCfg := cfg.Cfg{}
	err = cfg.ReadFromFile("configs/config.aliosstest.toml", &testCfg)
	require.NoError(t, err)
	// init service
	var aliOss infra.AliOss
	aliOss, err = infra.NewAliOss(&testCfg.Dfs)
	require.NoError(t, err)
	ctx := context.Background()
	testBucketName := "gxphoto"
	testFilename := "test.txt"
	t.Run("Test Save File", func(t *testing.T) {
		testReader := strings.NewReader("Hello, World!")
		err = aliOss.Save(ctx, testBucketName, testFilename, testReader)
		require.NoError(t, err)
	})
	t.Run("Test Get File", func(t *testing.T) {
		var rc io.ReadCloser
		rc, err = aliOss.Get(ctx, testBucketName, testFilename)
		require.NoError(t, err)
		// convert io.ReadCloser to string
		buf := new(strings.Builder)
		_, err = io.Copy(buf, rc)
		require.NoError(t, err)
		require.Equal(t, "Hello, World!", buf.String())
	})
	t.Run("Test Delete File", func(t *testing.T) {
		err = aliOss.Delete(ctx, testBucketName, testFilename)
		require.NoError(t, err)
		_, err = aliOss.Get(ctx, testBucketName, testFilename)
		require.Error(t, err)
	})
}
