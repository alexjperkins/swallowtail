package service

import "context"

var (
	svc AccountService
)

func Init() error {
	if svc != nil {
		return nil
	}
	svc = &accountService{}
	return nil
}

func UseMock() {}

type AccountService interface {
	CreateAccount(context.Context) error
}

func CreateAccount(ctx context.Context) error { return svc.CreateAccount(ctx) }
