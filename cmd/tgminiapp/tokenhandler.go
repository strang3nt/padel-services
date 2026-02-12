package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type data struct {
	Expiry time.Time
	Blob   any
}

type FileTokenHandler struct {
	store map[string]data
	mutex sync.Mutex
}

func makeFileTokenHandler() FileTokenHandler {
	return FileTokenHandler{
		store: map[string]data{},
	}
}

func (fm *FileTokenHandler) GenerateToken(
	duration time.Duration,
	blob any,
) string {
	b := make([]byte, 16)
	rand.Read(b) //nolint:all
	token := hex.EncodeToString(b)

	fm.mutex.Lock()
	fm.store[token] = data{
		Expiry: time.Now().Add(duration),
		Blob:   blob,
	}
	fm.mutex.Unlock()

	time.AfterFunc(duration, func() {
		fm.mutex.Lock()
		delete(fm.store, token)
		fm.mutex.Unlock()
	})

	return token
}

func (fm *FileTokenHandler) GetBlob(token string) (any, bool) {
	if data, ok := fm.store[token]; ok && time.Now().Before(data.Expiry) {
		fm.mutex.Lock()
		data := data.Blob
		fm.mutex.Unlock()
		return data, true
	}

	return "", false
}
