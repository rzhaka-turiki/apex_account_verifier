package api

import "errors"

var (
	ErrTryAgain        = errors.New("try again in a few minutes")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrPlayerNotFound  = errors.New("player not found")
	ErrExternalAPI     = errors.New("external api error")
	ErrUnknownPlatform = errors.New("unknown platform")
	ErrRateLimited     = errors.New("rate limit reached")
	ErrInternal        = errors.New("external api internal errror")
)
