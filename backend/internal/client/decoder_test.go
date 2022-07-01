package client

import (
	"encoding/hex"
	"github.com/htchan/BookSpider/internal/config"
	"golang.org/x/text/transform"
	"testing"
)

func TestDecoder_Load(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		decoderMethod string
		inputBytes    string
		want          string
	}{
		{
			name:          "load big5 decoder",
			decoderMethod: "big5",
			inputBytes:    "a440",
			want:          "一",
		},
		{
			name:          "load gbk decoder",
			decoderMethod: "gbk",
			inputBytes:    "d2bb",
			want:          "一",
		},
		{
			name:          "load nil decoder",
			decoderMethod: "",
			inputBytes:    "",
			want:          "",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			decoder := Decoder{DecoderConfig: config.DecoderConfig{Method: test.decoderMethod}}
			decoder.Load()
			if test.inputBytes == "" {
				if decoder.decoder != nil {
					t.Errorf("load non nil decoder")
				}
			} else {
				hexByte, _ := hex.DecodeString(test.inputBytes)

				got, _, err := transform.String(decoder.decoder, string(hexByte))

				if err != nil || got != test.want {
					t.Errorf("load wrong decoder decode %s to %s", got, test.want)
				}
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		decoderMethod string
		inputBytes    string
		want          string
		wantErr       bool
	}{
		{
			name:          "decode big5 string with big5 decoder",
			decoderMethod: "big5",
			inputBytes:    "a440",
			want:          "一",
			wantErr:       false,
		},
		{
			name:          "decode string with nil decoder",
			decoderMethod: "",
			inputBytes:    "41",
			want:          "A",
			wantErr:       false,
		},
		{
			name:          "decode big5 string with big5 decoder",
			decoderMethod: "gbk",
			inputBytes:    "d2bb",
			want:          "一",
			wantErr:       false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			decoder := Decoder{DecoderConfig: config.DecoderConfig{Method: test.decoderMethod}}
			decoder.Load()
			hexByte, _ := hex.DecodeString(test.inputBytes)
			got, err := decoder.Decode(string(hexByte))
			if (err != nil) != test.wantErr {
				t.Errorf("Decoder.Decode() return err %v, wantErr %v", err, test.wantErr)
			}
			if got != test.want {
				t.Errorf("Decoder.Decode() return %v; want: %v", got, test.want)
			}
		})
	}

	// t.Run("decode big5 string with big5 decoder", func (t *testing.T) {
	// 	t.Parallel()
	// 	decoder := Decoder{DecoderConfig: config.DecoderConfig{Method: "big5"}}
	// 	decoder.Load()
	// 	hexByte, _ := hex.DecodeString("a440")
	// 	result, err := decoder.Decode(string(hexByte))
	// 	if err != nil || result != "一" {
	// 		t.Errorf("decode wrong result: %v; error: %v", result, err)
	// 	}
	// })

	// t.Run("decode string with nil decoder", func (t *testing.T) {
	// 	t.Parallel()
	// 	decoder := Decoder{}
	// 	hexByte, _ := hex.DecodeString("a440")
	// 	result, err := decoder.Decode(string(hexByte))
	// 	if err != nil || result != string(hexByte) {
	// 		t.Errorf("decode wrong result: %v; error: %v", result, err)
	// 	}
	// })
}
