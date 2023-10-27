//go:build !linux
// +build !linux

package web

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/librespeed/speedtest/config"
	log "github.com/sirupsen/logrus"
)

func startListener(conf *config.Config, r *chi.Mux) error {
	var s error

	addr := net.JoinHostPort(conf.BindAddress, conf.Port)
	log.Infof("Starting backend server on %s", addr)

	// Create a ListenConfig and enable MultipathTCP
	lc := &net.ListenConfig{}
	lc.SetMultipathTCP(true)

	// Create a listener using the ListenConfig
	listener, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return err
	}

	// TLS and HTTP/2.
	if conf.EnableTLS {
		log.Info("Use TLS connection.")
		if !(conf.EnableHTTP2) {
			srv := &http.Server{
				Addr:         addr,
				Handler:      r,
				TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			}
			s = srv.ServeTLS(listener, conf.TLSCertFile, conf.TLSKeyFile)
		} else {
			s = http.ServeTLS(listener, r, conf.TLSCertFile, conf.TLSKeyFile)
		}
	} else {
		if conf.EnableHTTP2 {
			log.Errorf("TLS is mandatory for HTTP/2. Ignore settings that enable HTTP/2.")
		}
		s = http.Serve(listener, r)
	}

	return s
}
