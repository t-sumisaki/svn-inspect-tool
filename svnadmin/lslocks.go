package svnadmin

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

type SvnLockInfo struct {
	Path    string
	Owner   string
	Created string
	Comment string
}

const (
	svnlock_Prefix_Path    = "Path:"
	svnlock_Prefix_Owner   = "Owner:"
	svnlock_Prefix_Created = "Created:"
	svnlock_Prefix_Comment = "Comment"
)

func GetLockInfo(repoPath string) ([]SvnLockInfo, error) {

	cmd := exec.Command("svnadmin", "lslocks", repoPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("create pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("command execute error: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	var (
		lockInfos []SvnLockInfo
		current   SvnLockInfo
	)

	bCaptureComments := false

	for scanner.Scan() {

		line := scanner.Text()

		switch {

		case strings.HasPrefix(line, svnlock_Prefix_Path):
			if current.Path != "" {
				lockInfos = append(lockInfos, current)
				current = SvnLockInfo{}
			}
			current.Path = strings.TrimSpace(strings.TrimPrefix(line, svnlock_Prefix_Path))

		case strings.HasPrefix(line, svnlock_Prefix_Owner):
			current.Owner = strings.TrimSpace(strings.TrimPrefix(line, svnlock_Prefix_Owner))

		case strings.HasPrefix(line, svnlock_Prefix_Created):
			current.Created = strings.TrimSpace(strings.TrimPrefix(line, svnlock_Prefix_Created))

		case strings.HasPrefix(line, svnlock_Prefix_Comment):
			bCaptureComments = true
		case bCaptureComments:
			current.Comment = line
			bCaptureComments = false
		}
	}

	if current.Path != "" {
		lockInfos = append(lockInfos, current)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command finish error: %v", err)
	}

	return lockInfos, nil
}
