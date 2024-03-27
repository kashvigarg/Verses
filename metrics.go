package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiconfig) metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html;")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Verses Admin</h1>
	<p>Verses has been visited %d times!</p>
</body>

</html>
	`, cfg.fileservercounts)))
}

func (cfg *apiconfig) reqcounts(app http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileservercounts++
		app.ServeHTTP(w, r)
	})

}
