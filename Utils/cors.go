package utils

import "net/http"

const (
	options           string = "OPTIONS"
	allow_origin      string = "Access-Control-Allow-Origin"
	allow_methods     string = "Access-Control-Allow-Methods"
	allow_headers     string = "Access-Control-Allow-Headers"
	allow_credentials string = "Access-Control-Allow-Credentials"
	expose_headers    string = "Access-Control-Expose-Headers"
	credentials       string = "true"
	origin            string = "Origin"
	methods           string = "POST, GET, OPTIONS, PUT, DELETE, HEAD, PATCH"

	headers string = "Access-Control-Allow-Origin, Accept, Accept-Encoding, Authorization, Content-Length, Content-Type, X-CSRF-Token,ApiKey"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if o := r.Header.Get(origin); o != "" {
			w.Header().Set(allow_origin, o)
		} else {
			w.Header().Set(allow_origin, "*")
		}

		w.Header().Set(allow_headers, headers)
		w.Header().Set(allow_methods, methods)
		w.Header().Set(allow_credentials, credentials)
		w.Header().Set(expose_headers, headers)

		if r.Method == options {
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}
