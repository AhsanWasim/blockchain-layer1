package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Blockchain struct {
	Blocks []Block
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(newBlock Block) {
	bc.Blocks = append(bc.Blocks, newBlock)
}

// GetLastBlock gets the last block in the blockchain
func (bc *Blockchain) GetLastBlock() Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

// SaveBlockchain saves the blockchain to a file
func saveBlockchain(filePath string, bc Blockchain) error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling blockchain: %w", err)
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

// LoadBlockchain loads the blockchain from a file
func loadBlockchain(filePath string) (Blockchain, error) {
	var bc Blockchain

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return bc, fmt.Errorf("error reading blockchain file: %w", err)
	}

	if err := json.Unmarshal(file, &bc); err != nil {
		return bc, fmt.Errorf("error unmarshalling blockchain data: %w", err)
	}

	return bc, nil
}
