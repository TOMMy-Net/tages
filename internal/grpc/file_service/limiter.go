package fileservice

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc/peer"
)

type LimitsCheker struct {
	limit int64
	value atomic.Int64
}

type Limits struct {
	DownloadAndUploadConn *LimitsCheker
	ListFilesConn         *LimitsCheker
}

type Limiter struct {
	mu     sync.Mutex
	limits map[string]*Limits // key: HOST
}

func NewLimiter() *Limiter {
	return &Limiter{
		limits: make(map[string]*Limits, 10),
	}
}

func NewLimits() *Limits {
	return &Limits{
		DownloadAndUploadConn: &LimitsCheker{
			limit: 10,
		},
		ListFilesConn: &LimitsCheker{
			limit: 100,
		},
	}
}

func (l *Limiter) GetLimits(host string) *Limits {
	l.mu.Lock()
	defer l.mu.Unlock()
	if v, ok := l.limits[host]; ok {
		return v
	}

	l.limits[host] = NewLimits()
	return l.limits[host]

}

func (l *LimitsCheker) Add() bool {
	if v := l.value.Load(); v < l.limit {
		l.value.Add(1)
		return true
	}
	return false
}

func (l *LimitsCheker) Release() bool {
	if v := l.value.Load(); v > 0 {
		l.value.Add(-1)
		return true
	}
	return false
}


func getClientIP(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return "unknown"
	}
	host, _, _ := net.SplitHostPort(p.Addr.String())
	return host
}
