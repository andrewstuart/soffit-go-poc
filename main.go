package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"astuart.co/soffit-go"
	"astuart.co/vpki"

	"github.com/dgrijalva/jwt-go"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	conf map[string]string

	useVault = flag.Bool("use-vault", false, "use vault to obtain a certificate")
)

func init() {
	flag.Parse()
}

func main() {
	t := template.Must(template.ParseGlob("templates/**"))

	r := http.NewServeMux()

	secrets := map[string]string{}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		log.Println("Handling request")

		var sr soffit.Payload

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

		rs, err := soffit.RandHTMLID()
		if err != nil {
			log.Println(err)
		}

		err = t.Lookup("soffit.tmpl.html").Execute(w, map[string]interface{}{
			"srJson": string(bs),
			"sr":     sr,
			"conf":   conf,
			"jwt":    jwt,
			"secret": sec,
			"rs":     rs,
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

	r.Handle("/metrics", prometheus.Handler())

	if *useVault {
		cli := &vpki.Client{
			Addr:  "vault.astuart.co",
			Mount: "pki",
			Role:  "kube",
		}

		cli.SetToken(os.Getenv("VAULT_TOKEN"))

		go func() {
			err := vpki.ListenAndServeTLS(":8443", prometheus.InstrumentHandler("soffit", r), cli)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

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
