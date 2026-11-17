package models

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Session struct {
	ID       string
	UserID   int
	ExpireAt time.Time
}

var sessions = make(map[string]Session)
var mut sync.Mutex

func NewSession(userID int) (string, error) {
	mut.Lock()
	defer mut.Unlock()

	randByte := make([]byte, 32)

	if _, err := rand.Read(randByte); err != nil {
		return "", err
	}

	id := hex.EncodeToString(randByte)
	expireAt := time.Now().Add(24 * time.Hour)

	newSession := Session{
		ID:       id,
		UserID:   userID,
		ExpireAt: expireAt,
	}

	sessions[id] = newSession

	return id, nil
}

func GetSession(id string) (Session, bool) {
	mut.Lock()
	defer mut.Unlock()

	s, ok := sessions[id]

	if !ok {
		return Session{}, false
	}

	if time.Now().After(s.ExpireAt) {
		delete(sessions, id)
		return Session{}, false
	}
	return s, true
}

func DeleteSession(id string) {
	mut.Lock()
	defer mut.Unlock()
	delete(sessions, id)
}
