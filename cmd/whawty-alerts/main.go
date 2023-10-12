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
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/urfave/cli"
	"github.com/whawty/alerts/store"
)

var (
	wl  = log.New(os.Stdout, "[whawty.alerts]\t", log.LstdFlags)
	wdl = log.New(ioutil.Discard, "[whawty.alerts dbg]\t", log.LstdFlags)
)

func init() {
	if _, exists := os.LookupEnv("WHAWTY_ALERTS_DEBUG"); exists {
		wdl.SetOutput(os.Stderr)
	}
}

func cmdRun(c *cli.Context) error {
	s, err := store.Open("test.db")
	if err != nil {
		return cli.NewExitError(err.Error(), 3)
	}

	var webc *webConfig
	if c.String("web-config") != "" {
		webc, err = readWebConfig(c.String("web-config"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}
	webAddrs := c.StringSlice("web-addr")

	var wg sync.WaitGroup
	for _, webAddr := range webAddrs {
		a := webAddr
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := runWebAddr(a, webc, s); err != nil {
				fmt.Printf("warning running web interface(%s) failed: %s\n", a, err)
			}
		}()
	}
	wg.Wait()

	return cli.NewExitError(fmt.Sprintf("shutting down since all listening sockets have closed."), 0)
}

func main() {
	app := cli.NewApp()
	app.Name = "whawty-alerts"
	app.Version = "0.1"
	app.Usage = "simple alert manager"
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "run the alert manager",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "web-config",
					Value:  "",
					Usage:  "path to the web configuration file",
					EnvVar: "WHAWTY_ALERTS_WEB_CONFIG",
				},
				cli.StringSliceFlag{
					Name:   "web-addr",
					Usage:  "address to listen on for web API",
					EnvVar: "WHAWTY_ALERTS_WEB_ADDR",
				},
			},
			Action: cmdRun,
		},
	}

	wdl.Printf("calling app.Run()")
	app.Run(os.Args)
}
