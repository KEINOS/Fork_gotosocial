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

package dereferencing

import (
	"net/url"
	"sync"

	"codeberg.org/gruf/go-mutexes"
	"github.com/superseriousbusiness/gotosocial/internal/media"
	"github.com/superseriousbusiness/gotosocial/internal/state"
	"github.com/superseriousbusiness/gotosocial/internal/transport"
	"github.com/superseriousbusiness/gotosocial/internal/typeutils"
)

// Dereferencer wraps logic and functionality for doing dereferencing
// of remote accounts, statuses, etc, from federated instances.
type Dereferencer struct {
	state               *state.State
	converter           *typeutils.Converter
	transportController transport.Controller
	mediaManager        *media.Manager
	derefAvatars        map[string]*media.ProcessingMedia
	derefAvatarsMu      mutexes.Mutex
	derefHeaders        map[string]*media.ProcessingMedia
	derefHeadersMu      mutexes.Mutex
	derefEmojis         map[string]*media.ProcessingEmoji
	derefEmojisMu       mutexes.Mutex
	handshakes          map[string][]*url.URL
	handshakesMu        sync.Mutex // mutex to lock/unlock when checking or updating the handshakes map
}

// NewDereferencer returns a Dereferencer initialized with the given parameters.
func NewDereferencer(
	state *state.State,
	converter *typeutils.Converter,
	transportController transport.Controller,
	mediaManager *media.Manager,
) Dereferencer {
	return Dereferencer{
		state:               state,
		converter:           converter,
		transportController: transportController,
		mediaManager:        mediaManager,
		derefAvatars:        make(map[string]*media.ProcessingMedia),
		derefHeaders:        make(map[string]*media.ProcessingMedia),
		derefEmojis:         make(map[string]*media.ProcessingEmoji),
		handshakes:          make(map[string][]*url.URL),

		// use wrapped mutexes to allow safely deferring unlock
		// even when more granular locks are required (only unlocks once).
		derefAvatarsMu: mutexes.WithSafety(mutexes.New()),
		derefHeadersMu: mutexes.WithSafety(mutexes.New()),
		derefEmojisMu:  mutexes.WithSafety(mutexes.New()),
	}
}
