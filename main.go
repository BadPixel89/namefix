package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	filepath, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	files, err := os.ReadDir(filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	previewRename(files)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("proceed with changes?:(y/n)")
	text, _ := reader.ReadString('\n')

	if text != "y\n" {
		fmt.Println("no changes made. exiting program")
		os.Exit(0)
	}
	actuallyRename(files, filepath)
}

func previewRename(files []os.DirEntry) {
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
		rename := fixName(file.Name())
		if rename == file.Name() {
			fmt.Println("\t" + rename)
		} else {
			fmt.Println("*\t" + rename)
		}

	}
}

func actuallyRename(files []os.DirEntry, filepath string) {
	for _, file := range files {
		os.Rename(filepath+"/"+file.Name(), filepath+"/"+fixName(file.Name()))
	}
}

func fixName(filename string) string {
	filename = strings.ToLower(filename)

	filenamesplit := strings.Split(filename, ".")
	ext := ""
	if len(filenamesplit) > 1 {
		ext = "." + filenamesplit[len(filenamesplit)-1]
	}
	filename = strings.Replace(filename, ext, "", -1)

	filename = strings.Replace(filename, "1080p", "", -1)
	filename = strings.Replace(filename, "2160p", "", -1)
	filename = strings.Replace(filename, "bluray", "", -1)
	filename = strings.Replace(filename, "webrip", "", -1)
	filename = strings.Replace(filename, "x264", "", -1)
	filename = strings.Replace(filename, "x265", "", -1)
	filename = strings.Replace(filename, "yify", "", -1)
	filename = strings.Replace(filename, "rarbg", "", -1)
	filename = strings.Replace(filename, "hevc", "", -1)
	filename = strings.Replace(filename, "brrip", "", -1)
	filename = strings.Replace(filename, "dd5.1", "", -1)
	filename = strings.Replace(filename, "5.1", "", -1)
	filename = strings.Replace(filename, " ", ".", -1)

	removeDash := regexp.MustCompile(`([.\s-]){2,}`)

	remove := removeDash.FindString(filename)

	if remove != "" {
		filename = strings.Replace(filename, remove, "-", -1)
	}

	filename = strings.Replace(filename, "--", "-", -1)
	filename = strings.Replace(filename, ".", "-", -1)
	filename += ext
	filename = strings.Replace(filename, "-.", ".", -1)
	return filename
}
