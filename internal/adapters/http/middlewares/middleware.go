package middlewares

import (
	"log"
	"net/http"
	"regexp"
	"time"
)

func cleanURLPath(path string) string {
	uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	codeRegex := regexp.MustCompile(`[0-9]{8,15}`)
	postSlugRegex := regexp.MustCompile(`/cms(/v3|)/posts/([\w\p{L}\p{N}_-]+)`)

	path = uuidRegex.ReplaceAllString(path, "{uuid}")
	path = codeRegex.ReplaceAllString(path, "{code}")
	path = postSlugRegex.ReplaceAllString(path, "/cms/posts/{slug}")
	return path
}

type Middleware func(http.Handler) http.Handler

func MiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, v := range xs {
			next = v(next)
		}
		return next
	}
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := time.Now()
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.URL.Path, time.Since(n))
	})
}
