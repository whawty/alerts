//
// Copyright (c) 2023 whawty contributors (see AUTHORS file)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of whawty.alerts nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spreadspace/tlsconfig"
	apiV1 "github.com/whawty/alerts/api/v1"
	"github.com/whawty/alerts/store"
	"github.com/whawty/alerts/ui"
	"gopkg.in/yaml.v3"
)

const (
	WebUIPathPrefix = "/admin/"
	WebAPIv1Prefix  = "/api/v1/"
)

type webConfig struct {
	TLS *tlsconfig.TLSConfig `yaml:"tls"`
}

func readWebConfig(configfile string) (*webConfig, error) {
	file, err := os.Open(configfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	c := &webConfig{}
	if err = decoder.Decode(c); err != nil {
		return nil, fmt.Errorf("Error parsing config file: %s", err)
	}
	return c, nil
}

func runWeb(listener net.Listener, config *webConfig, st *store.Store) (err error) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.HandleMethodNotAllowed = true

	r.GET("/", func(c *gin.Context) { c.Redirect(http.StatusSeeOther, WebUIPathPrefix) })
	r.StaticFS(WebUIPathPrefix, ui.Assets)

	apiV1.InstallHTTPHandler(r.Group(WebAPIv1Prefix), st)

	server := &http.Server{Handler: r, WriteTimeout: 60 * time.Second, ReadTimeout: 60 * time.Second}
	if config != nil && config.TLS != nil {
		server.TLSConfig, err = config.TLS.ToGoTLSConfig()
		if err != nil {
			return
		}
		wl.Printf("web-api: listening on '%s' using TLS", listener.Addr())
		return server.ServeTLS(listener, "", "")

	}
	wl.Printf("web-api: listening on '%s'", listener.Addr())
	return server.Serve(listener)
}

func runWebAddr(addr string, config *webConfig, store *store.Store) (err error) {
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return runWeb(ln.(*net.TCPListener), config, store)
}

func runWebListener(listener *net.TCPListener, config *webConfig, store *store.Store) (err error) {
	return runWeb(listener, config, store)
}
