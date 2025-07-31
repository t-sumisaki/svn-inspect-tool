package diskutil

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type DFEntry struct {
	FileSystem string
	Size       string
	Used       string
	Available  string
	UsePercent string
	MountedOn  string
}

func FindByMount(entries []DFEntry, mountPoint string) *DFEntry {
	for _, entry := range entries {
		if entry.MountedOn == mountPoint {
			return &entry
		}
	}
	return nil
}

func GetDFResult() ([]DFEntry, error) {

	cmd := exec.Command("df", "-h")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, fmt.Errorf("create pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command execute error: %v", err)
	}

	scanner := bufio.NewScanner(stdout)

	var results []DFEntry

	headerSkipped := false

	for scanner.Scan() {

		if !headerSkipped {
			headerSkipped = true
			continue
		}

		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 6 {
			continue
		}

		entry := DFEntry{
			FileSystem: fields[0],
			Size:       fields[1],
			Used:       fields[2],
			Available:  fields[3],
			UsePercent: fields[4],
			MountedOn:  strings.Join(fields[5:], " "),
		}

		results = append(results, entry)
	}

	return results, nil
}

type DUEntry struct {
	Size string
	Path string
}

func GetDUResult(path string) ([]DUEntry, error) {

	cmd := exec.Command("du", "-sh", path)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command execute error: %v", err)
	}

	scanner := bufio.NewScanner(stdout)

	var result []DUEntry

	for scanner.Scan() {

		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) != 2 {
			continue
		}

		entry := DUEntry{
			Size: fields[0],
			Path: fields[1],
		}

		result = append(result, entry)
	}

	return result, nil
}
