package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	shell "github.com/ipfs/go-ipfs-api"
)

// FetchFromIPFS fetches data from IPFS given a CID
func fetchFromIPFS(cid string) ([]byte, error) {
	sh := shell.NewShell("http://localhost:5001")

	data, err := sh.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from IPFS: %w", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)

	return buf.Bytes(), nil
}

// WriteToFile writes the fetched data to a file
func writeToFile(filePath string, data []byte) error {
	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error saving data to file: %w", err)
	}
	return nil
}
