package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/fatih/color"
	"unicode"
)

// remove HEAD/.git from the path
func getRootGitFolderFromHeadFile(path string) string {
	return string(path[0:len(path)-9])
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

// Check if a binary file exists
func checkBinExists(bin string) {
	if _, err := exec.LookPath(bin); err != nil {
		log.Println("didn't find '" + bin + "' executable")
		os.Exit(1)
	}
}

func main() {
	checkBinExists("git")
	// gitRepository["path"] = "branch"
	gitRepository := make(map[string]string)

	err := filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
		// Cannot read the path -> skip
		if err != nil {
			return filepath.SkipDir
		}

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
					gitRepository[path] = res[0][1]
				}
			}
		}

		return nil
	})

	// If error during Walk
	if err != nil {
		log.Println(err)
	}

	// git --git-dir=/path/.git --work-tree=/path/ status --porcelain
	for key, value := range gitRepository {
		cmdName := "git"
		cmdArgs := []string{"--git-dir=" + key + ".git", "--work-tree=" + key, "status", "--porcelain"}

		// Execute git command
		out, err := exec.Command(cmdName, cmdArgs...).Output();

		if err != nil {
			fmt.Fprintln(os.Stderr, "There was an error running git status command: ", err)
			os.Exit(1)
		}

		// Display git path + current branch
		fmt.Print(key)
		if value == "master" {
			color.Green(" GIT[" + value + "]")
		} else {
			color.Yellow(" GIT[" + value + "]")
		}

		// Unicode
		sha := string(out)

		result := strings.Split(sha, "\n")
		for i := 0; i < len(result); i++ {
			if match, _ := regexp.MatchString("^\\s*[MADRCU]", result[i]); match {
				fmt.Print("  └─ ")
				color.Blue(sanitizeGitStatus(result[i]))
			}
		}
	}
}
