// Copyright 2019 The go-modernizingpark Authors
// This file is part of the go-modernizingpark library.
//
// The go-modernizingpark library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-modernizingpark library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-modernizingpark library. If not, see <http://www.gnu.org/licenses/>.

package les

import (
	"github.com/modernizingpark/go-modernizingpark/p2p"
	"github.com/modernizingpark/go-modernizingpark/p2p/dnsdisc"
	"github.com/modernizingpark/go-modernizingpark/p2p/enode"
	"github.com/modernizingpark/go-modernizingpark/rlp"
)

// lesEntry is the "les" ENR entry. This is set for LES servers only.
type lesEntry struct {
	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e lesEntry) ENRKey() string {
	return "les"
}

// setupDiscovery creates the node discovery source for the eth protocol.
func (eth *LightEthereum) setupDiscovery(cfg *p2p.Config) (enode.Iterator, error) {
	if /*cfg.NoDiscovery || */ len(eth.config.DiscoveryURLs) == 0 {
		return nil, nil
	}
	client := dnsdisc.NewClient(dnsdisc.Config{})
	return client.NewIterator(eth.config.DiscoveryURLs...)
}
