package diskutil

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
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

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finish error: %v", err)
	}

	return results, nil
}

type DUEntry struct {
	Size string
	Path string
}

func GetDUResult(path string) ([]DUEntry, error) {

	paths, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("glob error: %v", err)
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("no matching paths found")
	}

	args := append([]string{"-sh"}, paths...)

	cmd := exec.Command("du", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create pipe error: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("create stderr pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command execute error: %v", err)
	}

	var errBuf bytes.Buffer
	go func() {
		_, _ = errBuf.ReadFrom(stderr)
	}()

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

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finish error: %v, stderr: %s", err, strings.TrimSpace(errBuf.String()))
	}

	return result, nil
}
