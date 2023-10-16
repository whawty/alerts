//
// Copyright (c) 2023 whawty contributors (see AUTHORS file)
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//   - Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//   - Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//   - Neither the name of whawty.alerts nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
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
	"log"

	"github.com/whawty/alerts/store"
)

type EMailBackend struct {
	infoLog *log.Logger
	dbgLog  *log.Logger
	name    string
	conf    *NotifierBackendConfigEMail
	// TODO: add client config
}

func NewEMailBackend(name string, conf *NotifierBackendConfigEMail, infoLog, dbgLog *log.Logger) *EMailBackend {
	return &EMailBackend{name: name, conf: conf, infoLog: infoLog, dbgLog: dbgLog}
}

func (emb *EMailBackend) Init() (err error) {
	return fmt.Errorf("not yet implemented!")
}

func (smb *EMailBackend) Ready() bool {
	// TODO: implement this
	return false
}

func (emb *EMailBackend) Notify(ctx context.Context, target NotifierTarget, alert *store.Alert) error {
	return fmt.Errorf("not yet implemented!")
}

func (emb *EMailBackend) Close() error {
	// TODO: close client?
	return nil
}
