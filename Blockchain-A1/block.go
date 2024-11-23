package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// BlockHeader represents the header of a block
type BlockHeader struct {
	Version      int
	PreviousHash string
	MerkleRoot   string
	Timestamp    string
}

// Transaction represents a transaction containing the dataset and algorithm CID and data
type Transaction struct {
	DatasetCID    string
	AlgorithmCID  string
	DatasetData   []byte
	AlgorithmData []byte
}

// Block represents a block in the blockchain
type Block struct {
	Index       int
	Header      BlockHeader
	Transaction Transaction
	OutputHash  string
	Hash        string
	Nonce       int
}

// CalculateHash generates the hash for a block
func calculateHash(block Block) string {
	record := fmt.Sprintf("%d%s%s%s%s%s%d", block.Index, block.Header.PreviousHash, block.Header.Timestamp, block.Header.MerkleRoot, block.Transaction.DatasetCID, block.Transaction.AlgorithmCID, block.Nonce)
	hash := sha256.New()
	hash.Write([]byte(record))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// MineBlock mines a new block
func mineBlock(lastBlock Block, transaction Transaction) Block {
	nonce := 0
	var newBlock Block
	newBlock.Index = lastBlock.Index + 1
	newBlock.Header.Version = 1
	newBlock.Header.PreviousHash = lastBlock.Hash
	newBlock.Header.Timestamp = time.Now().String()
	newBlock.Transaction = transaction
	newBlock.Header.MerkleRoot = "placeholder_merkle_root"

	// Run Python script with the dataset to get the output hash
	output, err := runPythonScript("script.py", "dataset.csv")
	if err != nil {
		panic(err)
	}
	// Store the output hash in the block
	outputHash := fmt.Sprintf("%x", sha256.Sum256(output))
	newBlock.OutputHash = outputHash

	// Find a valid nonce
	for {
		newBlock.Nonce = nonce
		hash := calculateHash(newBlock)
		if hash[:4] == "0000" {
			newBlock.Hash = hash
			break
		}
		nonce++
	}

	return newBlock
}

// VerifyBlock verifies the integrity of a block
func verifyBlock(block Block) bool {
	// Fetch the dataset and algorithm data from IPFS
	datasetData, err := fetchFromIPFS(block.Transaction.DatasetCID)
	if err != nil {
		return false
	}

	algorithmData, err := fetchFromIPFS(block.Transaction.AlgorithmCID)
	if err != nil {
		return false
	}

	err = writeToFile("dataset_from_block.csv", datasetData)
	if err != nil {
		return false
	}

	err = writeToFile("algorithm_from_block.py", algorithmData)
	if err != nil {
		return false
	}

	// Run the Python script with the fetched dataset
	output, err := runPythonScript("algorithm_from_block.py", "dataset_from_block.csv")
	if err != nil {
		return false
	}

	// Compare the output hash with the block's output hash
	computedOutputHash := fmt.Sprintf("%x", sha256.Sum256(output))
	return computedOutputHash == block.OutputHash
}
