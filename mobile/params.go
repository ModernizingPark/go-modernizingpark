// Copyright 2016 The go-modernizingpark Authors
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

// Contains all the wrappers from the params package.

package gmpc

import (
	"github.com/modernizingpark/go-modernizingpark/p2p/discv5"
	"github.com/modernizingpark/go-modernizingpark/params"
)

// MainnetGenesis returns the JSON spec to use for the main Modernizingpark network. It
// is actually empty since that defaults to the hard coded binary genesis block.
func MainnetGenesis() string {
	return ""
}

// FoundationBootnodes returns the enode URLs of the P2P bootstrap nodes operated
// by the foundation running the V5 discovery protocol.
func FoundationBootnodes() *Enodes {
	nodes := &Enodes{nodes: make([]*discv5.Node, len(params.MainnetBootnodes))}
	for i, url := range params.MainnetBootnodes {
		nodes.nodes[i] = discv5.MustParseNode(url)
	}
	return nodes
}
