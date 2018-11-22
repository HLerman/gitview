package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

func processGit(gitRepositories Repository, key string, wg *sync.WaitGroup) {
	defer wg.Done()

	var content Git
	content.branch = gitRepositories[key].branch
	// Git status
	cmdName := "git"
	cmdArgs := []string{"--git-dir=" + key + ".git", "--work-tree=" + key, "status", "--porcelain"}

	// Execute git command
	out, err := exec.Command(cmdName, cmdArgs...).Output()

	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running git command: ", err)
		os.Exit(1)
	}

	// Unicode
	sha := string(out)

	// Git remote show origin
	cmdArgs = []string{"--git-dir=" + key + ".git", "--work-tree=" + key, "remote", "show", "origin"}

	// Execute git command
	out, err = exec.Command(cmdName, cmdArgs...).Output()

	if err != nil {
		content.status = "connection failed"
		gitRepositories[key] = &content
		return
	}

	content.status = "outdated"

	status := string(out)

	result := strings.Split(status, "\n")
	for i := 0; i < len(result); i++ {
		if match, _ := regexp.MatchString("^\\s*.+\\(up to date\\)", result[i]); match {
			content.status = "up-to-date"
		}
	}

	result = strings.Split(sha, "\n")
	for i := 0; i < len(result); i++ {
		if match, _ := regexp.MatchString("^\\s*[MADRCU]", result[i]); match {
			content.diff = append(content.diff, sanitizeGitStatus(result[i]))
		}
	}

	gitRepositories[key] = &content
}
