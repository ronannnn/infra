package officialaccount

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/url"
	"sort"
)

// signJsSdkConfig 生成 JS-SDK 配置签名
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
// 对以下内容进行 SHA1 签名
// noncestr=Wm3WZYTPz0wzccnW
// jsapi_ticket=sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg
// timestamp=1414587457
// url=http://mp.weixin.qq.com?params=value
func signJsSdkConfig(jsApiTicket, nonceStr string, timestamp int64, signedUrl string) string {
	params := url.Values{
		"jsapi_ticket": {jsApiTicket},
		"noncestr":     {nonceStr},
		"timestamp":    {fmt.Sprintf("%d", timestamp)},
		"url":          {signedUrl},
	}
	return sign(params)
}

func sign(params url.Values) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(params.Get(key))
		buf.WriteByte('&')
	}
	buf.Truncate(buf.Len() - 1)
	return fmt.Sprintf("%x", sha1.Sum(buf.Bytes()))
}
