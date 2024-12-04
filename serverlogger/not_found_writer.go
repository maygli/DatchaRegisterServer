package serverlogger

import "net/http"

type NotFoundResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (writer *NotFoundResponseWriter) WriteHeader(status int) {
	writer.Status = status // Store the status for our own use
	if status != http.StatusNotFound {
		writer.ResponseWriter.WriteHeader(status)
	}
}

func (writer *NotFoundResponseWriter) Write(data []byte) (int, error) {
	if writer.Status != http.StatusNotFound {
		return writer.ResponseWriter.Write(data)
	}
	return len(data), nil
}

//Have to has 'Flush' function to implemenet Flusher interface
func (writer *NotFoundResponseWriter) Flush() {
	if flusher, ok := writer.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
