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
	"io"
	"log"
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/serial"
	"github.com/whawty/alerts/store"
)

type SMSModemBackend struct {
	infoLog *log.Logger
	dbgLog  *log.Logger
	name    string
	conf    *NotifierBackendConfigSMSModem
	modem   io.ReadWriteCloser
	sms     *gsm.GSM
}

func NewSMSModemBackend(name string, conf *NotifierBackendConfigSMSModem, infoLog, dbgLog *log.Logger) *SMSModemBackend {
	if conf.Timeout <= 0 {
		conf.Timeout = 5 * time.Second
	}
	return &SMSModemBackend{name: name, conf: conf, infoLog: infoLog, dbgLog: dbgLog}
}

func (smb *SMSModemBackend) Init() (err error) {
	smb.modem, err = serial.New(serial.WithPort(smb.conf.Device), serial.WithBaud(smb.conf.Baudrate))
	if err != nil {
		smb.modem = nil
		return
	}

	a := at.New(smb.modem, at.WithTimeout(smb.conf.Timeout))
	if smb.conf.Pin != nil {
		var resp []string
		resp, err = a.Command(fmt.Sprintf("+CPIN=%d", smb.conf.Pin))
		if err != nil {
			smb.modem.Close()
			smb.modem = nil
			return
		}
		smb.dbgLog.Printf("SMSModem(%s): enter pin code response: %v", smb.name, resp)
	}

	smb.sms = gsm.New(a)
	err = smb.sms.Init()
	if err != nil {
		smb.modem.Close()
		smb.modem = nil
		smb.sms = nil
		return
	}

	err = smb.sms.StartMessageRx(
		func(msg gsm.Message) {
			smb.infoLog.Printf("SMSModem(%s): got SMS from '%s': %s", smb.name, msg.Number, msg.Message)
		},
		func(err error) {
			smb.infoLog.Printf("SMSModem(%s): got SMS rx error: %v", smb.name, err)
		})

	if err != nil {
		smb.modem.Close()
		smb.modem = nil
		smb.sms = nil
		return
	}
	return nil
}

func (smb *SMSModemBackend) Ready() bool {
	return smb.modem != nil && smb.sms != nil
}

func (smb *SMSModemBackend) Notify(ctx context.Context, target NotifierTarget, alert *store.Alert) error {
	// TODO: improve alert formatting
	message := fmt.Sprintf("%s / %s / %s", alert.State, alert.Severity, alert.Name)

	resp, err := smb.sms.SendLongMessage(target.SMS.Number, message)
	if err != nil {
		return err
	}
	smb.dbgLog.Printf("SMSModem(%s): send sms response: %v", smb.name, resp)
	return nil
}

func (smb *SMSModemBackend) Close() error {
	smb.sms.StopMessageRx()
	smb.modem.Close()
	smb.modem = nil
	smb.sms = nil
	return nil
}
