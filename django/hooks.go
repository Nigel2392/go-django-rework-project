package django

import (
	"net/http"

	"github.com/Nigel2392/django/core/except"
)

type markedResponseWriter struct {
	http.ResponseWriter
	wasWritten bool
}

func (w *markedResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *markedResponseWriter) WriteHeader(statusCode int) {
	w.wasWritten = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *markedResponseWriter) Write(data []byte) (int, error) {
	w.wasWritten = true
	return w.ResponseWriter.Write(data)
}

type (
	DjangoHook      func(*Application) error
	ServerErrorHook func(w http.ResponseWriter, r *http.Request, app *Application, err except.ServerError)
)

const (
	HOOK_SERVER_ERROR = "django.ServerError"
)
