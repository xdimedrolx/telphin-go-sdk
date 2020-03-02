package telphin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// Host points to the API
	Host = "https://apiproxy.telphin.ru"
	// Host points to the storage
	HostStorage = "https://storage.telphin.ru"
	// Websocket host
	WsHost = "sipproxy.telphin.ru"

	// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
	RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second
)

const (
	HangupDispositionCalleeBye      = "callee_bye"      // трубку положила принимающая сторона
	HangupDispositionCallerBye      = "caller_bye"      // трубку положила вызывающая сторона
	HangupDispositionCallerCancel   = "caller_cancel"   // вызывающая сторона отказалась ждать ответа
	HangupDispositionCalleeRefuse   = "callee_refuse"   // принимающая сторона отказалась отвечать (была занята, отсутствовала регистрация и т.п.)
	HangupDispositionInternalCancel = "internal_cancel" // вызов завершен сервером (обычно таймаут вызова или при пачке вызовов кто-то снял трубку)
)

const (
	CallResultBusy        = "busy"
	CallResultAnswered    = "answered"
	CallResultNotAnswered = "not answered"
	CallResultFailed      = "failed"
	CallResultRejected    = "rejected"
	CallResultBridged     = "bridged"

	CallCdrResultBusy              = "busy"
	CallCdrResultAnswered          = "answered"
	CallCdrResultNotAnswered       = "not answered"
	CallCdrResultAnsweredElsewhere = "answered elsewhere"
	CallCdrResultFailed            = "failed"
	CallCdrResultRejected          = "rejected"
	CallCdrResultVoicemail         = "voicemail"
)

const (
	EventStatusCalling     = "CALLING"     // поступил вызов
	EventStatusAnswer      = "ANSWER"      // вызов был отвечен
	EventStatusBusy        = "BUSY"        // вызов получил сигнал "занято"
	EventStatusNoAnswer    = "NOANSWER"    // звонок не отвечен (истек таймер ожидания на сервере)
	EventStatusCancel      = "CANCEL"      // звонящий отменил вызов до истечения таймера ожидания на сервере
	EventStatusCongestion  = "CONGESTION"  // произошла ошибка во время вызова
	EventStatusChanunavail = "CHANUNAVAIL" // у вызываемого абонента отсутствует регистрация
)

const (
	EventTypeDialIn  = "dial-in"
	EventTypeDialOut = "dial-out"
	EventTypeAnswer  = "answer"
	EventTypeHangup  = "hangup"
)

const (
	FlowIn       = "in"
	FlowOut      = "out"
	FlowTransfer = "transfer"
)

const (
	ExtensionTypePhone = "phone"
	ExtensionTypeIVR   = "ivr"
	ExtensionTypeQueue = "queue"
)

var extensionNumberRegexp = regexp.MustCompile(`\d+\*(\d+)`)
var TelphinStorageHostRegexp = regexp.MustCompile(`^` + HostStorage)

type (
	JSONTime time.Time

	Client struct {
		sync.Mutex
		Client         *http.Client
		ClientID       string
		Secret         string
		Host           string
		Token          *TokenResponse
		tokenExpiresAt time.Time
		Logger         FieldLogger
	}

	expirationTime int64

	// TokenResponse is for API response for the /oauth2/token endpoint
	TokenResponse struct {
		// RefreshToken string         `json:"refresh_token"`
		Token     string         `json:"access_token"`
		Type      string         `json:"token_type"`
		ExpiresIn expirationTime `json:"expires_in"`
	}

	// ErrorResponse a typed error returned by http Handlers and used for choosing error handlers
	ErrorResponse struct {
		Code    int
		Cause   string `json:"message"`
		Details string
	}

	// See: https://ringme-confluence.atlassian.net/wiki/spaces/RAL/pages/17367181/extension
	Extension struct {
		ID     uint32 `json:"id"`     // Уникальный идентификатор добавочного
		Status string `json:"status"` // Статус добавочного: 'active' - активен, 'blocked' - заблокирован
		Name   string `json:"name"`   // префикс_клиента*имя_добавочного или просто имя_добавочного(в этом случчае префикс будет дописан автоматически)
		Type   string `json:"type"`   // Тип добавочного: 'phone', 'queue', 'ivr', 'fax'
		Label  string `json:"label"`  // Display Name добавочного. Отображается на вызываемом терминале при исходящих вызовах (если поддерживается)
	}

	ExtensionCreateRequest struct {
		Status string `json:"status"`          // Статус добавочного: 'active' - активен, 'blocked' - заблокирован
		Name   string `json:"name"`            // префикс_клиента*имя_добавочного или просто имя_добавочного(в этом случчае префикс будет дописан автоматически)
		Type   string `json:"type"`            // Тип добавочного: 'phone', 'queue', 'ivr', 'fax'
		Label  string `json:"label,omitempty"` // Display Name добавочного. Отображается на вызываемом терминале при исходящих вызовах (если поддерживается)
	}

	// See: https://ringme-confluence.atlassian.net/wiki/spaces/RAL/pages/832602113/ani
	Ani struct {
		Default string `json:"default"` // Установленный ani, если не установлен - None
	}

	// See: https://ringme-confluence.atlassian.net/wiki/spaces/RAL/pages/20414492/did
	Did struct {
		ID          uint32 `json:"id"`           // Идентификатор внешнего номера
		Name        string `json:"name"`         // Имя внешнего номера, то есть сам номер
		Domain      string `json:"domain"`       // Домен внешнего номера
		ClientID    uint32 `json:"client_id"`    // Идентификатор клиента, которому назначен внешний номер
		ExtensionID uint32 `json:"extension_id"` // Идентификатор добавочного, которому назначен внешний номер
	}

	PhoneProperties struct {
		Password       string `json:"password"`        // Задает пароль для авторизации на sip-сервере
		AllowWebrtc    bool   `json:"allow_webrtc"`    // Если установлено значение "true", то позволяет добавочному работать через WebRTC
		RecordEnabled  bool   `json:"record_enabled"`  // Включить запись разговоров
		RecordTransfer bool   `json:"record_transfer"` // Включить запись в том числе и переадресованных вызовов
	}

	CallHistory struct {
		UUID              string    `json:"call_uuid"`
		Flow              string    `json:"flow"`
		InitTime          *JSONTime `json:"init_time_gmt"`
		StartTime         *JSONTime `json:"start_time_gmt"`
		BridgedTime       *JSONTime `json:"bridged_time_gmt"`
		HangupTime        *JSONTime `json:"hangup_time_gmt"`
		Duration          uint32    `json:"duration"`
		Bridged           bool      `json:"bridged"`
		BridgedDuration   uint32    `json:"bridged_duration"`
		ExtensionID       uint32    `json:"extension_id"`
		From              string    `json:"from_username"`
		To                string    `json:"to_username"`
		Result            string    `json:"result"`
		HangupCause       string    `json:"hangup_cause"`
		HangupDisposition *string   `json:"hangup_disposition,omitempty"`
		Cdr               *[]Cdr    `json:"cdr"`
	}

	Cdr struct {
		ExtensionType     string    `json:"extension_type"`
		RecordUUID        *string   `json:"record_uuid"`
		Result            string    `json:"result"`
		HangupCause       string    `json:"hangup_cause"`
		HangupDisposition *string   `json:"hangup_disposition,omitempty"`
		Duration          int       `json:"duration"`
		InitTime          *JSONTime `json:"init_time_gmt"`
		StartTime         *JSONTime `json:"start_time_gmt"`
		HangupTime        *JSONTime `json:"hangup_time_gmt"`
		Flow              string    `json:"flow"`
		ExtensionID       uint32    `json:"extension_id"`
		From              string    `json:"from_username"`
		To                string    `json:"to_username"`
	}

	CallHistories struct {
		Page        uint32        `json:"page"`
		PerPage     uint32        `json:"per_page"`
		Order       string        `json:"order"`
		CallHistory []CallHistory `json:"call_history"`
	}

	CallHistoryRequest struct {
		StartDatetime *string `url:"start_datetime"`
		EndDatetime   *string `url:"end_datetime"`
		Flow          *string `url:"flow"`
		ExtensionID   *uint32 `url:"extension_id"`
		ToUsername    *string `url:"to_username"`
		PerPage       uint16  `url:"per_page"`
	}

	RecordStorageUrl struct {
		RecordUrl string `json:"record_url"`
	}

	CreateEventRequest struct {
		URL       string `json:"url"`
		Method    string `json:"method"`
		EventType string `json:"event_type"`
	}

	Event struct {
		ID        int    `json:"id"`
		URL       string `json:"url"`
		Method    string `json:"method"`
		EventType string `json:"event_type"`
	}
)

func (e Extension) Number() (*int, error) {
	v := extensionNumberRegexp.FindStringSubmatch(e.Name)
	if len(v) != 2 {
		return nil, fmt.Errorf("Unknown number %s", e.Name)
	}
	number, err := strconv.ParseInt(v[1], 10, 64)
	if err != nil {
		return nil, err
	}
	n := int(number)
	return &n, nil
}

func (c *CallHistory) HasTransferredCdr() bool {
	for _, cdr := range *c.Cdr {
		if cdr.Flow == FlowTransfer {
			return true
		}
	}
	return false
}

func (c *CallHistory) HasExtensionPhoneTypeCdr() bool {
	for _, cdr := range *c.Cdr {
		if cdr.ExtensionType == ExtensionTypePhone {
			return true
		}
	}
	return false
}

func (c *CallHistory) HasAnsweredExtensionPhoneTypeCdr() bool {
	for _, cdr := range *c.Cdr {
		if cdr.ExtensionType == ExtensionTypePhone && cdr.Result == CallCdrResultAnswered {
			return true
		}
	}
	return false
}

func (c *CallHistory) HasExtensionPhoneTypeWithRecords() bool {
	for _, cdr := range *c.Cdr {
		if cdr.ExtensionType == ExtensionTypePhone && cdr.RecordUUID != nil {
			return true
		}
	}
	return false
}

// imeplement Marshaler und Unmarshalere interface
func (j *JSONTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}
	*j = JSONTime(t)
	return nil
}

func (j JSONTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (j JSONTime) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Code: %d. Error: %s. Details: %s", e.Code, e.Cause, e.Details)
}
