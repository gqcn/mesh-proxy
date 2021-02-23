package httpproxy

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

// ResponseWriter is the custom writer for http response.
type ResponseWriter struct {
	status      int                 // HTTP status.
	writer      http.ResponseWriter // The underlying ResponseWriter.
	buffer      *bytes.Buffer       // The output buffer.
	hijacked    bool                // Mark this request is hijacked or not.
	wroteHeader bool                // Is header wrote or not, avoiding error: superfluous/multiple response.WriteHeader call.
}

// NewResponseWriter creates and return a ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		buffer: bytes.NewBuffer(nil),
		writer: w,
	}
}

// RawWriter returns the underlying ResponseWriter.
func (w *ResponseWriter) RawWriter() http.ResponseWriter {
	return w.writer
}

// Status returns the status of ResponseWriter.
func (w *ResponseWriter) Status() int {
	return w.status
}

// Header implements the interface function of http.ResponseWriter.Header.
func (w *ResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write implements the interface function of http.ResponseWriter.Write.
func (w *ResponseWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	return len(data), nil
}

// WriteHeader implements the interface of http.ResponseWriter.WriteHeader.
func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
}

// Hijack implements the interface function of http.Hijacker.Hijack.
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	w.hijacked = true
	return w.writer.(http.Hijacker).Hijack()
}

// BufferString returns the buffered content as []byte.
func (w *ResponseWriter) Buffer() []byte {
	return w.buffer.Bytes()
}

// BufferString returns the buffered content as string.
func (w *ResponseWriter) BufferString() string {
	return w.buffer.String()
}

// OutputBuffer outputs the buffer to client and clears the buffer.
func (w *ResponseWriter) OutputBuffer() {
	if w.hijacked {
		return
	}
	if w.status != 0 && !w.wroteHeader {
		w.writer.WriteHeader(w.status)
	}
	// Default status text output.
	if w.status != http.StatusOK && w.buffer.Len() == 0 {
		w.buffer.WriteString(http.StatusText(w.status))
	}
	if w.buffer.Len() > 0 {
		w.writer.Write(w.buffer.Bytes())
		w.buffer.Reset()
	}
}
