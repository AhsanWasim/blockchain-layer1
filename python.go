package main

import (
	"fmt"
	"os/exec"
)

// RunPythonScript runs a Python script with the given dataset and returns the output
func runPythonScript(scriptPath, datasetPath string) ([]byte, error) {
	cmd := exec.Command("python", scriptPath, datasetPath)

	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running python script: %w\nOutput: %s", err, string(cmdOutput))
	}

	return cmdOutput, nil
}
