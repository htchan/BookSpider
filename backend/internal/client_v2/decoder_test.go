package client

import (
	"encoding/hex"
	"testing"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
)

func Test_NewDecoder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		decodeMethod  string
		expectDecoder Decoder
	}{
		{
			name:          "load big5 decoder",
			decodeMethod:  "big5",
			expectDecoder: Decoder{decoder: traditionalchinese.Big5.NewDecoder()},
		},
		{
			name:          "load gbk decoder",
			decodeMethod:  "gbk",
			expectDecoder: Decoder{decoder: simplifiedchinese.GBK.NewDecoder()},
		},
		{
			name:          "load nil decoder",
			decodeMethod:  "",
			expectDecoder: Decoder{},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			decoder := NewDecoder(test.decodeMethod)
			if !((decoder.decoder == nil && test.expectDecoder.decoder == nil) ||
				(decoder.decoder != nil && test.expectDecoder.decoder != nil && *decoder.decoder == *test.expectDecoder.decoder)) {
				t.Errorf("got decoder.decoder: %v; expect decoder: %v", decoder, test.expectDecoder)
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
			decoder:    NewDecoder("big5"),
			inputBytes: "a440",
			want:       "一",
			wantErr:    false,
		},
		{
			name:       "decode string with nil decoder",
			decoder:    NewDecoder(""),
			inputBytes: "41",
			want:       "A",
			wantErr:    false,
		},
		{
			name:       "decode gbk string with gbk decoder",
			decoder:    NewDecoder("gbk"),
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
