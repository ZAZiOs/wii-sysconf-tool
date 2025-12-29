package cmd

import (
	"fmt"
	"os"

	"sysconf-parser/sysconf"
)

func Encode(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	sys, err := sysconf.FromJSON(data)
	if err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %w", filename, err)
	}

	if sys == nil || len(sys.Items) == 0 {
		return fmt.Errorf("sysconf JSON is empty or invalid")
	}

	out, err := sysconf.Write(sys)
	if err != nil {
		return fmt.Errorf("failed to write sysconf binary: %w", err)
	}

	outFile := filename + ".bin"
	if err := os.WriteFile(outFile, out, 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outFile, err)
	}

	fmt.Printf("Encoded %s -> %s\n", filename, outFile)
	return nil
}

func Decode(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) != 0x4000 {
		return fmt.Errorf("invalid SYSCONF size: got %d, need 0x4000", len(data))
	}

	sys, err := sysconf.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse SYSCONF: %w", err)
	}

	jsonData, err := sysconf.ToJSON(sys)
	if err != nil {
		return fmt.Errorf("failed to convert to JSON: %w", err)
	}

	jsonFile := filename + ".json"
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	fmt.Printf("Decoded %s -> %s\n", filename, jsonFile)
	return nil
}

func Usage() {
	fmt.Println("Usage:")
	fmt.Println("  sysconfer.exe decode SYSCONF")
	fmt.Println("  sysconfer.exe encode sysconf.json")
}
