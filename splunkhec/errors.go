package splunkhec

import "errors"

// Response represents Splunk HEC 回傳的 JSON 結構，例如 {"text":"Success","code":0}
type Response struct {
	Text  string          `json:"text"`
	Code  int             `json:"code"`
	AckID *int            `json:"ackId"`
	Acks  map[string]bool `json:"acks"`
}

// Response status codes
const (
	StatusSuccess              = 0
	StatusTokenDisabled        = 1
	StatusTokenRequired        = 2
	StatusInvalidAuthorization = 3
	StatusInvalidToken         = 4
	StatusNoData               = 5
	StatusInvalidDataFormat    = 6
	StatusIncorrectIndex       = 7
	StatusInternalServerError  = 8
	StatusServerBusy           = 9
	StatusChannelMissing       = 10
	StatusInvalidChannel       = 11
	StatusEventFieldRequired   = 12
	StatusEventFieldBlank      = 13
	StatusAckDisabled          = 14
)

func retriable(code int) bool {
	return code == StatusServerBusy || code == StatusInternalServerError
}

var ErrEventTooLong = errors.New("event length is too long")
