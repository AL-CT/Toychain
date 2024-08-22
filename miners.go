package main

import (
	"fmt"
	"math"
	"sync"
)

type MiningPool struct {
	sync.Mutex
	chain     *Blockchain
	currIndex int
	count     uint64
	currPos   uint64
}

func (mp *MiningPool) startWorker(blk *Block, index int, workerID int) {

	for mp.currPos < math.MaxUint64 && index == mp.currIndex {
		mp.Lock()
		if mp.currPos >= math.MaxUint64 { //different worker reached maximum size of nonce
			mp.Unlock()
			return
		}
		start := mp.currPos
		mp.currPos = start + mp.count
		mp.Unlock()
		err := blk.mineBlock(start, mp.count)
		//if err != nil {
		//fmt.Printf("Worker %v - "+err.Error()+"\n", workerID)
		if err == nil {
			mp.Lock()
			if index != mp.currIndex { //different worker added block to chain
				mp.Unlock()
				return
			}
			err = mp.chain.addBlock(blk)
			if err != nil {
				fmt.Printf("Worker %v - "+err.Error()+"\n", workerID)
			} else {
				mp.currPos = 0
				mp.currIndex++
			}
			mp.Unlock()
		}

	}
	return
}
