package bookprocess

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Test_ParamsFromData(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		data         []byte
		expectParams *BookProcessParams
		expectErr    error
	}{
		{
			name: "happy flow/with trace",
			data: []byte(`{"book":{"site":"test","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"01234567890123456789012345678901","span_id":"0123456789012345","trace_flags":1}`),
			expectParams: &BookProcessParams{
				Book: model.Book{
					Site:          "test",
					ID:            1,
					HashCode:      0,
					Title:         "title",
					Writer:        model.Writer{ID: 1, Name: "writer"},
					Type:          "type",
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Status:        model.StatusInProgress,
					IsDownloaded:  true,
					Error:         model.Error{Err: errors.New("error")},
				},
				TraceID:    "01234567890123456789012345678901",
				SpanID:     "0123456789012345",
				TraceFlags: 0x1,
			},
			expectErr: nil,
		},
		{
			name: "happy flow/without trace",
			data: []byte(`{"book":{"site":"test","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"}}`),
			expectParams: &BookProcessParams{
				Book: model.Book{
					Site:          "test",
					ID:            1,
					HashCode:      0,
					Title:         "title",
					Writer:        model.Writer{ID: 1, Name: "writer"},
					Type:          "type",
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Status:        model.StatusInProgress,
					IsDownloaded:  true,
					Error:         model.Error{Err: errors.New("error")},
				},
			},
			expectErr: nil,
		},
		{
			name:         "error/invalid json",
			data:         []byte(`abc`),
			expectParams: nil,
			expectErr:    &json.SyntaxError{},
		},
		{
			name: "error/missing website",
			data: []byte(`{"trace_id":"01234567890123456789012345678901","span_id":"0123456789012345","trace_flags":1}`),
			expectParams: &BookProcessParams{
				TraceID:    "01234567890123456789012345678901",
				SpanID:     "0123456789012345",
				TraceFlags: 0x1,
			},
			expectErr: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			ctx, params, err := ParamsFromData(context.Background(), test.data)

			if test.expectErr != nil {
				assert.ErrorAs(t, err, &test.expectErr, "different error")
			}
			assert.Equal(t, test.expectParams, params, "different params")
			if params != nil && params.TraceID != "" && params.SpanID != "" {
				assert.Equal(t, trace.SpanContextFromContext(ctx).TraceID().String(), params.TraceID, "different trace id")
				assert.Equal(t, trace.SpanContextFromContext(ctx).SpanID().String(), params.SpanID, "different span id")
				assert.Equal(t, trace.SpanContextFromContext(ctx).TraceFlags(), trace.TraceFlags(params.TraceFlags), "different trace flags")
			}
		})
	}
}

func TestWebsiteUpdateParams_MarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		params      *BookProcessParams
		expect      string
		expectError error
	}{
		{
			name: "success",
			params: &BookProcessParams{
				Book: model.Book{
					Site:          "test",
					ID:            1,
					HashCode:      0,
					Title:         "title",
					Writer:        model.Writer{ID: 1, Name: "writer"},
					Type:          "type",
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Status:        model.StatusInProgress,
					IsDownloaded:  true,
					Error:         model.Error{Err: errors.New("error")},
				},
				TraceID:    "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				SpanID:     "XXXXXXXXXXXXXXXX",
				TraceFlags: 0x1,
			},
			expect:      `{"book":{"site":"test","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX","span_id":"XXXXXXXXXXXXXXXX","trace_flags":1}`,
			expectError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(test.params)
			assert.Equal(t, test.expect, string(data))
			assert.ErrorIs(t, err, test.expectError)
		})
	}
}

func TestWebsiteUpdateParams_ToData(t *testing.T) {
	emptyCtx := context.Background()
	spanCtx, span := otel.Tracer("test").Start(emptyCtx, "test")
	span.End()

	tests := []struct {
		name        string
		ctx         context.Context
		params      *BookProcessParams
		expect      string
		expectError error
	}{
		{
			name: "success/with span",
			ctx:  spanCtx,
			params: &BookProcessParams{
				Book: model.Book{
					Site:          "test",
					ID:            1,
					HashCode:      0,
					Title:         "title",
					Writer:        model.Writer{ID: 1, Name: "writer"},
					Type:          "type",
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Status:        model.StatusInProgress,
					IsDownloaded:  true,
					Error:         model.Error{Err: errors.New("error")},
				},
			},
			expect: fmt.Sprintf(
				`{"book":{"site":"test","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"%s","span_id":"%s","trace_flags":1}`,
				span.SpanContext().TraceID().String(),
				span.SpanContext().SpanID().String(),
			),
			expectError: nil,
		},
		{
			name: "success/without span",
			ctx:  emptyCtx,
			params: &BookProcessParams{
				Book: model.Book{
					Site:          "test",
					ID:            1,
					HashCode:      0,
					Title:         "title",
					Writer:        model.Writer{ID: 1, Name: "writer"},
					Type:          "type",
					UpdateDate:    "date",
					UpdateChapter: "chapter",
					Status:        model.StatusInProgress,
					IsDownloaded:  true,
					Error:         model.Error{Err: errors.New("error")},
				},
			},
			expect:      `{"book":{"site":"test","id":1,"hash_code":0,"title":"title","type":"type","update_date":"date","update_chapter":"chapter","status":"INPROGRESS","is_downloaded":true,"writer":{"id":1,"name":"writer"},"error":"error"},"trace_id":"00000000000000000000000000000000","span_id":"0000000000000000","trace_flags":0}`,
			expectError: nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			data, err := test.params.ToData(test.ctx)
			assert.Equal(t, test.expect, string(data))
			assert.ErrorIs(t, err, test.expectError)
		})
	}

}
