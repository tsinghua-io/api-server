package util

import (
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Limiter struct {
	sync.Mutex
	Window    int64 // In seconds.
	RateLimit int
	Count     map[string]int
	Reset     int64
}

func NewLimiter(window int64, rate int) *Limiter {
	return &Limiter{
		Window:    window,
		RateLimit: rate,
	}
}

func (l *Limiter) Request(id string) (ok bool, remain int, reset int64) {
	l.Lock()

	// Reset if needed.
	if now := time.Now().UTC().Unix(); now > l.Reset {
		l.Reset = now + l.Window
		l.Count = make(map[string]int)
	}
	reset = l.Reset

	count := l.Count[id]
	if count >= l.RateLimit {
		ok = false
	} else {
		ok = true
		count++
		l.Count[id] = count
	}
	remain = l.RateLimit - count

	l.Unlock()
	return
}

func (l *Limiter) Handler() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			host, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				glog.Errorf("Failed to get remote host from %s: %s", req.RemoteAddr, err)
				// Logged it, let it go.
			}
			ok, remain, reset := l.Request(host)

			rw.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.RateLimit))
			rw.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remain))
			rw.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset, 10))

			if !ok {
				err := fmt.Sprintf("API rate limit exceeded for %s.", host)
				Error(rw, err, http.StatusTooManyRequests)
			} else {
				h.ServeHTTP(rw, req)
			}
		})
	}
}
