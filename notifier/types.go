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
	"time"

	"github.com/whawty/alerts/store"
)

type NotifierBackendConfigEMail struct {
	From      string `yaml:"from"`
	Smarthost string `yaml:"smarthost"`
	// TODO: add auth and TLS support
}

type NotifierBackendConfigSMSModem struct {
	Device   string        `yaml:"device"`
	Baudrate int           `yaml:"baudrate"`
	Timeout  time.Duration `yaml:"timeout"`
	Pin      *uint         `yaml:"pin"`
}

type NotifierBackendConfig struct {
	Name     string
	EMail    *NotifierBackendConfigEMail    `yaml:"email"`
	SMSModem *NotifierBackendConfigSMSModem `yaml:"smsModem"`
}

type NotifierTargetSMS string
type NotifierTargetEMail string

type NotifierTarget struct {
	Name  string               `yaml:"name"`
	EMail *NotifierTargetEMail `yaml:"email"`
	SMS   *NotifierTargetSMS   `yaml:"sms"`
}

type Config struct {
	Interval time.Duration           `yaml:"interval"`
	Backends []NotifierBackendConfig `yaml:"backends"`
	Targets  []NotifierTarget        `yaml:"targets"`
}

// Interfaces

type NotifierBackend interface {
	Init() error
	Ready() bool
	Notify(context.Context, NotifierTarget, *store.Alert) (bool, error)
	Close() error
}
