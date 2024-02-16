//go:build wireinject
// +build wireinject

package di

import (
	respository "cleancode/pkg/respository"
	"cleancode/pkg/usecase"

	"github.com/google/wire"
)

func InitializeUserUseCase() (*usecase.UserUseCase, error) {
	wire.Build(
		usecase.NewUserUseCase,
		respository.NewUserRepository,
	)
	return &usecase.UserUseCase{}, nil
}
