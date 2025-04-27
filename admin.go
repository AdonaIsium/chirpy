package main

import (
	"html/template"
	"net/http"
)

type MetricsData struct {
	Hits int32
}

func (cfg *apiConfig) adminMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	const tpl = `
<!DOCTYPE html>
<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited {{.Hits}} times!</p>
</body>
</html>`

	tmpl := template.Must(template.New("page").Parse(tpl))
	data := MetricsData{
		Hits: cfg.fileserverHits.Load(),
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
