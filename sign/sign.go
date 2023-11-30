package sign

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Generate 生成签名
func Generate(timestamp, nonce, secretKey string) string {
	strs := []string{timestamp, nonce, secretKey}
	sort.Strings(strs)

	h := sha1.New()
	h.Write([]byte(strings.Join(strs, "")))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Verify 验证签名
func Verify(timestamp, nonce, signature, secretKey string, expired int64) error {
	now := time.Now().Unix()

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return err
	}
	if ts > now {
		return errors.New("sign generation time is greater than validation time")
	}

	if now - ts >= expired {
		return errors.New("sign expired")
	}

	// 校验签名
	if Generate(timestamp, nonce, secretKey) != signature {
		return errors.New("sign verify fail")
	}

	return nil
}
