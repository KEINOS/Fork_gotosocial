// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apiutil "github.com/superseriousbusiness/gotosocial/internal/api/util"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
)

const textCSSUTF8 = string(apiutil.TextCSS + "; charset=utf-8")

func (m *Module) customCSSGETHandler(c *gin.Context) {
	if !config.GetAccountsAllowCustomCSS() {
		err := errors.New("accounts-allow-custom-css is not enabled on this instance")
		apiutil.WebErrorHandler(c, gtserror.NewErrorNotFound(err), m.processor.InstanceGetV1)
		return
	}

	if _, err := apiutil.NegotiateAccept(c, apiutil.TextCSS); err != nil {
		apiutil.WebErrorHandler(c, gtserror.NewErrorNotAcceptable(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	// usernames on our instance will always be lowercase
	username := strings.ToLower(c.Param(usernameKey))
	if username == "" {
		err := errors.New("no account username specified")
		apiutil.WebErrorHandler(c, gtserror.NewErrorBadRequest(err, err.Error()), m.processor.InstanceGetV1)
		return
	}

	customCSS, errWithCode := m.processor.Account().GetCustomCSSForUsername(c.Request.Context(), username)
	if errWithCode != nil {
		apiutil.WebErrorHandler(c, errWithCode, m.processor.InstanceGetV1)
		return
	}

	c.Header(cacheControlHeader, cacheControlNoCache)
	c.Data(http.StatusOK, textCSSUTF8, []byte(customCSS))
}
