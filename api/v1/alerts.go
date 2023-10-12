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

	"github.com/gin-gonic/gin"
)

func (api *API) ListAlerts(c *gin.Context) {
	offset, limit, ok := getPaginationParameter(c)
	if !ok {
		return
	}

	alerts, err := api.store.ListAlerts(offset, limit)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, AlertsListing{alerts})
}

func (api *API) ReadAlert(c *gin.Context) {
	id := c.Param("alert-id")

	alert, err := api.store.GetAlert(id)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, alert)
}

func (api *API) DeleteAlert(c *gin.Context) {
	id := c.Param("alert-id")

	if err := api.store.DeleteAlert(id); err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
