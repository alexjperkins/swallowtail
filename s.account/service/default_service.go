package service

import (
	"context"
	"swallowtail/s.account/dao"
	"swallowtail/s.account/domain"
)

type accountService struct{}

func (a *accountService) CreateAccount(ctx context.Context) error {
	account := &domain.Account{}
	return dao.CreateAccount(account)
}
