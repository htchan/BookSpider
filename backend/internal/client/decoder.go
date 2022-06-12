package client

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"github.com/htchan/BookSpider/internal/config"
)

type Decoder struct {
	config.DecoderConfig
	decoder *encoding.Decoder
}

func (decoder *Decoder) Load() {
	if decoder.Method == "big5" {
		decoder.decoder = traditionalchinese.Big5.NewDecoder()
	}
}

func (decoder Decoder) Decode(str string) (string, error) {
	if decoder.decoder == nil {
		return str, nil
	}
	str, _, err :=  transform.String(decoder.decoder, str)
	return str, err
}