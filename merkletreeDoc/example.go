package main

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/iden3/go-iden3/core"
	"github.com/iden3/go-iden3/db"
	"github.com/iden3/go-iden3/merkletree"
)

func main() {
	storage, err := db.NewLevelDbStorage("./path", false)
	if err != nil {
		panic(err)
	}
	// new merkletree of 140 levels of maximum depth using the defined
	// storage
	mt, err := merkletree.NewMerkleTree(storage, 140)
	if err != nil {
		panic(err)
	}
	defer mt.Storage().Close()

	name0 := "alice@iden3.io"
	ethAddr0 := common.HexToAddress("0x7b471a1bdbd3b8ac98f3715507449f3a8e1f3b22")
	claim0 := core.NewClaimAssignName(name0, ethAddr0)
	claimEntry0 := claim0.Entry()

	name1 := "bob@iden3.io"
	ethAddr1 := common.HexToAddress("0x28f8267fb21e8ce0cdd9888a6e532764eb8d52dd")
	claim1 := core.NewClaimAssignName(name1, ethAddr1)
	claimEntry1 := claim1.Entry()

	fmt.Println("adding claim0")
	err = mt.Add(claimEntry0)
	if err != nil {
		panic(err)
	}
	fmt.Println("merkle root: " + mt.RootKey().Hex())
	fmt.Println("adding claim1")
	err = mt.Add(claimEntry1)
	if err != nil {
		panic(err)
	}

	mp, err := mt.GenerateProof(claimEntry0.HIndex())
	if err != nil {
		panic(err)
	}
	fmt.Println("merkle root: " + mt.RootKey().Hex())

	fmt.Println("merkle proof: ", mp)
	checked := merkletree.VerifyProof(mt.RootKey(), mp,
		claimEntry0.HIndex(), claimEntry0.HValue())
	fmt.Println("merkle proof checked:", checked)

	claimDataInPos, err := mt.GetDataByIndex(claimEntry0.HIndex())
	if err != nil {
		panic(err)
	}
	claimEntryInPos := merkletree.Entry{Data: *claimDataInPos}
	// print true if the claimInPosBytes is the same than claimEntry0.Bytes()
	fmt.Println("claim in position equals to the original:",
		bytes.Equal(claimEntry0.Bytes(), claimEntryInPos.Bytes()))

	name2 := "eve@iden3.io"
	ethAddr2 := common.HexToAddress("0x29a6a240e2d8f8bf39b5338b9664d414c5d793f4")
	claim2 := core.NewClaimAssignName(name2, ethAddr2)
	claimEntry2 := claim2.Entry()

	mp, err = mt.GenerateProof(claimEntry2.HIndex())
	if err != nil {
		panic(err)
	}

	fmt.Println("merkle proof: ", mp)

	checked = merkletree.VerifyProof(mt.RootKey(), mp, claimEntry2.HIndex(), claimEntry2.HValue())

	fmt.Println("merkle proof of non existence checked:", checked)
}
