package huaweihilink

import "errors"

var (
	ErrSessionFileNotFound = errors.New("session file not found")
	ErrSessionExpired      = errors.New("session expired")
)
