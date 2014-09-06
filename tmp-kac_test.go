// Copyright (c) 2014-2014 PPCD developers.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package ppcutil_test

import (
	"github.com/conformal/btcdb"
	"github.com/conformal/btcnet"
	"github.com/mably/ppcutil"
	"testing"

	_ "github.com/conformal/btcdb/ldb" // init only
)

func TestPoWTargetCalculation(t *testing.T) {
	params := btcnet.MainNetParams
	db, err := btcdb.OpenDB("leveldb", "testdata/db_512")
	if err != nil {
		t.Errorf("db error %v", err)
		return
	}
	defer db.Close()

	lastBlock, _ := db.FetchBlockBySha(params.GenesisHash)
	for height := 1; height < 512; height++ {
		sha, _ := db.FetchBlockShaByHeight(int64(height))
		block, _ := db.FetchBlockBySha(sha)
		if !block.MsgBlock().IsProofOfStake() {
			targetRequired := ppcutil.GetNextTargetRequired(params, db, lastBlock, false)
			if targetRequired != block.MsgBlock().Header.Bits {
				t.Errorf("bad target for block #%d %v, have %x want %x", height, sha, targetRequired, block.MsgBlock().Header.Bits)
				return
			}
		}
		lastBlock = block
	}
	return
}
