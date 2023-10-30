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
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/whawty/alerts/store"
)

func (api *API) ListHeartbeats(c *gin.Context) {
	offset, limit, ok := getPaginationParameter(c)
	if !ok {
		return
	}

	heartbeats, err := api.store.ListHeartbeats(offset, limit)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, HeartbeatsListing{heartbeats})
}

func (api *API) CreateHeartbeat(c *gin.Context) {
	heartbeat := &store.Heartbeat{}
	err := json.NewDecoder(c.Request.Body).Decode(heartbeat)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "error decoding heartbeat: " + err.Error()})
		return
	}
	heartbeat.ID = ulid.Make().String()

	if heartbeat, err = api.store.CreateHeartbeat(heartbeat); err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusCreated, heartbeat)
}

func (api *API) ReadHeartbeat(c *gin.Context) {
	id := c.Param("heartbeat-id")

	heartbeat, err := api.store.GetHeartbeat(id)
	if err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusOK, heartbeat)
}

func (api *API) DeleteHeartbeat(c *gin.Context) {
	id := c.Param("heartbeat-id")

	if err := api.store.DeleteHeartbeat(id); err != nil {
		sendError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
