package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
	"strconv"
	"time"
)

const diff = 20 //Difficulty

type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      []byte
	PrevHash  []byte
	Nonce     uint64
}

type Blockchain struct {
	blocks []*Block
}

func calculateHash(block Block, nonce uint64) []byte {
	record := string(rune(block.Index)) + block.Timestamp + block.Data + string(block.PrevHash) + strconv.Itoa(int(nonce))
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hashed
}

func generateBlock(index int, prevHash []byte, Data string) Block {

	newBlock := Block{
		Index:     index,
		Timestamp: time.Now().String(),
		Data:      Data,
		PrevHash:  prevHash,
	}
	/*
		err := newBlock.mineBlock()
		if err != nil {
			return newBlock, errors.New("error: could not create a block")
		}*/

	return newBlock
}

func (blk *Block) mineBlock(start uint64, count uint64) error {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-diff))
	nonce := start
	for nonce < start+count && nonce < math.MaxUint64 {

		hash := calculateHash(*blk, nonce)

		hashInt := new(big.Int)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			blk.Nonce = nonce
			blk.Hash = hash
			return nil
		} else {
			nonce++
		}
	}
	if nonce >= math.MaxUint64 {
		return errors.New("error: block has no solution") //Could play around with extra nonce, etc.}
	} else {
		return errors.New("error: could not find solution within given nonce range")
	}

}

func (chain *Blockchain) addBlock(blk *Block) error {
	var prevHash []byte
	if chain.getSize() == 0 {
		var temp [sha256.Size]byte
		prevHash = temp[:]
	}
	if chain.getSize() > 0 {
		prevHash = chain.blocks[len(chain.blocks)-1].Hash
	}
	if err := verify(blk, prevHash); err != nil {
		return err
	}
	chain.blocks = append(chain.blocks, blk)
	return nil
}

func (chain *Blockchain) getBlock(index int) (*Block, error) {
	if index >= chain.getSize() {
		return &Block{}, errors.New("error: index out of range")
	}
	return chain.blocks[index], nil
}

func verify(newBlock *Block, prevHash []byte) error {

	if !bytes.Equal(prevHash, newBlock.PrevHash) {
		//fmt.Printf("Hash1: %x\nHash2: %x\n", prevHash, newBlock.PrevHash)
		return errors.New("error: previous hash is different")
	}
	if !bytes.Equal(calculateHash(*newBlock, newBlock.Nonce), newBlock.Hash) {
		return errors.New("error: Hash with nonce is invalid")
	}
	return nil
}

func (chain Blockchain) getSize() int {
	return len(chain.blocks)
}

// to be used with multiple chains being received
func (chain Blockchain) replaceChain(newBlocks *Blockchain) error {
	if len(newBlocks.blocks) > len(chain.blocks) {
		chain.blocks = newBlocks.blocks
		return nil
	}
	return errors.New("error: received chain is shorter than the original chain")
}
