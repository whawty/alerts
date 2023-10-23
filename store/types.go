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

package store

import (
	"errors"
	"fmt"
	"time"

	"github.com/enescakir/emoji"
)

// Configuration

type Config struct {
	Path string `yaml:"path"`
}

// Errors

var (
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound       = errors.New("not found")
)

type ErrInvalidStateTransition struct {
	old AlertState
	new AlertState
}

func (e ErrInvalidStateTransition) Error() string {
	return fmt.Sprintf("invalid alert state transition: %s -> %s", e.old.String(), e.new.String())
}

// Alerts

type AlertState uint

const (
	StateNew AlertState = iota
	StateOpen
	StateAcknowledged
	StateStale
	StateClosed
)

func (s AlertState) String() string {
	switch s {
	case StateNew:
		return "new"
	case StateOpen:
		return "open"
	case StateAcknowledged:
		return "acknowledged"
	case StateStale:
		return "stale"
	case StateClosed:
		return "closed"
	}
	return "unknown"
}

func (s *AlertState) FromString(str string) error {
	switch str {
	case "new":
		*s = StateNew
	case "open":
		*s = StateOpen
	case "acknowledged":
		*s = StateAcknowledged
	case "stale":
		*s = StateStale
	case "closed":
		*s = StateClosed
	default:
		return errors.New("invalid alert state: '" + str + "'")
	}
	return nil
}

func (s AlertState) Emoji() emoji.Emoji {
	switch s {
	case StateNew:
		return emoji.GlowingStar
	case StateOpen:
		return emoji.Bell
	case StateAcknowledged:
		return emoji.BellWithSlash
	case StateStale:
		return emoji.QuestionMark
	case StateClosed:
		return emoji.CheckMarkButton
	}
	return emoji.WhiteQuestionMark
}

func (s AlertState) MarshalText() (data []byte, err error) {
	data = []byte(s.String())
	return
}

func (s *AlertState) UnmarshalText(data []byte) (err error) {
	return s.FromString(string(data))
}

type AlertSeverity uint

const (
	SeverityCritical AlertSeverity = iota
	SeverityWarning
	SeverityInformational
)

func (s AlertSeverity) String() string {
	switch s {
	case SeverityCritical:
		return "critical"
	case SeverityWarning:
		return "warning"
	case SeverityInformational:
		return "informational"
	}
	return "unknown"
}

func (s *AlertSeverity) FromString(str string) error {
	switch str {
	case "critical":
		*s = SeverityCritical
	case "warning":
		*s = SeverityWarning
	case "informational":
		*s = SeverityInformational
	default:
		return errors.New("invalid alert severity: '" + str + "'")
	}
	return nil
}

func (s AlertSeverity) Emoji() emoji.Emoji {
	switch s {
	case SeverityCritical:
		return emoji.DoubleExclamationMark
	case SeverityWarning:
		return emoji.Warning
	case SeverityInformational:
		return emoji.Information
	}
	return emoji.WhiteQuestionMark
}

func (s AlertSeverity) MarshalText() (data []byte, err error) {
	data = []byte(s.String())
	return
}

func (s *AlertSeverity) UnmarshalText(data []byte) (err error) {
	return s.FromString(string(data))
}

type Alert struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"created"`
	UpdatedAt time.Time     `json:"updated"`
	Name      string        `json:"name"`
	State     AlertState    `json:"state"`
	Severity  AlertSeverity `json:"severity"`
}

func (a Alert) String() string {
	return a.ID
}
