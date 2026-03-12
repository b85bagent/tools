package splunkhec

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	retryWaitTime = 1 * time.Second

	defaultMaxContentLength       = 1_000_000
	defaultAcknowledgementTimeout = 90 * time.Second
)

type Client struct {
	HEC

	httpClient *http.Client
	serverURL  string
	token      string

	keepAlive bool
	channel   string
	retries   int
	maxLength int

	ackIDs []int
	ackMux sync.Mutex

	compression string
}

func NewClient(serverURL string, token string) HEC {
	return &Client{
		httpClient: http.DefaultClient,
		serverURL:  serverURL,
		token:      token,
		keepAlive:  true,
		channel:    newChannelID(),
		retries:    2,
		maxLength:  defaultMaxContentLength,
	}
}

func (hec *Client) SetHTTPClient(client *http.Client) {
	hec.httpClient = client
}

func (hec *Client) SetKeepAlive(enable bool) {
	hec.keepAlive = enable
}

func (hec *Client) SetChannel(channel string) {
	hec.channel = channel
}

func (hec *Client) SetMaxRetry(retries int) {
	hec.retries = retries
}

func (hec *Client) SetMaxContentLength(size int) {
	hec.maxLength = size
}

func (hec *Client) SetCompression(compression string) {
	hec.compression = compression
}

func (hec *Client) WriteEventWithContext(ctx context.Context, event *Event) error {
	if event.empty() {
		return nil
	}

	endpoint := "/services/collector?channel=" + hec.channel
	data, _ := json.Marshal(event)

	if len(data) > hec.maxLength {
		return ErrEventTooLong
	}
	return hec.write(ctx, endpoint, data)
}

func (hec *Client) WriteEvent(event *Event) error {
	return hec.WriteEventWithContext(context.Background(), event)
}

func (hec *Client) WriteBatchWithContext(ctx context.Context, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	endpoint := "/services/collector?channel=" + hec.channel
	var buffer bytes.Buffer
	var tooLongs []int

	for index, event := range events {
		if event.empty() {
			continue
		}

		data, _ := json.Marshal(event)
		if len(data) > hec.maxLength {
			tooLongs = append(tooLongs, index)
			continue
		}
		if buffer.Len()+len(data) > hec.maxLength {
			if err := hec.write(ctx, endpoint, buffer.Bytes()); err != nil {
				return err
			}
			buffer.Reset()
		}
		buffer.Write(data)
	}

	if buffer.Len() > 0 {
		if err := hec.write(ctx, endpoint, buffer.Bytes()); err != nil {
			return err
		}
	}
	if len(tooLongs) > 0 {
		return ErrEventTooLong
	}
	return nil
}

func (hec *Client) WriteBatch(events []*Event) error {
	return hec.WriteBatchWithContext(context.Background(), events)
}

type EventMetadata struct {
	Host       *string
	Index      *string
	Source     *string
	SourceType *string
	Time       *time.Time
}

func (hec *Client) WriteRawWithContext(ctx context.Context, reader io.ReadSeeker, metadata *EventMetadata) error {
	endpoint := rawHecEndpoint(hec.channel, metadata)

	return breakStream(reader, hec.maxLength, func(chunk []byte) error {
		if err := hec.write(ctx, endpoint, chunk); err != nil {
			if res, ok := err.(*Response); !ok || res.Code != StatusNoData {
				return err
			}
		}
		return nil
	})
}

func (hec *Client) WriteRaw(reader io.ReadSeeker, metadata *EventMetadata) error {
	return hec.WriteRawWithContext(context.Background(), reader, metadata)
}

type acknowledgementRequest struct {
	Acks []int `json:"acks"`
}

func (hec *Client) WaitForAcknowledgementWithContext(ctx context.Context) error {
	hec.ackMux.Lock()
	ackIDs := hec.ackIDs
	hec.ackIDs = nil
	hec.ackMux.Unlock()

	if len(ackIDs) == 0 {
		return nil
	}

	endpoint := "/services/collector/ack?channel=" + hec.channel

	for {
		ackRequestData, _ := json.Marshal(acknowledgementRequest{Acks: ackIDs})

		response, err := hec.makeRequest(ctx, endpoint, ackRequestData)
		if err != nil {
			hec.ackMux.Lock()
			hec.ackIDs = append(hec.ackIDs, ackIDs...)
			hec.ackMux.Unlock()
			return err
		}

		for ackIDString, status := range response.Acks {
			if status {
				ackID, err := strconv.Atoi(ackIDString)
				if err != nil {
					return fmt.Errorf("could not convert ack ID to int: %v", err)
				}

				ackIDs = removeInt(ackIDs, ackID)
			}
		}

		if len(ackIDs) == 0 {
			break
		}

		select {
		case <-time.After(retryWaitTime):
			continue
		case <-ctx.Done():
			hec.ackMux.Lock()
			hec.ackIDs = append(hec.ackIDs, ackIDs...)
			hec.ackMux.Unlock()
			return ctx.Err()
		}
	}

	return nil
}

func (hec *Client) WaitForAcknowledgement() error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAcknowledgementTimeout)
	defer cancel()
	return hec.WaitForAcknowledgementWithContext(ctx)
}

func breakStream(reader io.ReadSeeker, max int, callback func(chunk []byte) error) error {
	buf := make([]byte, max+1)
	var writeAt int
	for {
		n, err := reader.Read(buf[writeAt:max])
		if n == 0 && err == io.EOF {
			break
		}

		if err == io.EOF && buf[writeAt+n-1] != '\n' {
			n++
			buf[writeAt+n-1] = '\n'
		}

		data := buf[0 : writeAt+n]

		cut := bytes.LastIndexByte(data, '\n') + 1
		if cut == 0 {
			cut = len(data)
		}
		if err := callback(buf[:cut]); err != nil {
			return err
		}

		writeAt = copy(buf, data[cut:])

		if err != nil && err != io.EOF {
			return err
		}
	}

	if writeAt != 0 {
		return callback(buf[:writeAt])
	}

	return nil
}

func responseFrom(body []byte) *Response {
	var res Response
	_ = json.Unmarshal(body, &res)
	return &res
}

func (res *Response) Error() string {
	return res.Text
}

func (res *Response) String() string {
	b, _ := json.Marshal(res)
	return string(b)
}

func (hec *Client) makeRequest(ctx context.Context, endpoint string, data []byte) (*Response, error) {
	retries := 0
RETRY:
	var reader io.Reader
	if hec.compression == "gzip" {
		var buffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&buffer)
		if _, err := gzipWriter.Write(data); err != nil {
			return nil, err
		}
		gzipWriter.Close()
		reader = &buffer
	} else {
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(http.MethodPost, hec.serverURL+endpoint, reader)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if hec.keepAlive {
		req.Header.Set("Connection", "keep-alive")
	}
	req.Header.Set("Authorization", "Splunk "+hec.token)
	if hec.compression == "gzip" {
		req.Header.Set("Content-Encoding", "gzip")
	}
	res, err := hec.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := responseFrom(body)

	if res.StatusCode != http.StatusOK {
		if retriable(response.Code) && retries < hec.retries {
			retries++
			time.Sleep(retryWaitTime)
			goto RETRY
		}
	}

	return response, nil
}

func (hec *Client) write(ctx context.Context, endpoint string, data []byte) error {
	response, err := hec.makeRequest(ctx, endpoint, data)
	if err != nil {
		return err
	}

	// 有些 HEC/代理回傳 code=0 但 text 為空，視同成功
	if response.Text != "Success" && response.Code != StatusSuccess {
		return response
	}

	if response.AckID != nil {
		hec.ackMux.Lock()
		hec.ackIDs = append(hec.ackIDs, *response.AckID)
		hec.ackMux.Unlock()
	}

	return nil
}

func rawHecEndpoint(channel string, metadata *EventMetadata) string {
	var buffer bytes.Buffer
	buffer.WriteString("/services/collector/raw?channel=" + channel)
	if metadata == nil {
		return buffer.String()
	}
	if metadata.Host != nil {
		buffer.WriteString("&host=" + *metadata.Host)
	}
	if metadata.Index != nil {
		buffer.WriteString("&index=" + *metadata.Index)
	}
	if metadata.Source != nil {
		buffer.WriteString("&source=" + *metadata.Source)
	}
	if metadata.SourceType != nil {
		buffer.WriteString("&sourcetype=" + *metadata.SourceType)
	}
	if metadata.Time != nil {
		buffer.WriteString("&time=" + epochTime(metadata.Time))
	}
	return buffer.String()
}

func removeInt(list []int, target int) []int {
	for i, v := range list {
		if v == target {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func newChannelID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
