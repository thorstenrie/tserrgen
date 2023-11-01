// Copyright (c) 2023 thorstenrie.
// All Rights Reserved. Use is governed with GNU Affero General Public License v3.0
// that can be found in the LICENSE file.
package tserrgen_test

import (
	"testing"

	"github.com/thorstenrie/tserrgen"
)

func TestGenerate(t *testing.T) {
	if e := tserrgen.Generate("tserr.json"); e != nil {
		t.Error(e)
	}
}
