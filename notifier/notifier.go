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

package notifier

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/whawty/alerts/store"
)

type Notifier struct {
	conf     *Config
	store    *store.Store
	infoLog  *log.Logger
	dbgLog   *log.Logger
	ctx      context.Context
	cancel   context.CancelFunc
	backends map[string]NotifierBackend
}

func (n *Notifier) Close() error {
	n.cancel()
	for _, backend := range n.backends {
		if backend.Ready() {
			backend.Close()
		}
	}
	return nil
}

func NewNotifier(conf *Config, st *store.Store, infoLog, dbgLog *log.Logger) (n *Notifier, err error) {
	if infoLog == nil {
		infoLog = log.New(io.Discard, "", 0)
	}
	if dbgLog == nil {
		dbgLog = log.New(io.Discard, "", 0)
	}

	n = &Notifier{conf: conf, store: st, infoLog: infoLog, dbgLog: dbgLog}
	n.ctx, n.cancel = context.WithCancel(context.Background())
	if n.conf.Interval <= 0 {
		n.conf.Interval = 1 * time.Minute
	}

	n.backends = make(map[string]NotifierBackend)
	for idx, backend := range n.conf.Backends {
		if backend.Name == "" {
			err = fmt.Errorf("found unnamed backend at config index %d", idx)
			return
		}
		if _, exists := n.backends[backend.Name]; exists {
			err = fmt.Errorf("found duplicate backend name at config index %d", idx)
			return
		}

		var b NotifierBackend
		cnt := 0
		if backend.EMail != nil {
			b = NewEMailBackend(backend.Name, backend.EMail, infoLog, dbgLog)
			cnt = cnt + 1
		}
		if backend.SMSModem != nil {
			b = NewSMSModemBackend(backend.Name, backend.SMSModem, infoLog, dbgLog)
			cnt = cnt + 1
		}
		if cnt == 0 {
			err = fmt.Errorf("no valid backend config found for backend '%s'", backend.Name)
			return
		}
		if cnt > 1 {
			err = fmt.Errorf("backend '%s' has ambiguous backend config", backend.Name)
			return
		}
		n.backends[backend.Name] = b

		if err := b.Init(); err != nil {
			infoLog.Printf("notifier: failed to initialize backend '%s': %v", backend.Name, err)
		} else {
			infoLog.Printf("notifier: backend '%s' successfully initialized", backend.Name)
		}
	}

	// TODO: start go-routine to re-initialize failed backends
	// TODO: start go-routine to handle notfications

	infoLog.Printf("notifier: started with %d backends and evaluation interval %s", len(n.backends), conf.Interval.String())
	return
}
