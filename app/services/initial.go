package service

import "context"

type Service struct {
	ShutdownCtx context.Context
}

var AppService *Service

func InitializeData(ctx context.Context) error {
	if err := initEnvironment(); err != nil {
		return err
	}

	AppService = &Service{
		ShutdownCtx: ctx,
	}
	return nil
}

func initEnvironment() error {
	return nil
}
