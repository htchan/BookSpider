package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/rs/zerolog"
)

func (s *ServiceImpl) BookInfo(ctx context.Context, bk *model.Book) string {
	bytes, err := json.Marshal(bk)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to marshal book info")

		return fmt.Sprintf("%s-%v#%v", bk.Site, bk.ID, bk.FormatHashCode())
	}

	return string(bytes)
}
