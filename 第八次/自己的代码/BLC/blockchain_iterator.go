package BLC

import (
	"github.com/boltdb/bolt"
	"log"
)

type Rwq_BlockchainIterator struct {
	rwq_currentHash []byte
	rwq_db          *bolt.DB
}

func (i *Rwq_BlockchainIterator) Rwq_Next() *Rwq_Block {
	var block *Rwq_Block

	err := i.rwq_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.rwq_currentHash)
		block = Rwq_DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.rwq_currentHash = block.Rwq_PrevBlockHash

	return block
}