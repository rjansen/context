package http

import (
	"context"
	"github.com/rjansen/l"
	"net/http"
)

var (
	ContextKeys = Keys{
		LOG:      "requestLog",
		TOKEN:    "requestToken",
		IDENTITY: "requestIdentity",
		AUDITOR:  "requestAuditor",
	}
)

type Keys struct {
	LOG      string
	TOKEN    string
	IDENTITY string
	AUDITOR  string
}

type Auditor struct {
	Log      l.Logger
	Identity *Identity
}

type Identity struct {
	Token string      `json:"token"`
	Value interface{} `jsno:"value"`
}

func set(r *http.Request, key interface{}, val interface{}) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), key, val))
}

func Get(r *http.Request, key interface{}) interface{} {
	return r.Context().Value(key)
}

func GetLog(r *http.Request) l.Logger {
	return Get(r, ContextKeys.LOG).(l.Logger)
}

func GetToken(r *http.Request) string {
	return Get(r, ContextKeys.TOKEN).(string)
}

func GetIdentity(r *http.Request) *Identity {
	return Get(r, ContextKeys.IDENTITY).(*Identity)
}

func GetAuditor(r *http.Request) *Auditor {
	return Get(r, ContextKeys.AUDITOR).(*Auditor)
}
