package userdb

import (
	"fmt"
	"strings"
	"sync"
)

type UserInfo struct {
	Username string
	Password string

	// TBC
}

type DB interface {
	SyncFromDB() error
	GetUserInfo(user string) (*UserInfo, bool)
}

type DBBase struct {
	DBType string
	DBPath string

	UserDB map[string]*UserInfo
	DBLock sync.RWMutex
}

func NewDB(dbpath string) (DB, error) {
	if strings.HasPrefix(dbpath, "file://") {
		return NewFileDB(dbpath);
	}
	return nil, fmt.Errorf("dbpath %s does not support.", dbpath)
}
