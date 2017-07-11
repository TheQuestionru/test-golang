package types

import (
	"github.com/TheQuestionru/test-golang/server/lib/logger"
	"strings"
)

type Req struct {
	AccountId int64
	WidgetId  int64
	Ip        string
	Ua        string
	Device    string
	Referer   string
	Hostname  string
	Params    ReqParams
	Lang      LanguageType
}

type ReqParams struct {
	Qml bool // Enables QML in responses.
}

func NewEmptyReq() Req {
	return Req{Lang: LanguageRU}
}

func NewReq(accountId int64) Req {
	return Req{
		AccountId: accountId,
		Lang:      LanguageRU,
	}
}

func (c Req) IsAuthorized() bool {
	return c.AccountId != 0
}

func (c Req) Authorize() error {
	if !c.IsAuthorized() {
		return ErrUnauthorized
	}
	return nil
}

func (c Req) WithAccount(account int64) Req {
	var r Req
	r = c
	r.AccountId = account
	return r
}

func (c *Req) Log() logger.Fields {
	if c.AccountId == 0 && c.WidgetId == 0 && len(c.Ip) == 0 {
		return nil
	}

	ret := logger.Fields{}

	if c.AccountId != 0 {
		ret["accountId"] = c.AccountId
	}

	if c.WidgetId != 0 {
		ret["widgetId"] = c.WidgetId
	}

	if len(c.Ip) != 0 {
		ret["ip"] = c.Ip
	}

	if len(c.Lang) != 0 {
		ret["lang"] = c.Lang
	}

	return ret
}

func (c *Req) IsAndroidApp() bool {
	if strings.Contains(c.Ua, "thequestion_android") {
		return true
	}
	return false
}
