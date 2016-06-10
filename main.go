package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/andrewstuart/soffit-go-poc/pkg/soffit"
	"github.com/dgrijalva/jwt-go"
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

			var jwt = "{{ .jwt }}"

			$.ajax({
				method: 'GET',
				url: '{{ .conf.endpoint }}/data',
				headers: {
					Authorization: "JWT " + jwt}
				})
				.then(function(d) {
					ele.find('#remote-data').append('<pre>' + d + '</pre>');
				});
		})(up.jQuery);
	</script>
</div>
`

func main() {
	t := template.Must(template.New("soffit").Parse(pageTpl))

	r := http.NewServeMux()

	secrets := map[string]string{}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		log.Println("Handling request")

		var sr soffit.Request

		err := json.NewDecoder(io.TeeReader(r.Body, os.Stdout)).Decode(&sr)
		if err != nil {
			log.Println("Error decoding JSON", err)
			http.Error(w, "Error decoding JSON", http.StatusNotAcceptable)
			return
		}

		bs, _ := json.MarshalIndent(sr, "", "  ")

		sec := newSecret()
		secrets[sr.UserName] = sec

		jwt, err := getJWT(sr, sec)

		if err != nil {
			log.Println("Error signing jwt", err)
			http.Error(w, "jwt signing error", 500)
			return
		}

		err = t.Execute(w, map[string]interface{}{
			"srJson": string(bs),
			"sr":     sr,
			"conf":   conf,
			"jwt":    jwt,
			"secret": sec,
		})

		if err != nil {
			log.Printf("Error executing template: %v\n", err)
		}
	})

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling %s request at /data\n", r.Method)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")

		if r.Method == http.MethodOptions {
			defer log.Println("finished options")

			w.Write(nil)
			return
		}

		au := strings.Split(r.Header.Get("Authorization"), " ")

		if len(au) < 2 || au[0] != "JWT" {
			log.Println("Only JWT Authorization is acceptable")
			http.Error(w, "Only JWT Authorization is acceptable", 401)
			return
		}

		reqJWT, err := jwt.Parse(au[1], func(t *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})

		if err != nil {
			log.Println("Error parsing jwt", err)
			http.Error(w, "Invalid JWT", 403)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"hello":        "world",
			"jwtSecret":    reqJWT.Claims["secret"].(string),
			"storedSecret": secrets[reqJWT.Claims["sub"].(string)],
		})
	})

	log.Fatal(http.ListenAndServe(":8089", r))
}

func newSecret() string {
	bs := make([]byte, 20)
	rand.Read(bs)

	return base64.StdEncoding.EncodeToString(bs)
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
