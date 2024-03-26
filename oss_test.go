package infra_test

import (
	"context"
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
	testBucketName := "test"
	testFilename := "test.txt"
	t.Run("Test Save File", func(t *testing.T) {
		testReader := strings.NewReader("Hello, World!")
		err = aliOss.Save(ctx, testBucketName, testFilename, testReader)
		require.NoError(t, err)
	})
}
