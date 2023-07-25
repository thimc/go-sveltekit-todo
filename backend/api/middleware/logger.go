package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			log.Println(r.RemoteAddr, r.Method, r.RequestURI, time.Since(start))

			// TODO: The way the logger is setup currently makes this impossible,
			//       the Context is always empty since the user isn't passed down
			//       to the Logger middleware.
			//
			// user := r.Context().Value("user").(*types.User)
			// if user != nil {
			// 	log.Printf("%+v\n", user)
			// }

		}(time.Now())
		next.ServeHTTP(w, r)
	})
}
