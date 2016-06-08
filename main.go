package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
)

var scriptTpl = `
(function() {
	var i = 0;
	function incr() {
		$("#%s-response #ticktock").html('' + i++);
	}

	setInterval(incr, 1000);
})(up.jQuery);
`

func script(sr soffit.SoffitRequest) string {
	return fmt.Sprintf(scriptTpl, sr.Namespace)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var sr soffit.SoffitRequest

		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			log.Println("Error decoding JSON", err)
		}

		fmt.Fprintf(w, `<div id="%s-response">`, sr.Namespace)
		fmt.Fprintf(w, "<h1>Your portal is running %s major version %d</h1>", sr.Portal.Provider, sr.Portal.Version.Major)
		if bs, err := json.MarshalIndent(sr, "", "  "); err == nil {
			os.Stdout.Write(bs)
			fmt.Fprintf(w, "<h2>Full details:</h2>\n<pre>%s</pre>", bs)
		}
		fmt.Fprintf(w, `<script type="text/javascript">%s</script>`, script(sr))
		fmt.Fprintf(w, `<h2>Javascript example</h2>`)
		fmt.Fprintf(w, `<div id="ticktock"></div>`)
		fmt.Fprintf(w, `</div>`)
	})

	log.Fatal(http.ListenAndServe(":8089", nil))
}
