package book

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/htchan/BookSpider/internal/model"
)

func Info(bk model.Book) string {
	bytes, err := json.Marshal(bk)
	if err != nil {
		return fmt.Sprintf("%s-%v#%v", bk.Site, bk.ID, strconv.FormatInt(int64(bk.HashCode), 36))
	}
	return string(bytes)
}
