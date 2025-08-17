package dash

import (
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`
<!DOCTYPE html>
<head>
	<meta charset="utf-8">
	<title>Gophormula</title>
</head>
<body>
	<h1>dash</h1>
</body>
</html>
`))
	if err != nil {
		log.Println(err)
	}
}
