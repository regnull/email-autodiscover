package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/regnull/email-autodiscover/template"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

const (
	HTTP_PORT  = 80
	HTTPS_PORT = 443
)

type CmdArgs struct {
	HttpPort   int
	HttpsPort  int
	CertFile   string
	KeyFile    string
	ConfigFile string
}

type Server struct {
	args *template.Args
}

func NewServer(templateArgs *template.Args) *Server {
	return &Server{args: templateArgs}
}

func (s *Server) HandleThunderbirdConfig(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("thunderbird config request")
	reply, err := template.Thunderbird(s.args)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate Thunderbird response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(reply))
}

func (s *Server) HandleIOSConfig(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("iOS config request")
	email := r.URL.Query().Get("email")
	if email == "" {
		log.Warn().Msg("no email in iOS mail config request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug().Str("email", email).Msg("got email")

	newArgs := *s.args
	newArgs.Email = email
	reply, err := template.IOSMail(&newArgs)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate iOS response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(reply))
}

func (s *Server) HandleOutlookConfig(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("outlook config request")
	// Request comes as:
	// `<Autodiscover xmlns="https://schemas.microsoft.com/exchange/autodiscover/outlook/requestschema/2006">
	// <Request>
	//   <EMailAddress>user@contoso.com</EMailAddress>
	//   <AcceptableResponseSchema>https://schemas.microsoft.com/exchange/autodiscover/outlook/responseschema/2006a</AcceptableResponseSchema>
	// </Request>
	// </Autodiscover>`

	if r.Method != "POST" {
		log.Warn().Msg("not POST request to outlook config")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("failed to read request body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var xmlRec struct {
		XMLName xml.Name `xml:"Autodiscover"`
		Request struct {
			XMLName                  xml.Name `xml:"Request"`
			EMailAddress             string   `xml:"EMailAddress"`
			AcceptableResponseSchema string   `xml:"AcceptableResponseSchema"`
		}
	}
	if err := xml.Unmarshal(body, &xmlRec); err != nil {
		log.Warn().Err(err).Msg("failed to parse request xml")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debug().Str("schema", xmlRec.Request.AcceptableResponseSchema).Msg("got schema")
	newArgs := *s.args
	newArgs.Email = xmlRec.Request.EMailAddress
	reply, err := template.OutlookMail(&newArgs)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate iOS response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-type", "text/xml")
	w.Write([]byte(reply))
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05", NoColor: true})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var args CmdArgs
	flag.IntVar(&args.HttpPort, "http-port", HTTP_PORT, "HTTP port")
	flag.IntVar(&args.HttpsPort, "https-port", HTTPS_PORT, "HTTPS port")
	flag.StringVar(&args.CertFile, "cert-file", "", "certificate file")
	flag.StringVar(&args.KeyFile, "key-file", "", "key file")
	flag.StringVar(&args.ConfigFile, "config", "", "config file")
	flag.Parse()

	if args.ConfigFile == "" {
		log.Fatal().Msg("--config-file must be specified")
	}

	log.Info().Str("config-file", args.ConfigFile).Msg("using config file")
	config, err := ioutil.ReadFile(args.ConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config file")
	}

	templateArgs := &template.Args{}
	err = yaml.Unmarshal(config, templateArgs)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse config file")
	}

	log.Debug().Interface("config", templateArgs).Msg("got config")

	server := NewServer(templateArgs)
	http.HandleFunc("/mail/config-v1.1.xml", server.HandleThunderbirdConfig)
	http.HandleFunc("/email.mobileconfig", server.HandleIOSConfig)
	http.HandleFunc("/Autodiscover/Autodiscover.xml", server.HandleOutlookConfig)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Int("port", args.HttpPort).Msg("starting http server...")
		err := http.ListenAndServe(fmt.Sprintf(":%d", args.HttpPort), nil)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	if args.CertFile != "" && args.KeyFile != "" {
		wg.Add(1)

		go func() {
			defer wg.Done()
			log.Info().Int("port", args.HttpsPort).Msg("starting https server...")
			err := http.ListenAndServeTLS(fmt.Sprintf(":%d", args.HttpsPort), args.CertFile, args.KeyFile, nil)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}()
	}
	wg.Wait()
}
