package global

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunTerraformCommand runs a terraform command in a folder
func RunTerraformCommand(folder string, command []string) error {
	cmd := exec.Command("terraform", command...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = folder
	return cmd.Run()
}

func GetTerraformOutputs(folder string, outputs map[string]string) (err error) {
	cmd := exec.Command("terraform", "output")
	cmd.Dir = folder
	outputBytes, err := cmd.Output()
	if err != nil {
		return
	}
	text := string(outputBytes)
	scanner := bufio.NewScanner(strings.NewReader(string(text)))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		txt := scanner.Text()
		txt = strings.TrimSpace(txt)
		parts := strings.Split(txt, "=")
		if len(parts) != 2 {
			continue
		}
		outputs[strings.TrimSpace(parts[0])] = TrimAndRemoveQuotes(parts[1])
	}
	return
}

func TerraformDeploy(buildFolder, craftName string) (err error) {
	fmt.Printf("Deploying '%s' from '%s'\n", craftName, buildFolder)
	err = RunTerraformCommand(buildFolder, []string{"fmt", "-list=false"})
	if err != nil {
		return
	}

	// Now deploy
	err = RunTerraformCommand(buildFolder, []string{"init", "-migrate-state"})
	if err != nil {
		return
	}

	err = RunTerraformCommand(buildFolder, []string{"apply", "--auto-approve"})
	if err != nil {
		return
	}
	fmt.Printf("Done deploying '%s'\n", craftName)
	return
}
