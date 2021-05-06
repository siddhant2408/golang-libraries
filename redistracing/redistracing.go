// Package redistracing provides tracing for Redis.
package redistracing

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	opentracing_ext "github.com/opentracing/opentracing-go/ext"
	"github.com/siddhant2408/golang-libraries/closeutils"
	"github.com/siddhant2408/golang-libraries/errors"
	"github.com/siddhant2408/golang-libraries/tracingutils"
)

const (
	externalServiceName = "go-redis"
)

// Hook implements redis.Hook and traces commands.
type Hook struct {
	Addrs []string
	DB    *int
}

// BeforeProcess implements redis.Hook.
func (h *Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	span := h.startSpan(&ctx, "redis.command")
	tracingutils.SetSpanResourceName(span, cmd.FullName())
	h.setSpanTagKey(span, cmd)
	return ctx, nil
}

// AfterProcess implements redis.Hook.
func (h *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span, spanFinish := getSpanFromContext(ctx)
	defer spanFinish()
	err := cmd.Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		tracingutils.SetSpanError(span, err)
	}
	return nil
}

// BeforeProcessPipeline implements redis.Hook.
func (h *Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	span := h.startSpan(&ctx, "redis.pipeline")
	span.SetTag("redis.commands.count", len(cmds))
	return ctx, nil
}

// AfterProcessPipeline implements redis.Hook.
func (h *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span, spanFinish := getSpanFromContext(ctx)
	defer spanFinish()
	for _, cmd := range cmds {
		err := cmd.Err()
		if err != nil && !errors.Is(err, redis.Nil) {
			tracingutils.SetSpanError(span, err)
			break
		}
	}
	return nil
}

func (h *Hook) startSpan(pctx *context.Context, operationName string) opentracing.Span {
	span, spanFinish := tracingutils.StartChildSpan(pctx, operationName, nil)
	*pctx = setSpanToContext(*pctx, span, spanFinish)
	tracingutils.SetSpanServiceName(span, externalServiceName)
	tracingutils.SetSpanType(span, tracingutils.SpanTypeRedis)
	opentracing_ext.SpanKindRPCClient.Set(span)
	if len(h.Addrs) > 0 {
		span.SetTag("redis.addrs", strings.Join(h.Addrs, " "))
	}
	if h.DB != nil {
		span.SetTag("redis.db", *h.DB)
	}
	return span
}

func (h *Hook) setSpanTagKey(span opentracing.Span, cmd redis.Cmder) {
	args := cmd.Args()
	if len(args) == 0 {
		return
	}
	args = args[1:] // The first argument is the name of the command.
	switch cmd.Name() {
	// Single first
	case
		// Geo
		"geoadd", "geohash", "geopos", "geodist", "georadius", "georadiusbymember",
		// Hashes
		"hdel", "hexists", "hget", "hgetall", "hincrby", "hincrbyfloat", "hkeys", "hlen", "hmget", "hmset", "hset", "hsetnx", "hstrlen", "hvals", "hscan",
		// HyperLogLog
		"pfadd",
		// Keys
		"dump", "expire", "expireat", "move", "persist", "pexpire", "pexpireat", "pttl", "restore", "sort", "ttl", "type",
		// Lists
		"lindex", "linsert", "llen", "lpop", "lpos", "lpush", "lpushx", "lrange", "lrem", "lset", "ltrim", "rpop", "rpush", "rpushx",
		// Sets
		"sadd", "scard", "sismember", "smismember", "smembers", "spop", "srandmember", "srem", "sscan",
		// Sorted Sets
		"zadd", "zcard", "zcount", "zincrby", "zlexcount", "zpopmax", "zpopmin", "zrange", "zrangebylex", "zrevrangebylex", "zrangebyscore", "zrank", "zrem", "zremrangebylex", "zremrangebyrank", "zremrangebyscore", "zrevrange", "zrevrangebyscore", "zrevrank", "zscore", "zmscore", "zunionstore", "zcan",
		// Streams
		"xadd", "xtrim", "xdel", "xrange", "xrevrange", "xlen", "xack", "xclaim", "xpending",
		// Strings
		"append", "bitcount", "bitfield", "bitpos", "decr", "decrby", "get", "getbit", "getrange", "getset", "incr", "incrby", "incrbyfloat", "psetex", "set", "setbit", "setex", "setnx", "setrange", "strlen":
		h.setSpanTagKeySingleFirst(span, args)
	// Single offsetStart=1
	case
		// Streams
		"xinfo", "xgroup":
		h.setSpanTagKeySingleOffsetStart(span, args, 1)
	// Multi all
	case
		// HyperLogLog
		"pfcount", "pfmerge",
		// Keys
		"del", "exists", "rename", "renamenx", "touch", "unlink",
		// Lists
		"rpoplpush",
		// Sets
		"sdiff", "sdiffstore", "sinter", "sinterstore", "sunion", "sunionstore",
		// Strings
		"mget":
		h.setSpanTagKeyMultiAll(span, args)
	// Multi offsetStart=1
	case
		// Strings
		"bitop":
		h.setSpanTagKeyMultiOffset(span, args, 1, 0)
	// Multi offsetEnd=1
	case
		// Lists
		"blpop", "brpoplpush", "brpop",
		// Sets
		"smove",
		// Sorted Sets
		"bzpopmin", "bzpopmax":
		h.setSpanTagKeyMultiOffset(span, args, 0, 1)
	// Mset
	case "mset", "msetnx":
		h.setSpanTagKeyMset(span, args)
	}
}

func (h *Hook) setSpanTagKeySingleFirst(span opentracing.Span, args []interface{}) {
	if len(args) == 0 {
		return
	}
	arg := args[0]
	key, ok := arg.(string)
	if !ok {
		return
	}
	span.SetTag("redis.key", key)
}

func (h *Hook) setSpanTagKeySingleOffsetStart(span opentracing.Span, args []interface{}, offsetStart int) {
	if len(args) < offsetStart {
		return
	}
	args = args[offsetStart:]
	h.setSpanTagKeySingleFirst(span, args)
}

func (h *Hook) setSpanTagKeyMultiAll(span opentracing.Span, args []interface{}) {
	keys := make([]string, 0, len(args))
	for _, arg := range args {
		key, ok := arg.(string)
		if !ok {
			continue
		}
		keys = append(keys, key)
	}
	h.setSpanTagKeys(span, keys)
}

func (h *Hook) setSpanTagKeyMultiOffset(span opentracing.Span, args []interface{}, offsetStart int, offsetEnd int) {
	if len(args) < offsetStart+offsetEnd {
		return
	}
	args = args[offsetStart : len(args)-offsetEnd]
	h.setSpanTagKeyMultiAll(span, args)
}

func (h *Hook) setSpanTagKeyMset(span opentracing.Span, args []interface{}) {
	keys := make([]string, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		arg := args[i]
		key, ok := arg.(string)
		if !ok {
			continue
		}
		keys = append(keys, key)
	}
	h.setSpanTagKeys(span, keys)
}

func (h *Hook) setSpanTagKeys(span opentracing.Span, keys []string) {
	if len(keys) == 0 {
		return
	}
	span.SetTag("redis.keys", keys)
}

type spanContainer struct {
	span   opentracing.Span
	finish closeutils.F
}

type spanContainerContextKey struct{}

func setSpanToContext(ctx context.Context, span opentracing.Span, finish closeutils.F) context.Context {
	return context.WithValue(ctx, spanContainerContextKey{}, &spanContainer{
		span:   span,
		finish: finish,
	})
}

func getSpanFromContext(ctx context.Context) (opentracing.Span, closeutils.F) {
	sc := ctx.Value(spanContainerContextKey{}).(*spanContainer) //nolint:errcheck
	return sc.span, sc.finish
}
