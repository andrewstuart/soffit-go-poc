package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
)

const pageTpl = `
<div id="{{ .sr.Namespace }}-response">
	<h1>Your portal is running {{ .sr.Portal.Provider }} major version {{ .sr.Portal.Version.Major }}.</h1>
	<h2>Full Details:</h2>
	<pre>{{ .srJson }}</pre>

	<h2>JavaScript example</h2>
	<div id="ticktock"></div>

	<script type="text/javascript">
		(function() {
			var i = 0;
			function incr() {
				$("#{{ .sr.Namespace }}-response #ticktock").html('' + i++);
			}

			setInterval(incr, 1000);
		})(up.jQuery);
	</script>
</div>
`

func main() {
	t := template.Must(template.New("soffit").Parse(pageTpl))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		log.Println("Handling request")

		var sr soffit.Request

		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			log.Println("Error decoding JSON", err)
		}

		bs, _ := json.MarshalIndent(sr, "", "  ")

		err = t.Execute(w, map[string]interface{}{
			"srJson": string(bs),
			"sr":     sr,
		})

		if err != nil {
			log.Println("Error executing template: %v", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8089", nil))
}
