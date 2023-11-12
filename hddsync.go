package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	err := loadEnv(".env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// from .env
	source := os.Getenv("SOURCE")
	destination := os.Getenv("DESTINATION")
	foldersStr := os.Getenv("FOLDERS")

	folders := strings.Split(foldersStr, ",")

	// Exec rsync
	for _, folder := range folders {
		cmd := exec.Command("rsync", "-avz", "--partial", "--delete", source+folder+"/", destination+folder+"/")
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to sync: %s\n", folder)
			os.Exit(1)
		}
	}

	fmt.Println("Finished!")
}

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		os.Setenv(parts[0], parts[1])
	}
	return scanner.Err()
}
