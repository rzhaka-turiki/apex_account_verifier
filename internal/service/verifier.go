package service

import (
	"context"
	"fmt"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/api"
	"github.com/rzhaka-turiki/apex_account_verifier/internal/model"
)

type Verifier struct {
	client *api.Client
}

func NewVerifier(client *api.Client) *Verifier {
	return &Verifier{client: client}
}

func (v *Verifier) VerifyAccount(ctx context.Context, player, platform string, level int) (*model.Account, error) {
	account, err := v.client.GetAccount(ctx, player, platform)
	if err != nil {
		return nil, err
	}
	if account.Level != level {
		return nil, fmt.Errorf("level mismatch: expected %d, got %d", level, account.Level)
	}
	return account, nil
}
