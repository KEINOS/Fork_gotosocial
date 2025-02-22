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

package typeutils

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/superseriousbusiness/gotosocial/internal/ap"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/regexes"
)

type statusInteractions struct {
	Faved      bool
	Muted      bool
	Bookmarked bool
	Reblogged  bool
	Pinned     bool
}

func (c *Converter) interactionsWithStatusForAccount(ctx context.Context, s *gtsmodel.Status, requestingAccount *gtsmodel.Account) (*statusInteractions, error) {
	si := &statusInteractions{}

	if requestingAccount != nil {
		faved, err := c.state.DB.IsStatusFavedBy(ctx, s.ID, requestingAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("error checking if requesting account has faved status: %s", err)
		}
		si.Faved = faved

		reblogged, err := c.state.DB.IsStatusBoostedBy(ctx, s.ID, requestingAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("error checking if requesting account has reblogged status: %s", err)
		}
		si.Reblogged = reblogged

		muted, err := c.state.DB.IsStatusMutedBy(ctx, s, requestingAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("error checking if requesting account has muted status: %s", err)
		}
		si.Muted = muted

		bookmarked, err := c.state.DB.IsStatusBookmarkedBy(ctx, s, requestingAccount.ID)
		if err != nil {
			return nil, fmt.Errorf("error checking if requesting account has bookmarked status: %s", err)
		}
		si.Bookmarked = bookmarked

		// The only time 'pinned' should be true is if the
		// requesting account is looking at its OWN status.
		if s.AccountID == requestingAccount.ID {
			si.Pinned = !s.PinnedAt.IsZero()
		}
	}
	return si, nil
}

func misskeyReportInlineURLs(content string) []*url.URL {
	m := regexes.MisskeyReportNotes.FindAllStringSubmatch(content, -1)
	urls := make([]*url.URL, 0, len(m))
	for _, sm := range m {
		url, err := url.Parse(sm[1])
		if err == nil && url != nil {
			urls = append(urls, url)
		}
	}
	return urls
}

// getURI is a shortcut/util function for extracting
// the JSONLDId URI of an Activity or Object.
func getURI(withID ap.WithJSONLDId) (*url.URL, string, error) {
	idProp := withID.GetJSONLDId()
	if idProp == nil {
		return nil, "", errors.New("id prop was nil")
	}

	if !idProp.IsIRI() {
		return nil, "", errors.New("id prop was not an IRI")
	}

	id := idProp.Get()
	return id, id.String(), nil
}
