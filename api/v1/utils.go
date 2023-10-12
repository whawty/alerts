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

package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/whawty/alerts/store"
)

func statusCodeFromError(err error) (code int, response ErrorResponse) {
	code = http.StatusInternalServerError
	response = ErrorResponse{Error: err.Error()}

	switch err {
	case store.ErrNotImplemented:
		code = http.StatusNotImplemented
	case store.ErrNotFound:
		code = http.StatusNotFound
	}
	return
}

func sendError(c *gin.Context, err error) {
	code, response := statusCodeFromError(err)
	c.JSON(code, response)
}

func parsePositiveIntegerParameter(c *gin.Context, name string) (int, bool) {
	valueStr := c.Query(name)
	if valueStr == "" {
		return -1, true
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "query parameter " + name + " is invalid: " + err.Error()})
		return -1, false
	}
	if value < 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "query parameter " + name + " must be >= 0"})
		return -1, false
	}
	return value, true
}

func getPaginationParameter(c *gin.Context) (offset, limit int, ok bool) {
	if offset, ok = parsePositiveIntegerParameter(c, "offset"); !ok {
		return
	}
	limit, ok = parsePositiveIntegerParameter(c, "limit")
	return
}
