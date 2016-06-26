package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var conf map[string]string

//<h1>Your portal is running {{ .sr.Request.Portal.Provider }} major version {{ .sr.Request.Portal.Version.Major }}.</h1>
const pageTpl = `
<div id="{{ .sr.Request.Namespace }}-response">
	<h2>Hello there {{ .sr.User.Username }}. Welcome to Soffit.</h2>
	<h2>Full Details:</h2>
	<pre>{{ .srJson }}</pre>

	<h2>JavaScript example</h2>
	<div id="ticktock"></div>

	<h2>Remote Data Example</h2>
	<div id="remote-data"></div>

	<script type="text/javascript">
		(function() {
			var i = 0;
			var ele = $('#{{ .sr.Request.Namespace }}-response');
			function incr() {
				 ele.find('#ticktock').html('' + i++);
			}

			setInterval(incr, 1000);

			var jwt = "{{ .jwt }}"

			$.ajax({
				method: 'GET',
				url: '{{ .conf.endpoint }}/data',
				headers: {
					Authorization: "Bearer " + jwt}
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

	tl, err := tls.Listen("tcp", ":8443", &tls.Config{
		Certificates: []tls.Certificate{*tlsCert},
	})
	if err != nil {
		log.Fatal(err)
	}

	r := http.NewServeMux()

	secrets := map[string]string{}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		log.Println("Handling request")

		var sr Payload

		err := json.NewDecoder(io.TeeReader(r.Body, os.Stdout)).Decode(&sr)
		if err != nil {
			log.Println("Error decoding JSON", err)
			http.Error(w, "Error decoding JSON", http.StatusNotAcceptable)
			return
		}

		bs, _ := json.MarshalIndent(sr, "", "  ")

		sec := newSecret()
		secrets[sr.User.Username] = sec

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

		if len(au) < 2 || au[0] != "Bearer" {
			log.Println("Only Bearer Authorization is acceptable")
			http.Error(w, "Only Bearer Authorization is acceptable", 401)
			return
		}

		reqJWT, err := jwt.Parse(au[1], func(t *jwt.Token) (interface{}, error) {
			return &signingKey.PublicKey, nil
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

	go func() {
		err := http.Serve(tl, r)
		if err != nil {
			log.Fatal(err)
		}
	}()

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
