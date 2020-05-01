package main

import (
	"github.com/chrisvdg/gotiny/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	listAddr := pflag.StringP("listenaddr", "l", ":8080", "http listen address")
	tlsListAddr := pflag.StringP("tlsaddr", "t", "8443", "https listen address")
	tlsKey := pflag.StringP("tlskey", "k", "", "TLS private key file path")
	tlsCert := pflag.StringP("tlscert", "c", "", "TLS certificate file path")
	readToken := pflag.StringP("readtoken", "r", "", "Read authorization token")
	writeToken := pflag.StringP("writetoken", "w", "", "Write authorization token")
	allowPublicCreate := pflag.BoolP("allowpubliccreate", "p", false, "Allows creation of generated tiny URLs without authorization when write token is set")
	idLen := pflag.IntP("idlen", "i", 5, "Length of generated tiny URL IDs")
	prettyJSON := pflag.BoolP("prettyjson", "j", false, "API outputs more readable JSON")
	fileBackendPath := pflag.StringP("filebackend", "f", "", "File to store file backend data")
	verbose := pflag.BoolP("verbose", "v", false, "Verbose output")

	pflag.Parse()

	c := &server.Config{
		ListenAddr:    *listAddr,
		TLSListenAddr: *tlsListAddr,
		TLS: &server.TLSConfig{
			KeyFile:  *tlsKey,
			CertFile: *tlsCert,
		},
		ReadAuthToken:              *readToken,
		WriteAuthToken:             *writeToken,
		AllowPublicCreateGenerated: *allowPublicCreate,
		GeneratedIDLen:             *idLen,
		PrettyJSON:                 *prettyJSON,
		FileBackendPath:            *fileBackendPath,
		Verbose:                    *verbose,
	}

	if c.Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05 02/01/2006",
	})

	s, err := server.New(c)
	if err != nil {
		log.Fatalf("Failed to init server: %s", err)
	}
	err = s.ListenAndServeFileBackedAPI()
	if err != nil {
		log.Fatalf("Failed to run server: %s", err)
	}
}
