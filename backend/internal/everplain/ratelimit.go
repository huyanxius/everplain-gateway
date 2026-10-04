package everplain

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// One atomic, expiring counter per service key, shared across gateway replicas.
// Redis failure is fail-closed. A TTL window can burst at its boundary and is
// not a cost reservation or a hard token budget.
var requestCounter = redis.NewScript(`
local count = redis.call('INCR', KEYS[1])
if count == 1 then redis.call('PEXPIRE', KEYS[1], ARGV[1]) end
local ttl = redis.call('PTTL', KEYS[1])
if ttl < 0 then redis.call('PEXPIRE', KEYS[1], ARGV[1]); ttl = tonumber(ARGV[1]) end
return {count, ttl}
`)

func RateLimit(client redis.UniversalClient, perMinute int, keyID func(*gin.Context) int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := keyID(c)
		if id <= 0 {
			Abort(c, http.StatusUnauthorized, "authentication_failed", "A valid server-side service key is required.", false)
			return
		}
		if client == nil {
			Abort(c, http.StatusServiceUnavailable, "limiter_unavailable", "The service-key rate limiter is unavailable.", true)
			return
		}
		values, err := requestCounter.Run(c.Request.Context(), client, []string{"everplain:rpm:" + strconv.FormatInt(id, 10)}, time.Minute.Milliseconds()).Int64Slice()
		if err != nil || len(values) != 2 {
			Abort(c, http.StatusServiceUnavailable, "limiter_unavailable", "The service-key rate limiter is unavailable.", true)
			return
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(perMinute))
		if values[0] > int64(perMinute) {
			retry := max(int64(1), (values[1]+999)/1000)
			c.Header("Retry-After", strconv.FormatInt(retry, 10))
			Abort(c, http.StatusTooManyRequests, "rate_limited", "The service-key request limit was exceeded.", true)
			return
		}
		c.Next()
	}
}
