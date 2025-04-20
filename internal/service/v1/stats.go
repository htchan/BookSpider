package service

import (
	"context"
	"database/sql"

	"github.com/htchan/BookSpider/internal/repo"
)

func (s *ServiceImpl) Stats(ctx context.Context) repo.Summary {
	return s.rpo.Stats()
}

func (s *ServiceImpl) DBStats(ctx context.Context) sql.DBStats {
	return s.rpo.DBStats()
}
