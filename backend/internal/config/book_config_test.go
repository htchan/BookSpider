package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadBookConfigs(t *testing.T) {
	t.Parallel()

	t.Cleanup(func() {
		os.Remove("book_test")
	})

	tests := []struct {
		name           string
		yamlBookString string
		want           map[string]*BookConfig
		wantErr        bool
	}{
		{
			name: "load book config success",
			yamlBookString: `testSourceKey:
        url:
          base: base url
          download: download url
          chapterPrefix: chapter prefix url
        maxChapterError: 1
        updateDateLayout: layout`,
			want: map[string]*BookConfig{
				"testSourceKey": &BookConfig{
					URLConfig:        URLConfig{Base: "base url", Download: "download url", ChapterPrefix: "chapter prefix url"},
					MaxChaptersError: 1,
					UpdateDateLayout: "layout",
					SourceKey:        "testSourceKey",
				},
			},
			wantErr: false,
		},
	}

	for i, test := range tests {
		dirName := fmt.Sprintf("book_test/%d", i)
		os.MkdirAll(dirName, 0750)
		os.WriteFile(fmt.Sprintf("./%v/book_configs.yaml", dirName), []byte(test.yamlBookString), 0644)
		t.Run(test.name, func(t *testing.T) {
			t.Cleanup(func() {
				os.Remove(fmt.Sprintf("./%v/book_configs.yaml", dirName))
				os.Remove(dirName)
			})

			got, err := LoadBookConfigs(dirName)
			if (err != nil) != test.wantErr {
				t.Errorf("LoadBookConfigs() return error %v, wantErr %v", err, test.wantErr)
			}

			if !cmp.Equal(got, test.want) {
				t.Error(cmp.Diff(got, test.want))
			}
		})
	}
}
