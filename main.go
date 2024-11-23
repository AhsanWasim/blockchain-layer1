// Testing changes

package main

import (
	"log"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		bc, err := loadBlockchain("blockchain.json")
		if err != nil {
			log.Println("No existing blockchain found. Creating a new one.")
			bc = Blockchain{}
		}

		pythonScriptCID := "QmSYAcrNoFFvvd81yuUaCf4VdoN9N9Y59thiiwbTzxjXDf"
		datasetCID := "QmR8gF9DpDpGtmKdNDjKk1FbUs3ea9tkQeAahrCcn9Qce1"

		// Fetch dataset and script from IPFS
		pythonScriptData, err := fetchFromIPFS(pythonScriptCID)
		if err != nil {
			log.Fatal(err)
		}
		err = writeToFile("script.py", pythonScriptData)
		if err != nil {
			log.Fatal(err)
		}

		datasetData, err := fetchFromIPFS(datasetCID)
		if err != nil {
			log.Fatal(err)
		}
		err = writeToFile("dataset.csv", datasetData)
		if err != nil {
			log.Fatal(err)
		}

		// Prepare the transaction
		transaction := Transaction{
			DatasetCID:    datasetCID,
			AlgorithmCID:  pythonScriptCID,
			DatasetData:   datasetData,
			AlgorithmData: pythonScriptData,
		}

		// Mine and add the block
		var lastBlock Block
		if len(bc.Blocks) == 0 {
			lastBlock = Block{Index: 0, Header: BlockHeader{Version: 1, PreviousHash: "", Timestamp: time.Now().String()}}
		} else {
			lastBlock = bc.GetLastBlock()
		}

		newBlock := mineBlock(lastBlock, transaction)
		bc.AddBlock(newBlock)

		// Save blockchain and verify block
		if err := saveBlockchain("blockchain.json", bc); err != nil {
			log.Fatal("Error saving blockchain: ", err)
		}

		if verifyBlock(newBlock) {
			log.Println("Block verified successfully!")
		} else {
			log.Println("Block verification failed.")
		}
	}
}
