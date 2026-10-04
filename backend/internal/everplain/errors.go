package everplain

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// errorWriter discards only pre-stream HTTP error bodies, so provider errors
// cannot leak prompts, tokens or account details. Successful SSE is never
// buffered and remains the upstream protocol, including its terminal errors.
type errorWriter struct {
	gin.ResponseWriter
	status int
}

func (w *errorWriter) WriteHeader(code int) {
	if w.ResponseWriter.Written() {
		return
	}
	w.status = code
	if code < 400 {
		w.ResponseWriter.WriteHeader(code)
	}
}
func (w *errorWriter) Status() int {
	if w.status != 0 {
		return w.status
	}
	return w.ResponseWriter.Status()
}
func (w *errorWriter) WriteHeaderNow() {
	if w.Status() < 400 {
		w.ResponseWriter.WriteHeaderNow()
	}
}
func (w *errorWriter) Write(b []byte) (int, error) {
	if w.Status() >= 400 && !w.ResponseWriter.Written() {
		return len(b), nil
	}
	return w.ResponseWriter.Write(b)
}
func (w *errorWriter) WriteString(s string) (int, error) { return w.Write([]byte(s)) }
func (w *errorWriter) Flush() {
	if w.Status() < 400 {
		w.ResponseWriter.Flush()
	}
}

func NormalizeErrors(enabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enabled || !IsGatewayPath(c.Request.URL.Path) {
			c.Next()
			return
		}
		original := c.Writer
		w := &errorWriter{ResponseWriter: original}
		c.Writer = w
		c.Header("X-Everplain-Contract-Version", ContractVersion)
		c.Next()
		c.Writer = original
		if w.Status() < 400 || original.Written() {
			return
		}
		e := ClassifyError(w.Status())
		if v, ok := c.Get("everplain_error"); ok {
			if known, ok := v.(Error); ok {
				e = known
			}
		}
		e.RequestID = original.Header().Get("X-Client-Request-ID")
		// Recompute framing: the original provider body has been discarded.
		for _, key := range []string{"Content-Length", "Content-Encoding", "ETag"} {
			original.Header().Del(key)
		}
		original.Header().Set("Cache-Control", "no-store")
		original.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.JSON(w.Status(), gin.H{"error": e})
	}
}

var _ http.ResponseWriter = (*errorWriter)(nil)
