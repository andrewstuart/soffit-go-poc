package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"astuart.co/vpki"

	"github.com/Masterminds/sprig"
	"github.com/dgrijalva/jwt-go"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	conf map[string]string

	useVault = flag.Bool("use-vault", false, "use vault to obtain a certificate")
)

const (
	saltLen = 8
	ivLen   = saltLen
)

func init() {
	flag.Parse()
}

type SoffitOpts struct {
	Preferences, Definition, Request map[string]interface{}
}

func main() {
	t := template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseGlob("templates/**.html"))

	r := http.NewServeMux()

	secrets := map[string]string{}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		s := SoffitOpts{}

		for k := range r.Header {
			if strings.Index(k, "X-Soffit") != 0 {
				continue
			}
			bs, err := base64.StdEncoding.DecodeString(r.Header.Get(k))
			if err != nil {
				log.Println(err)
				continue
			}
			dec, err := decrypt(bs, "CHANGEME")
			if err != nil {
				log.Println(err)
				continue
			}

			token, err := jwt.Parse(string(dec), nil)

			if err != nil && !strings.Contains(err.Error(), "Keyfunc") {
				log.Println("Error parsing jwt", err)
				http.Error(w, "Invalid JWT", 403)
				return
			}

			switch k {
			case "X-Soffit-Portalrequest":
				s.Request = token.Claims
			case "X-Soffit-Definition":
				s.Definition = token.Claims
			case "X-Soffit-Preferences":
				s.Preferences = token.Claims
			}

		}

		err := t.Lookup("soffit.tmpl.html").Execute(w, s)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error parsing template", 500)
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

		if err != nil && !strings.Contains(err.Error(), "Keyfunc") {
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

	log.Println("listening")

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
