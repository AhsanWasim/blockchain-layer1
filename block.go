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

	transactions := []string{transaction.DatasetCID, transaction.AlgorithmCID}
	newBlock.Header.MerkleRoot = computeMerkleRoot(transactions)

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
	transactions := []string{block.Transaction.DatasetCID, block.Transaction.AlgorithmCID}
	computedMerkleRoot := computeMerkleRoot(transactions)

	if err != nil {
		return false
	}

	if computedMerkleRoot != block.Header.MerkleRoot {
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

func computeMerkleRoot(transactions []string) string {
	if len(transactions) == 0 {
		return ""
	}

	var hashPairs []string

	// leaf nodes are being hashed for the merkle tree (leaf nodes contain transactions)
	for _, tx := range transactions {
		hash := sha256.Sum256([]byte(tx))
		hashPairs = append(hashPairs, fmt.Sprintf("%x", hash))
	}

	// hashes are combined as a pair and then combined hash is calculated for parent node.
	// for odd num of hashes, the odd number moves up without pairing
	// it happens till the root
	// in the end merkle root is returned

	for len(hashPairs) > 1 {
		var newLevel []string
		for i := 0; i < len(hashPairs); i += 2 {
			if i+1 < len(hashPairs) {
				combinedHash := sha256.Sum256([]byte(hashPairs[i] + hashPairs[i+1]))
				newLevel = append(newLevel, fmt.Sprintf("%x", combinedHash))
			} else {
				newLevel = append(newLevel, hashPairs[i])
			}
		}
		hashPairs = newLevel
	}
	fmt.Printf("\nMerkle Root: %s\n", hashPairs[0])

	return hashPairs[0]
}
