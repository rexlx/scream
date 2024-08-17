package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	flag.Parse()
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12, // Or tls.VersionTLS13 for stricter security
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,

			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		fmt.Println("error loading cert", err)
		return
	}
	cfg.Certificates = []tls.Certificate{cert}

	s := NewServer(*url, *firstUserMode)
	server := &http.Server{
		Addr:      *url,
		Handler:   s.Session.LoadAndSave(s.Gateway),
		TLSConfig: cfg,
	}
	s.CleanUpTokens()
	ticker := time.NewTicker(2 * time.Minute)
	go func() {
		for range ticker.C {
			s.UpdateGraphs()
		}
	}()
	s.Logger.Println("server started")
	fmt.Println("server started", s.URL)
	err = server.ListenAndServeTLS("", "")
	// err = server.ListenAndServeTLS(*certFile, *keyFile)
	if err != nil {
		fmt.Println("error starting server", err)
	}
}
