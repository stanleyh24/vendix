package webhooks

import (
	"github.com/stanleyh24/vendix/internal/config"
	"github.com/stanleyh24/vendix/internal/database"
)

type Service struct {
	db  *database.DB
	cfg *config.Config
}

func NewService(db *database.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}
