package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
)

const jsonPath = "gitview.json"

var args struct {
	Pull    bool `arg:"--pull" help:"Git pull on all repositories"`
	Refresh bool `arg:"--refresh" help:"Create json file which contain the repositories path. This Json can be used to avoid searching phase"`
}

func returnStringFromFile(path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return string(dat)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
	arg.MustParse(&args)
	checkBinExists("git")

	gitRepositories := make(Repository)

	if fileExists(jsonPath) {
		jsonString := returnStringFromFile(jsonPath)

		var jsonDecode []string
		if err := json.Unmarshal([]byte(jsonString), &jsonDecode); err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(jsonDecode); i++ {
			path := jsonDecode[i] + ".git/HEAD"

			fileStat, err := os.Stat(path)
			if err != nil {
				log.Fatal(err)
			}

			writeRepositoryInformation(path, &gitRepositories, fileStat)
		}
	} else {
		err := filepath.Walk("/", func(path string, info os.FileInfo, err error) error {
			// Cannot read the path -> skip
			if err != nil {
				return filepath.SkipDir
			}

			writeRepositoryInformation(path, &gitRepositories, info)

			return nil
		})

		// If error during Walk
		if err != nil {
			log.Println(err)
		}
	}

	// End, write Json file and stop.
	if args.Refresh {
		var repositories []string

		for key := range gitRepositories {
			repositories = append(repositories, key)
		}

		json, err := json.Marshal(repositories)
		if err != nil {
			log.Fatal(err)
		}

		err = ioutil.WriteFile(jsonPath, json, 0644)
		if err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}

	// Start Go routine to check repositories
	var wg sync.WaitGroup

	for key := range gitRepositories {
		wg.Add(1)
		go gitProcess(gitRepositories, key, &wg)
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
