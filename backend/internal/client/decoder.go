package client

import (
	"github.com/htchan/BookSpider/internal/config"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Decoder struct {
	config.DecoderConfig
	decoder *encoding.Decoder
}

func (decoder *Decoder) Load() {
	if decoder.Method == "big5" {
		decoder.decoder = traditionalchinese.Big5.NewDecoder()
	} else if decoder.Method == "gbk" {
		decoder.decoder = simplifiedchinese.GBK.NewDecoder()
	}
}

func (decoder Decoder) Decode(str string) (string, error) {
	if decoder.decoder == nil {
		return str, nil
	}
	str, _, err := transform.String(decoder.decoder, str)
	return str, err
}
