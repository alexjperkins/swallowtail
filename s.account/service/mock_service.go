package service

import "context"

type mockService struct{}

func (m *mockService) CreateAccount(ctx context.Context) error { return nil }
