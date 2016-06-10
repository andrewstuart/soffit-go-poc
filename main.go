package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
)

var conf map[string]string

const pageTpl = `
<div id="{{ .sr.Namespace }}-response">
	<h1>Your portal is running {{ .sr.Portal.Provider }} major version {{ .sr.Portal.Version.Major }}.</h1>
	<h2>Full Details:</h2>
	<pre>{{ .srJson }}</pre>

	<h2>JavaScript example</h2>
	<div id="ticktock"></div>

	<h2>Remote Data Example</h2>
	<div id="remote-data"></div>

	<script type="text/javascript">
		(function() {
			var i = 0;
			var ele = $('#{{ .sr.Namespace }}-response');
			function incr() {
				 ele.find('#ticktock').html('' + i++);
			}

			setInterval(incr, 1000);

			$.get('{{ .conf.endpoint }}/data')
				.success(function(d) {
					ele.find('#remote-data').append('<pre>' + d + '</pre>');
				});
		})(up.jQuery);
	</script>
</div>
`

func main() {
	t := template.Must(template.New("soffit").Parse(pageTpl))

	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			"conf":   conf,
		})

		if err != nil {
			log.Println("Error executing template: %v", err)
		}
	})

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		json.NewEncoder(w).Encode(map[string]string{
			"hello": "world",
		})
	})

	log.Fatal(http.ListenAndServe(":8089", r))
}

func init() {
	// bs, err := ioutil.ReadFile(os.Getenv("CONF_FILE"))
	// if err != nil {
	// 	log.Fatal("Could not read CONF_FILE", err)
	// }
	// err = yaml.Unmarshal(bs, &conf)
	// if err != nil {
	// 	log.Fatal("Could not unmarshal yaml config", err)
	// }
	//
	// conf["endpoint"] = os.Getenv("ENDPOINT")

	conf = map[string]string{
		"endpoint": os.Getenv("ENDPOINT"),
	}
}
