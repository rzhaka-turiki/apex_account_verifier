package api

import "context"

type ApexAPIClinet interface {
	GetUIDByName(ctx context.Context, playerName string, platform string) (string, error)
	GetPlayer(ctx context.Context, uid string, platform string) (*model.ApexPlayer, error)
}
