package networking_demo

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"

	"time"
)

func RestyDemo(checkUrl string) {
	// 创建一个Resty客户端
	client := resty.New()
	// 禁用HTTP/2
	client.SetTransport(&http.Transport{
		ForceAttemptHTTP2: false,
	})
	suitsArray := []uint16{
		tls.TLS_RSA_WITH_RC4_128_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
	}

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 使用Perm函数打乱数组顺序
	rand.Shuffle(len(suitsArray), func(i, j int) {
		suitsArray[i], suitsArray[j] = suitsArray[j], suitsArray[i]
	})
	suitsArray = suitsArray[:10]

	fmt.Print(suitsArray)
	// 关闭证书校验
	client.SetTLSClientConfig(&tls.Config{
		CipherSuites: suitsArray,

		InsecureSkipVerify: true})

	// 发送GET请求
	resp, err := client.R().Get(checkUrl)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}

	// 检查响应状态码
	if resp.StatusCode() != 200 {
		fmt.Println("请求失败，状态码:", resp.StatusCode())
		return
	}

	// 获取响应内容
	body := resp.Body()
	fmt.Println("restyDemo info start  ================ ")

	fmt.Println("响应内容:", string(body))
	fmt.Println("restyDemo info end  ================ ")
}
