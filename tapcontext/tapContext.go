package tapcontext

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"go.elastic.co/apm"
)

type customContextType string

const (
	//TAPCtx - defining a separate type to avoid colliding with basic type
	TAPCtx customContextType = "tapCtx"
)

// TapContext contains context of client
type TapContext struct {
	UserEmail      string              // Email of the user
	RequestID      string              // RequestID - used to track logs across a request-response cycle
	PermissionsMap map[string][]string // this map will help in flagging
	TapApiToken    string              // TapApiToken - used to authenticate the session/request
	Application    string              // application for dynamic application auth
	Locale         string              // Locale for language
}

// TContext is the combination of native context and TapContext
type TContext struct {
	context.Context
	TapContext
}

// GetTapCtx returns the tap context from the context provided
func GetTapCtx(ctx context.Context) (TapContext, bool) {
	if ctx == nil {
		return TapContext{}, false
	}
	tapCtx, exists := ctx.Value(TAPCtx).(TapContext)
	return tapCtx, exists
}

// WithTapCtx returns a new context with the tap context provided
func WithTapCtx(ctx context.Context, tapctx TapContext) context.Context {
	return context.WithValue(ctx, TAPCtx, tapctx)
}

// UpgradeCtx embeds native context and TapContext
func UpgradeCtx(ctx context.Context) TContext {
	var tContext TContext
	tapCtx, _ := GetTapCtx(ctx)

	tContext.Context = ctx
	tContext.TapContext = tapCtx
	return tContext
}

func NewTapContext() TContext {
	return TContext{
		Context:    context.Background(),
		TapContext: TapContext{},
	}
}

func CopyTapContext(ctx context.Context) TContext {
	tapCtx, _ := GetTapCtx(ctx)
	return TContext{
		Context:    context.Background(),
		TapContext: tapCtx,
	}
}

func New(id ...string) TContext {
	var requestID string
	if len(id) > 0 {
		requestID = id[0]
	}
	if len(requestID) == 0 {
		requestID = strings.ReplaceAll(uuid.NewString(), "-", "")
	}
	tapCtx := TapContext{
		RequestID: requestID,
	}
	ctx := UpgradeCtx(WithTapCtx(context.Background(), tapCtx))
	return ctx
}

func NewTapContextAndAPMTransactionOps() (TContext, apm.TransactionOptions) {
	var traceID apm.TraceID
	uuid := uuid.New()
	var spanID apm.SpanID
	var traceOptions apm.TraceOptions
	copy(traceID[:], uuid[:])
	copy(spanID[:], traceID[8:])
	tapCtx := TapContext{
		RequestID: traceID.String(),
	}
	return UpgradeCtx(WithTapCtx(context.Background(), tapCtx)), apm.TransactionOptions{
		TraceContext: apm.TraceContext{
			Trace:   traceID,
			Span:    spanID,
			Options: traceOptions.WithRecorded(true),
		},
		TransactionID: spanID,
	}
}
