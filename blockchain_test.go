package main

import (
	"crypto/sha256"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

var chain Blockchain

func (blk Block) print(t *testing.T) {

	t.Logf("#########################\n")
	t.Logf("Index: 		%v\n", blk.Index)
	t.Logf("Timestamp: 	%s\n", blk.Timestamp)
	t.Logf("Data: 		%s\n", blk.Data)
	t.Logf("Hash: 		%x\n", blk.Hash[:])
	t.Logf("PrevHash: 	%x\n", blk.PrevHash[:])
	t.Logf("Nonce: 		%x\n", blk.Nonce)
	t.Logf("#########################\n\n")

}

const MiningRange = 200000
const MaxWorkers = 5
const ChainSize = 8

func TestCreateBlockchain(t *testing.T) {

	var wg sync.WaitGroup

	chain = Blockchain{make([]*Block, 0)}
	mp := MiningPool{chain: &chain, currIndex: 0, count: MiningRange, currPos: 0}
	var temp [sha256.Size]byte
	newBlock := generateBlock(mp.currIndex, temp[:], "Genesis")
	index := mp.currIndex
	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mp.startWorker(&newBlock, index, i)
		}()
	}

	wg.Wait()

	chain.blocks[chain.getSize()-1].print(t)

	nameArr := []string{"John", "Alice", "Bob", "Charlie", "Marzia", "Peter", "Shiv"}
	var message string

	for i := 0; i < ChainSize-1; i++ {

		i1 := rand.Intn(len(nameArr))
		i2 := rand.Intn(len(nameArr))
		for i2 == i1 {
			i2 = rand.Intn(len(nameArr))
		}

		message = nameArr[i1] + " sends " + nameArr[i2] + " " + strconv.Itoa(rand.Intn(100)+1) + " LD"
		newBlock := generateBlock(mp.currIndex, mp.chain.blocks[mp.chain.getSize()-1].Hash, message)
		index = mp.currIndex
		for i := 0; i < MaxWorkers-1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mp.startWorker(&newBlock, index, i)
			}()
		}
		wg.Wait()
		newBlock.print(t)
	}
}

/*func TestCreateBlockchain(t *testing.T) {

	chain, err = createBlockchain("Alice sends Bob 10 LD")
	if err != nil {
		t.Fatal(err)
	}
	newBlock, _ := generateBlock(chain.blocks[len(chain.blocks)-1], "John sends Alice 20 LD")
	if err := verify(newBlock, *chain.blocks[len(chain.blocks)-1]); err != nil {
		t.Fatal(err.Error())
	} else {
		chain.blocks = append(chain.blocks, &newBlock)
	}
	newBlock1, _ := generateBlock(chain.blocks[len(chain.blocks)-1], "Alice sends John 5 LD")
	if errVer := verify(newBlock1, *chain.blocks[len(chain.blocks)-1]); errVer != nil {
		t.Fatal(errVer.Error())
	} else {
		chain.blocks = append(chain.blocks, &newBlock1)
	}
	for _, blk := range chain.blocks {
		blk.print(t)
	}
}

func TestRejectInvalidBlock(t *testing.T) {
	newBlock, _ := generateBlock(chain.blocks[len(chain.blocks)-2], "Bob sends John 30 LD")
	if err := verify(newBlock, *chain.blocks[len(chain.blocks)-1]); err != nil {
		t.Logf(err.Error())
	} else {
		t.Fatal("Invalid block was accepted!")
	}
}*/
