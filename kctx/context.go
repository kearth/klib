// 提供增强版上下文管理，基于标准 context.Context 扩展，
// 支持元数据键值对存储、追踪ID（TraceID）自动生成与继承，
// 并保证并发安全的读写操作。
//
// 核心特性：
// - 完全兼容标准 context 接口，可无缝替换原生 context
// - 内置 TraceID 用于分布式追踪，支持从父上下文继承
// - 提供线程安全的 Set/Get 方法管理元数据
// - 支持 WithCancel/WithTimeout 等衍生上下文创建
package kctx

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// TraceIDKey 用于在上下文中存储/获取 TraceID 的键，
	// 可通过 context.Value(TraceIDKey) 从父上下文继承 TraceID。
	TraceIDKey = "TraceID"
	// MetaMapKey 用于在上下文中存储/获取元数据映射的键，
	// 通过 context.Value(MetaMapKey) 可获取所有元数据的副本。
	MetaMapKey = "MetaMap"
)

type (

	// Context 扩展标准 context.Context 接口，增加元数据和追踪ID管理能力。
	// 实现了标准库 context.Context 的所有方法，可直接作为标准上下文使用。
	Context interface {
		context.Context
		Get(key string) string
		Set(key string, val string)
		Values() map[string]string
		Context() context.Context
		SetContext(ctx context.Context)
		TraceID() string
	}

	// 上下文实现
	kCtx struct {
		ctx     context.Context   // 底层标准context
		traceID string            // 不可变TraceID（无需锁保护）
		meta    map[string]string // 字符串元数据映射
		mu      sync.RWMutex      // 保护 metamap 和 context 的并发访问
	}
)

// 创建上下文
func New(parent ...context.Context) Context {
	// 解析父上下文，默认使用Background
	baseCtx := context.Background()
	if len(parent) > 0 && parent[0] != nil {
		baseCtx = parent[0]
	}
	// 继承或生成TraceID：优先从父上下文获取，其次生成新UUID
	traceID := uuid.NewString()
	if parentTraceID, ok := baseCtx.Value(TraceIDKey).(string); ok && parentTraceID != "" {
		traceID = parentTraceID
	}

	meta := make(map[string]string)
	if parentMeta, ok := baseCtx.Value(MetaMapKey).(map[string]string); ok {
		meta = parentMeta
	}

	return &kCtx{
		ctx:     baseCtx,
		traceID: traceID,
		meta:    meta,
	}
}

func (k *kCtx) Get(key string) string {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.meta[key]
}

func (k *kCtx) Set(key string, val string) {
	// 空键直接忽略，避免无效数据
	if key == "" {
		return
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	// 过滤"值未变更"的更新，减少map复制开销
	if k.meta[key] == val {
		return
	}

	// 预分配容量（原长度+1），避免扩容，提升复制效率
	newMeta := make(map[string]string, len(k.meta)+1)
	for k, v := range k.meta {
		newMeta[k] = v
	}
	newMeta[key] = val
	k.meta = newMeta
}

// Values 返回元数据副本，彻底杜绝外部修改内部状态
func (k *kCtx) Values() map[string]string {
	k.mu.RLock()
	defer k.mu.RUnlock()

	copied := make(map[string]string, len(k.meta))
	for k, v := range k.meta {
		copied[k] = v
	}
	return copied
}

// Context 获取底层标准context（并发安全）
func (k *kCtx) Context() context.Context {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.ctx
}

// TraceID 返回不可变TraceID（无需锁，初始化后不再修改）
func (k *kCtx) TraceID() string {
	return k.traceID
}

// SetContext 替换底层标准context，增加nil校验
func (k *kCtx) SetContext(ctx context.Context) {
	if ctx == nil {
		return
	}

	k.mu.Lock()
	defer k.mu.Unlock()
	k.ctx = ctx
}

// --------------- 实现context.Context接口 ---------------
func (k *kCtx) Done() <-chan struct{} {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.ctx.Done()
}

func (k *kCtx) Err() error {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.ctx.Err()
}

func (k *kCtx) Deadline() (deadline time.Time, ok bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.ctx.Deadline()
}

func (k *kCtx) Value(key any) any {
	// 优先处理内部关键键，避免锁开销
	keyStr, isStrKey := key.(string)
	if isStrKey {
		switch keyStr {
		case TraceIDKey:
			return k.traceID
		case MetaMapKey:
			return k.Values() // 返回副本，安全无副作用
		}
	}

	// 其他键从底层context获取，加锁保证并发安全
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.ctx.Value(key)
}

// --------------- 内部辅助函数 ---------------
// copyCtx 深拷贝kCtx实例，仅在内部调用，保证入参为*kCtx
func copyCtx(src *kCtx) *kCtx {
	src.mu.RLock()
	defer src.mu.RUnlock()

	// 深拷贝元数据（string不可变，直接赋值即为深拷贝）
	newMeta := make(map[string]string, len(src.meta))
	for k, v := range src.meta {
		newMeta[k] = v
	}

	return &kCtx{
		ctx:     src.ctx,     // 复用底层context（引用类型，符合context设计理念）
		traceID: src.traceID, // TraceID不可变，直接复用
		meta:    newMeta,     // 元数据深拷贝，隔离变更
	}
}

// --------------- 上下文衍生函数 ---------------
// WithCancel 基于父上下文创建可取消的新上下文，并发安全
func WithCancel(parent Context) (Context, context.CancelFunc) {
	// 类型断言失败时，使用默认上下文兜底
	parentImpl, ok := parent.(*kCtx)
	if !ok {
		newCtx := New()
		baseCtx, cancel := context.WithCancel(newCtx.Context())
		newCtx.SetContext(baseCtx)
		return newCtx, cancel
	}

	// 1. 基于父上下文创建可取消的标准context
	parentImpl.mu.RLock()
	baseCtx, cancel := context.WithCancel(parentImpl.ctx)
	parentImpl.mu.RUnlock()

	// 2. 复制父上下文元数据，创建新实例
	newCtx := copyCtx(parentImpl)
	newCtx.SetContext(baseCtx)

	return newCtx, cancel
}

// WithTimeout 包级函数：基于父上下文创建带超时的新上下文（核心调整点）
func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	// 处理无效超时（≤0），降级为可取消上下文
	if timeout <= 0 {
		return WithCancel(parent)
	}

	parentImpl, ok := parent.(*kCtx)
	if !ok {
		newCtx := New()
		baseCtx, cancel := context.WithTimeout(newCtx.Context(), timeout)
		newCtx.SetContext(baseCtx)
		return newCtx, cancel
	}

	parentImpl.mu.RLock()
	baseCtx, cancel := context.WithTimeout(parentImpl.ctx, timeout)
	parentImpl.mu.RUnlock()

	newCtx := copyCtx(parentImpl)
	newCtx.SetContext(baseCtx)
	return newCtx, cancel
}
