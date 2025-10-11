package usecase

import (
	"safe_talk/internal/repository"
	"safe_talk/pkg/logger"
)

type UseCase struct {
	l logger.Logger
	r repository.Repository
}

func New(l logger.Logger, r repository.Repository) UseCase {
	return UseCase{l: l, r: r}
}
