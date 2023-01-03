package client

import (
	"github.com/htchan/BookSpider/internal/config"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Decoder struct {
	decoder *encoding.Decoder
}

func NewDecoder(conf config.DecoderConfig) Decoder {
	var decoder *encoding.Decoder
	if conf.Method == "big5" {
		decoder = traditionalchinese.Big5.NewDecoder()
	} else if conf.Method == "gbk" {
		decoder = simplifiedchinese.GBK.NewDecoder()
	}

	return Decoder{decoder: decoder}
}

func NewDecoderV2(decodeMethod string) Decoder {
	var decoder *encoding.Decoder
	if decodeMethod == "big5" {
		decoder = traditionalchinese.Big5.NewDecoder()
	} else if decodeMethod == "gbk" {
		decoder = simplifiedchinese.GBK.NewDecoder()
	}

	return Decoder{decoder: decoder}
}

func (decoder Decoder) Decode(str string) (string, error) {
	if decoder.decoder == nil {
		return str, nil
	}
	str, _, err := transform.String(decoder.decoder, str)
	return str, err
}
