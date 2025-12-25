package bookprocess

import (
	"context"
	"encoding/json"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"
)

type BookProcessParams struct {
	Book       model.Book `json:"book"`
	TraceID    string     `json:"trace_id"`
	SpanID     string     `json:"span_id"`
	TraceFlags byte       `json:"trace_flags"`
}

func ParamsFromData(ctx context.Context, data []byte) (context.Context, *BookProcessParams, error) {
	// parse message body
	params := new(BookProcessParams)
	if jsonErr := json.Unmarshal(data, params); jsonErr != nil {
		return ctx, nil, jsonErr
	}

	ctx = log.With().
		Str("trace_id", params.TraceID).
		Str("site", params.Book.Site).
		Int("book_id", params.Book.ID).
		Str("book_hash", params.Book.FormatHashCode()).
		Str("book_title", params.Book.Title).
		Str("book_writer", params.Book.Writer.Name).
		Str("book_status", params.Book.Status.String()).
		Logger().WithContext(ctx)

	if params.TraceID != "" && params.SpanID != "" {
		traceID, traceErr := trace.TraceIDFromHex(params.TraceID)
		spanID, spanErr := trace.SpanIDFromHex(params.SpanID)
		if traceErr != nil || spanErr != nil {
			return ctx, params, nil
		}

		spanContext := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.TraceFlags(params.TraceFlags),
			Remote:     true, // Indicate that this span context is from a remote service
		})
		ctx := trace.ContextWithSpanContext(ctx, spanContext)

		return ctx, params, nil
	}

	return ctx, params, nil
}

func (params *BookProcessParams) ToData(ctx context.Context) ([]byte, error) {
	spanCtx := trace.SpanContextFromContext(ctx)
	params.TraceID = spanCtx.TraceID().String()
	params.SpanID = spanCtx.SpanID().String()
	params.TraceFlags = byte(spanCtx.TraceFlags())

	return json.Marshal(params)
}
