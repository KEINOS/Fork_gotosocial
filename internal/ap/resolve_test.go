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

package ap_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/superseriousbusiness/gotosocial/internal/ap"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
)

type ResolveTestSuite struct {
	APTestSuite
}

func (suite *ResolveTestSuite) TestResolveDocumentAsStatusable() {
	b := []byte(suite.typeToJson(suite.document1))

	statusable, err := ap.ResolveStatusable(context.Background(), b)
	suite.NoError(err)
	suite.NotNil(statusable)
}

func (suite *ResolveTestSuite) TestResolveDocumentAsAccountable() {
	b := []byte(suite.typeToJson(suite.document1))

	accountable, err := ap.ResolveAccountable(context.Background(), b)
	suite.True(gtserror.WrongType(err))
	suite.EqualError(err, "ResolveAccountable: cannot resolve vocab type *typedocument.ActivityStreamsDocument as accountable")
	suite.Nil(accountable)
}

func TestResolveTestSuite(t *testing.T) {
	suite.Run(t, &ResolveTestSuite{})
}
