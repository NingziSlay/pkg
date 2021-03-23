package middleware

import (
	"context"
	log2 "github.com/NingziSlay/pkg/log"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

var (
	// TracingComponentTag tags
	TracingComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "gRPC"}
)

// MDReaderWriter metadata Reader and Writer
type MDReaderWriter struct {
	metadata.MD
}

// ForeachKey range all keys to call handler
func (c MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter
func (c MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}

// OpenTracingClientInterceptor  rewrite client's interceptor with open tracing
func OpenTracingClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, resp interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var parentCtx opentracing.SpanContext
		if parent := opentracing.SpanFromContext(ctx); parent != nil {
			parentCtx = parent.Context()
		}
		cliSpan := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			TracingComponentTag,
			ext.SpanKindRPCClient,
		)
		defer cliSpan.Finish()
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		err := tracer.Inject(cliSpan.Context(), opentracing.TextMap, MDReaderWriter{md})
		if err != nil {
			l := log2.GetLogger()
			l.Warn().Err(err).Msg("can not inject to metadata")
		}
		ctx = metadata.NewOutgoingContext(ctx, md)
		err = invoker(ctx, method, req, resp, cc, opts...)
		if err != nil {
			cliSpan.LogFields(log.String("err", err.Error()))
		}
		return err
	}
}

// LogInterceptor rpc 接口日志中间件，用来记录每个 RPC 调用的
// 请求参数和返回值
func LogInterceptor(log zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		start := time.Now()
		defer func(s time.Time) {
			log.Info().
				Str("gRPC.request.method", info.FullMethod).
				AnErr("gRPC.response.err", err).
				Dur("gRPC.response.took_ms", time.Since(s)).
				Send()
		}(start)
		resp, err = handler(ctx, req)
		return
	}
}
