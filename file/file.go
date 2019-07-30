package file

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
)

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

func checkBinExists(bin string) {
	if _, err := exec.LookPath(bin); err != nil {
		log.Println("didn't find '" + bin + "' executable")
		os.Exit(1)
	}
}

func getJSONPath() (string, error) {
	usr, err := user.Current()

	if err != nil {
		return "", err
	}

	return usr.HomeDir + "/gitview.json", nil
}

func createJSONFile(repositories []string) error {
	json, err := json.Marshal(repositories)

	if err != nil {
		return err
	}

	jsonPath, err := getJSONPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(jsonPath, json, 0644)
	if err != nil {
		return err
	}

	return nil
}
