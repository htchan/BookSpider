package client

import (
	"testing"
	"encoding/hex"
	"golang.org/x/text/transform"
	"github.com/htchan/BookSpider/internal/config"
)

func TestDecoder_Load(t *testing.T) {
	t.Parallel()

	t.Run("load big5 decoder", func (t *testing.T) {
		t.Parallel()
		decoder := Decoder{DecoderConfig: config.DecoderConfig{Method: "big5"}}
		decoder.Load()
		hexByte, _ := hex.DecodeString("a440")
		text, _, err := transform.String(decoder.decoder, string(hexByte))
		if err != nil || text != "一" {
			t.Errorf("load wrong decoder decode \\ua440 to %s", text)
		}
	})

	t.Run("load nil decoder", func (t *testing.T) {
		t.Parallel()
		decoder := Decoder{DecoderConfig: config.DecoderConfig{Method: ""}}
		decoder.Load()
		if decoder.decoder != nil {
			t.Errorf("load wrong decoder")
		}
	})
}

func TestDecoder_Decode(t *testing.T) {
	t.Parallel()

	t.Run("decode big5 string with big5 decoder", func (t *testing.T) {
		t.Parallel()
		decoder := Decoder{DecoderConfig: config.DecoderConfig{Method: "big5"}}
		decoder.Load()
		hexByte, _ := hex.DecodeString("a440")
		result, err := decoder.Decode(string(hexByte))
		if err != nil || result != "一" {
			t.Errorf("decode wrong result: %v; error: %v", result, err)
		}
	})

	t.Run("decode string with nil decoder", func (t *testing.T) {
		t.Parallel()
		decoder := Decoder{}
		hexByte, _ := hex.DecodeString("a440")
		result, err := decoder.Decode(string(hexByte))
		if err != nil || result != string(hexByte) {
			t.Errorf("decode wrong result: %v; error: %v", result, err)
		}
	})
}