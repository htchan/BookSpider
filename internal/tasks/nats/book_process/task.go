package bookprocess

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/htchan/BookSpider/internal/model"
	"github.com/htchan/BookSpider/internal/service/v1"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type BookProcessTask struct {
	Site    string
	nc      *nats.Conn
	Service service.BookService
}

func getTracer() trace.Tracer {
	return otel.Tracer("htchan/BookSpider/book-process")
}

func NewTask(
	site string,
	nc *nats.Conn,
	serv service.BookService,
) *BookProcessTask {
	return &BookProcessTask{
		Site:    site,
		nc:      nc,
		Service: serv,
	}
}

func (task *BookProcessTask) subject() string {
	return fmt.Sprintf("book_spider.books.process.%s", task.Site)
}

func (task *BookProcessTask) Publish(
	ctx context.Context,
	bk *model.Book,
) error {
	params := BookProcessParams{
		Book: *bk,
	}

	data, err := params.ToData(ctx)
	if err != nil {
		return err
	}

	err = task.nc.Publish(task.subject(), data)
	if err != nil {
		return err
	}

	return nil
}

func (task *BookProcessTask) Subscribe(ctx context.Context) (jetstream.ConsumeContext, error) {
	js, err := jetstream.New(task.nc)
	if err != nil {
		return nil, fmt.Errorf("init jetstream fail: %v", err)
	}

	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     strings.ReplaceAll(task.subject(), ".", "-"),
		Subjects: []string{task.subject()},
		MaxAge:   time.Hour * 24 * 7,
	})
	if err != nil {
		return nil, fmt.Errorf("create / update stream fail: %v", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:      strings.ReplaceAll(task.subject(), ".", "-"),
		Durable:   strings.ReplaceAll(task.subject(), ".", "-"),
		AckPolicy: jetstream.AckExplicitPolicy,
		AckWait:   time.Minute * 60,
	})
	if err != nil {
		return nil, fmt.Errorf("create / update consumer fail: %v", err)
	}

	return consumer.Consume(task.handler)
}

func (task *BookProcessTask) Validate(ctx context.Context, params *BookProcessParams) error {
	// validate params
	_, validateSpan := getTracer().Start(ctx, "Validate Params")
	defer validateSpan.End()

	validateSpan.SetAttributes(
		attribute.String("vendor_name", task.Site),
		attribute.String("book_site", params.Book.Site),
		attribute.Int("book_id", params.Book.ID),
		attribute.String("book_hash_code", params.Book.FormatHashCode()),
		attribute.String("book_name", params.Book.Title),
		attribute.String("book_writer", params.Book.Writer.Name),
		attribute.String("book_status", params.Book.Status.String()),
	)

	if !task.Service.SupportBook(&params.Book) || params.Book.Site != task.Site {
		validateSpan.SetStatus(codes.Error, ErrNotSupportedBook.Error())
		validateSpan.RecordError(ErrNotSupportedBook)

		return ErrNotSupportedBook
	}

	return nil
}

func (task *BookProcessTask) handler(msg jetstream.Msg) {
	ctx := log.With().
		Str("task", "book-process").
		Str("book_site", task.Site).
		Logger().WithContext(context.Background())

	defer func() {
		ackErr := msg.Ack()
		if ackErr != nil {
			zerolog.Ctx(ctx).Error().Err(ackErr).Msg("ack failed")
		}
	}()

	// parse message body
	ctx, params, err := ParamsFromData(ctx, msg.Data())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).
			Str("data", string(msg.Data())).
			Msg("failed to parse message body")

		return
	}

	ctx, span := getTracer().Start(ctx, "Process Book", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	span.SetAttributes(
		append(
			params.Book.OtelAttributes(),
		)...,
	)

	// validate params
	validateErr := task.Validate(ctx, params)
	if validateErr != nil {
		zerolog.Ctx(ctx).Error().Err(validateErr).Msg("validate params failed")

		return
	}

	// sleep after update
	defer func() {
		time.Sleep(time.Second)
	}()

	// call vendor service to update website
	updateCtx, updateSpan := getTracer().Start(ctx, "Vendor Service Call")
	defer updateSpan.End()

	processErr := task.Service.ProcessBook(updateCtx, &params.Book)
	if processErr != nil {
		zerolog.Ctx(ctx).Error().Err(processErr).Msg("process book failed")
		updateSpan.SetStatus(codes.Error, processErr.Error())
		updateSpan.RecordError(processErr)
		return
	}

	updateSpan.End()
}
