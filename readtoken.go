package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// 只读令牌：让「数据展示」不再需要 CPA 管理密钥。
//
// 背景：CPA 对 /v0/management/* 的门禁是无状态 Bearer 校验——每次请求都要带
// 管理密钥，没有会话、没有 Cookie，插件也无法按路由绕过（中间件先于插件
// 执行）。而 /v0/resource/plugins/<id>/* 是免鉴权的 GET 通道，CPA 会把请求
// 动态分发给插件（ServeResourceHTTP，仅 GET）。
//
// 因此读数据改走资源路径，但不做成完全公开：由插件颁发一个只读令牌——
//   - 令牌 256bit 随机，生成后持久化在数据目录（0600），删除文件即重新生成；
//   - 只放行 GET 读取（用量汇总/明细/单价/额度/识别结果），不能写、不能改
//     任何设置——写操作仍在 /v0/management/* 后面，由 CPA 管理密钥把守；
//   - 浏览器首次进入时用管理密钥调 /read-token（该路由本身仍由 CPA 管理
//     密钥保护）换取令牌，存入 localStorage，之后数据展示不再需要密钥；
//   - 管理密钥本身不落浏览器盘：换到令牌后仅保留在 sessionStorage（当前
//     标签页），关闭即失效。
type readTokenState struct {
	Token string `json:"token"`
}

var (
	readTokenMu    sync.Mutex
	readTokenCache string
)

func readTokenPath() string {
	return filepath.Join(storeDir(), "stepfun-credit-read-token.json")
}

// currentReadToken 返回（必要时生成并持久化）只读令牌。
func currentReadToken() string {
	readTokenMu.Lock()
	defer readTokenMu.Unlock()
	if readTokenCache != "" {
		return readTokenCache
	}
	if raw, err := os.ReadFile(readTokenPath()); err == nil {
		var st readTokenState
		if json.Unmarshal(raw, &st) == nil && st.Token != "" {
			readTokenCache = st.Token
			return readTokenCache
		}
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	token := hex.EncodeToString(buf)
	if err := os.MkdirAll(storeDir(), 0o700); err != nil {
		return ""
	}
	raw, err := json.Marshal(readTokenState{Token: token})
	if err != nil {
		return ""
	}
	tmp := readTokenPath() + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return ""
	}
	if err := os.Rename(tmp, readTokenPath()); err != nil {
		return ""
	}
	readTokenCache = token
	return token
}

// readTokenOK 常量时间比较；空令牌一律拒绝。
func readTokenOK(provided string) bool {
	if provided == "" {
		return false
	}
	expected := currentReadToken()
	if expected == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) == 1
}
