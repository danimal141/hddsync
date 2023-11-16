package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Load the .env file
	err := loadEnv(".env")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get paths from environment variables
	source := os.Getenv("SOURCE")
	destination := os.Getenv("DESTINATION")
	foldersStr := os.Getenv("FOLDERS")

	// Split FOLDERS environment variable to create a list of folders
	folders := strings.Split(foldersStr, ",")

	// Confirm HDD and folders to sync
	fmt.Printf("Source: %s\n", source)
	fmt.Printf("Destination: %s\n", destination)
	fmt.Printf("Folders to sync: %v\n", folders)
	fmt.Print("Do you want to start sync with these settings? (y/n): ")

	// Wait for user input
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	if strings.TrimSpace(response) != "y" {
		fmt.Println("Sync cancelled.")
		os.Exit(0)
	}

	// logging
	logFile, err := os.OpenFile("rsync_errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	// Execute rsync for each folder
	for _, folder := range folders {
		cmd := exec.Command("rsync", "-avzPv", "--partial", "--delete", "--exclude", ".DS_Store", source+folder+"/", destination+folder+"/")
		cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		cmd.Stderr = logFile
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error occurred during sync: %s\n", folder)
			// do not stop
			// os.Exit(1)
		}
	}

	fmt.Println("Sync completed")
}

// Function to load .env file
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
