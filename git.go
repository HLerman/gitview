package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

func gitProcess(gitRepositories Repository, key string, wg *sync.WaitGroup) {
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

	if args.Pull {
		if content.status == "outdated" && len(content.diff) == 0 && content.branch == "master" {
			// Git pull
			cmdArgs = []string{"--git-dir=" + key + ".git", "--work-tree=" + key, "pull"}

			// Execute git command
			_, err = exec.Command(cmdName, cmdArgs...).Output()

			if err != nil {
				content.status += ", failed to update"
				gitRepositories[key] = &content
				return
			}

			content.status = "up-to-date"
		}
	}

	gitRepositories[key] = &content
}

// remove HEAD/.git from the path
func getRootGitFolderFromHeadFile(path string) string {
	return string(path[0 : len(path)-9])
}

func sanitizeGitStatus(data string) string {
	data = strings.TrimSpace(data)
	newData := data

	// Check if we have more than 1 space between [MADRCU] and the file/folder
	// For example : A  test
	if match, _ := regexp.MatchString("^[MADRCU]\\s{2,}", data); match {
		// Loop all characters execept the two frist bytes (beacause it's [MADRCU] + space and we
		// would like to sav 1 space
		for i := 2; i < len(data); i++ {
			// Check if the current iteration is a space
			if unicode.IsSpace(int32(data[i])) {
				newData = newData[0:i] + newData[i+1:len(data)]
			} else {
				return newData
			}
		}
	}

	return data
}

func writeRepositoryInformation(path string, gitRepositories *Repository, info os.FileInfo) {
	// Check if the path is a HEAD git file, if yes we can open it
	if match, _ := regexp.MatchString("\\.git\\/HEAD$", path); match && !info.IsDir() {
		f, err := os.Open(path)

		// Cannot open the file
		if err != nil {
			log.Println(err)
		}

		// Try to determine the branch
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {

			// If the content file is correctly formed
			if match, _ := regexp.MatchString("^ref: refs\\/heads\\/.+$", scanner.Text()); match {
				// Get the branch
				match, _ := regexp.Compile("^ref: refs\\/heads\\/(.+)$")
				res := match.FindAllStringSubmatch(scanner.Text(), -1)

				// Rewrite path
				path = getRootGitFolderFromHeadFile(path)
				// Add the branch into gitRepository[path]

				var content Git
				content.branch = res[0][1]
				(*gitRepositories)[path] = &content
			}
		}
	}
}
