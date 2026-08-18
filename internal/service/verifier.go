package service

import (
	"context"
	"fmt"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/api"
	"github.com/rzhaka-turiki/apex_account_verifier/internal/model"
	"github.com/rzhaka-turiki/apex_account_verifier/internal/queue"
)

type Verifier struct {
	client *api.Client
	queue  *queue.Queue
}

func NewVerifier(client *api.Client, q *queue.Queue) *Verifier {
	return &Verifier{
		client: client,
		queue:  q,
	}
}

func (v *Verifier) VerifyAccount(ctx context.Context, player, platform string, level int) (*model.Account, error) {
	account, err := v.queue.Submit(
		ctx,
		func(ctx context.Context) (*model.Account, error) {
			return v.client.GetAccount(ctx, player, platform)
		},
	)
	if err != nil {
		return nil, err
	}
	if account.Level != level {
		return nil, fmt.Errorf("level mismatch: expected %d, got %d", level, account.Level)
	}
	return account, nil
}
