package client

import (
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type DecodeMethod string

const (
	DecodeMethodGBK  DecodeMethod = "gbk"
	DecodeMethodBig5 DecodeMethod = "big5"
	DecodeMethodUTF8 DecodeMethod = "utf8"
)

type Decoder struct {
	decoder *encoding.Decoder
}

func NewDecoder(decodeMethod DecodeMethod) Decoder {
	var decoder *encoding.Decoder
	switch decodeMethod {
	case DecodeMethodGBK:
		decoder = simplifiedchinese.GBK.NewDecoder()
	case DecodeMethodBig5:
		decoder = traditionalchinese.Big5.NewDecoder()
	default:
		decoder = nil
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
