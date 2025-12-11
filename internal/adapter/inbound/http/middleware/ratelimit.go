package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/okinn/service-presensi/pkg/httputil"
	"golang.org/x/time/rate"
)

type RateLimiterConfig struct {
	// RequestsPerSecond adalah jumlah request yang diizinkan per detik
	RequestsPerSecond rate.Limit
	// BurstSize adalah jumlah request yang diizinkan dalam burst
	BurstSize int
	// CleanupInterval adalah interval untuk membersihkan limiter yang tidak aktif
	CleanupInterval time.Duration
	// MaxAge adalah durasi maksimal limiter disimpan setelah terakhir digunakan
	MaxAge time.Duration
}

func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         20,
		CleanupInterval:   time.Minute,
		MaxAge:            3 * time.Minute,
	}
}

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	config   RateLimiterConfig
}

func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		config:   config,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.config.RequestsPerSecond, rl.config.BurstSize)
		rl.visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.config.CleanupInterval)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.config.MaxAge {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			httputil.Error(w, http.StatusTooManyRequests, "Terlalu banyak request, coba lagi nanti")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// LoginRateLimiter adalah rate limiter khusus untuk endpoint login
// dengan limit yang lebih ketat untuk mencegah brute force
type LoginRateLimiter struct {
	*RateLimiter
}

func NewLoginRateLimiter() *LoginRateLimiter {
	config := RateLimiterConfig{
		RequestsPerSecond: 1, // 1 request per second
		BurstSize:         5, // Max 5 attempts dalam burst
		CleanupInterval:   time.Minute,
		MaxAge:            5 * time.Minute,
	}
	return &LoginRateLimiter{
		RateLimiter: NewRateLimiter(config),
	}
}
