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
	"sync"
	"unicode"

	"github.com/fatih/color"
)

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

// Check if a binary file exists
func checkBinExists(bin string) {
	if _, err := exec.LookPath(bin); err != nil {
		log.Println("didn't find '" + bin + "' executable")
		os.Exit(1)
	}
}

type Git struct {
	branch string
	diff   []string
	status string
}

type Repository map[string]*Git

func main() {
	checkBinExists("git")

	gitRepositories := make(Repository)

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

					var content Git
					content.branch = res[0][1]
					gitRepositories[path] = &content
				}
			}
		}

		return nil
	})

	// If error during Walk
	if err != nil {
		log.Println(err)
	}

	// Start Go routine to check repositories
	var wg sync.WaitGroup

	for key := range gitRepositories {
		wg.Add(1)
		go processGit(gitRepositories, key, &wg)
	}

	wg.Wait()

	// Check column length
	lengthFirstColumn := 0
	lengthSecondColumn := 0
	for key := range gitRepositories {
		if len(key) >= lengthFirstColumn {
			lengthFirstColumn = len(key)
		}

		if len(gitRepositories[key].branch) >= lengthSecondColumn {
			lengthSecondColumn = len(gitRepositories[key].branch)
		}
	}

	for key := range gitRepositories {
		fmt.Print(key)

		// Add space to have column
		if len(key) < lengthFirstColumn {
			space := lengthFirstColumn - len(key)

			for i := 0; i < space; i++ {
				fmt.Print(" ")
			}
		}

		gitBranch := " GIT[" + gitRepositories[key].branch + "]"
		if gitRepositories[key].branch == "master" {
			fmt.Printf("%s", color.GreenString(gitBranch))
		} else {
			fmt.Printf("%s", color.YellowString(gitBranch))
		}

		// Add space to have column
		if len(gitRepositories[key].branch) < lengthSecondColumn {
			space := lengthSecondColumn - len(gitRepositories[key].branch)

			for i := 0; i < space; i++ {
				fmt.Print(" ")
			}
		}

		if gitRepositories[key].status == "outdated" {
			color.Yellow(" " + gitRepositories[key].status)
		} else if gitRepositories[key].status == "up-to-date" {
			color.Green(" " + gitRepositories[key].status)
		} else {
			color.Red(" " + gitRepositories[key].status)
		}

		for i := 0; i < len(gitRepositories[key].diff); i++ {
			fmt.Print("  └─ ")
			color.Cyan(gitRepositories[key].diff[i])
		}
	}
}
