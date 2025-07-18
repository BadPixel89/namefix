package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type replacement struct {
	Match       string `json:"Match"`
	Replacement string `json:"Replacement"`
}

var configDir string

const configName string = "namefix.conf"

func main() {
	filepath, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	replacements := readconfig()
	if replacements == nil {
		os.Exit(1)
	}

	files, err := os.ReadDir(filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	previewRename(files, replacements)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("proceed with changes?:(y/n)")
	text, _ := reader.ReadString('\n')
	//account for newline of windows and unix
	if text == "y\n" || text == "y\r\n" {
		actuallyRename(files, filepath, replacements)
		os.Exit(0)
	}
	fmt.Println("no changes made exiting program")
	os.Exit(0)
}

func previewRename(files []os.DirEntry, replacements []replacement) {
	fmt.Println("current file names")
	for _, file := range files {
		if file.Type().IsDir() {
			continue
		}
		fmt.Println("\t" + file.Name())
	}
	fmt.Println("names after renaming * indicates change")
	for _, file := range files {
		if file.Type().IsDir() {
			continue
		}
		rename := fixName(file.Name(), replacements)
		if rename == file.Name() {
			fmt.Println("\t" + rename)
		} else {
			fmt.Println("*\t" + rename)
		}

	}
}

func actuallyRename(files []os.DirEntry, filepath string, replacements []replacement) {
	for _, file := range files {
		if file.Type().IsDir() {
			continue
		}
		os.Rename(filepath+"/"+file.Name(), filepath+"/"+fixName(file.Name(), replacements))
	}
}
func fixName(filename string, replacers []replacement) string {
	filename = strings.ToLower(filename)
	filenamesplit := strings.Split(filename, ".")
	ext := ""
	if len(filenamesplit) > 1 {
		ext = "." + filenamesplit[len(filenamesplit)-1]
	}
	filename = strings.Replace(filename, ext, "", -1)

	for _, replace := range replacers {
		filename = strings.Replace(filename, replace.Match, replace.Replacement, -1)
	}

	//	keep regex here until we implement into config
	//	matches sequencs of 2 or more '.' '-' or space
	removechars := regexp.MustCompile(`([.\s-]){2,}`)
	remove := removechars.FindAllString(filename, -1)

	for _, substr := range remove {
		//match one to stop replace interfering with itself
		//each substring should occurr at least once if regex matched it
		filename = strings.Replace(filename, substr, "-", 1)
	}
	filename += ext
	//clean up in case last character is a -, -.mp4 looks weird
	filename = strings.Replace(filename, "-.", ".", -1)
	return filename
}
func readconfig() []replacement {
	homepath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("could not find home directory")
		return nil
	}
	configDir = filepath.Join(homepath, ".config")
	_, err = os.Stat(filepath.Join(configDir, configName))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err.Error())
			return nil
		}
		createBlankConfig(configDir)
		return nil
	}

	config, err := os.ReadFile(filepath.Join(configDir, configName))
	if err != nil {
		fmt.Println("failed to read config from: " + configDir)
		fmt.Println(err.Error())
		return nil
	}
	var replacementlist []replacement
	err = json.Unmarshal(config, &replacementlist)
	if err != nil {
		fmt.Println("marshalling fail")
		return nil
	}
	return replacementlist
}
func createBlankConfig(confdir string) {
	err := os.MkdirAll(confdir, 0755)
	if err != nil {
		fmt.Println("unable to make config directory")
		fmt.Println(err.Error())
		return
	}
	_, err = os.Create(filepath.Join(confdir, configName))
	if err != nil {
		fmt.Println("unable to make config file")
		fmt.Println(err.Error())
		return
	}

	var blankconfig []replacement
	dummy := replacement{Match: "", Replacement: ""}
	blankconfig = append(blankconfig, dummy)
	confcontents, err := json.MarshalIndent(blankconfig, "", "\t")
	if err != nil {
		fmt.Println("failed marshalling struct to create blank config")
		fmt.Println(err.Error())
		return
	}
	err = os.WriteFile(filepath.Join(confdir, configName), confcontents, 0755)
	if err != nil {
		fmt.Println("write blank config fail")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("blank config file created: " + filepath.Join(confdir, configName))
	fmt.Println("open with a text editor to add or remove replacemnt definitions")
}
