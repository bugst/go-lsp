//
// Copyright 2021 Cristian Maglie. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package lsp

type Unimplemented struct{}

func (*Unimplemented) UnmarshalJSON([]byte) error {
	panic("Unimplemented!")
}
