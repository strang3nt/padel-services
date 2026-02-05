// 1. Telegram Verification Logic
package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func verifyTelegram(initData string) (int64, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, err
	}

	receivedHash := values.Get("hash")
	values.Del("hash")

	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheck []string
	for _, k := range keys {
		dataCheck = append(dataCheck, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataStr := strings.Join(dataCheck, "\n")

	// HMAC-SHA256 chain
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secret := h.Sum(nil)

	h = hmac.New(sha256.New, secret)
	h.Write([]byte(dataStr))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	if expectedHash != receivedHash {
		return 0, fmt.Errorf("invalid hash")
	}

	// Extract User ID
	var user struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(values.Get("user")), &user); err != nil {
		return 0, err
	}

	return user.ID, nil
}
