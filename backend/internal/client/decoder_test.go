package client

import (
	"encoding/hex"
	"testing"

	"github.com/htchan/BookSpider/internal/config"
)

func Test_NewDecoder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		conf      config.DecoderConfig
		expectNil bool
	}{
		{
			name:      "load big5 decoder",
			conf:      config.DecoderConfig{Method: "big5"},
			expectNil: false,
		},
		{
			name:      "load gbk decoder",
			conf:      config.DecoderConfig{Method: "gbk"},
			expectNil: false,
		},
		{
			name:      "load nil decoder",
			conf:      config.DecoderConfig{},
			expectNil: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			decoder := NewDecoder(test.conf)
			if (decoder.decoder == nil) != test.expectNil {
				t.Errorf("got decoder.decoder: %v; expect nil: %v", decoder.decoder, test.expectNil)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		decoder    Decoder
		inputBytes string
		want       string
		wantErr    bool
	}{
		{
			name:       "decode big5 string with big5 decoder",
			decoder:    NewDecoder(config.DecoderConfig{Method: "big5"}),
			inputBytes: "a440",
			want:       "一",
			wantErr:    false,
		},
		{
			name:       "decode string with nil decoder",
			decoder:    NewDecoder(config.DecoderConfig{Method: ""}),
			inputBytes: "41",
			want:       "A",
			wantErr:    false,
		},
		{
			name:       "decode gbk string with gbk decoder",
			decoder:    NewDecoder(config.DecoderConfig{Method: "gbk"}),
			inputBytes: "d2bb",
			want:       "一",
			wantErr:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			hexByte, _ := hex.DecodeString(test.inputBytes)
			got, err := test.decoder.Decode(string(hexByte))
			if (err != nil) != test.wantErr {
				t.Errorf("Decoder.Decode() return err %v, wantErr %v", err, test.wantErr)
			}
			if got != test.want {
				t.Errorf("Decoder.Decode() return %v; want: %v", got, test.want)
			}
		})
	}
}
