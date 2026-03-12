package splunkhec

import (
	"context"
	"io"
	"net/http"
)

type HEC interface {
	SetHTTPClient(client *http.Client)
	SetKeepAlive(enable bool)
	SetChannel(channel string)
	SetMaxRetry(retries int)
	SetMaxContentLength(size int)
	SetCompression(compression string)

	WriteEvent(event *Event) error
	WriteBatch(events []*Event) error
	WriteBatchWithContext(ctx context.Context, events []*Event) error
	WriteRaw(reader io.ReadSeeker, metadata *EventMetadata) error
	WriteRawWithContext(ctx context.Context, reader io.ReadSeeker, metadata *EventMetadata) error
	WaitForAcknowledgement() error
	WaitForAcknowledgementWithContext(ctx context.Context) error
}
