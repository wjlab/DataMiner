package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"path/filepath"
)

// TNSEntry structure of one entry from tnsnames.ora
type TNSEntry struct {
	Name     string
	Desc     string
	Service  string
	Servers  []TNSAddress
}

// TNSAddress  the host and database port of an address section
type TNSAddress struct {
	Host string
	Port string
}

// TNSEntries Map of tns entries
type TNSEntries map[string]TNSEntry

// Load tnsnames.ora and then parse it, return TNSEntries information
func LoadTNS(tnsFile string) (TNSEntries,error){
	// load available tns entries
	tnsEntries1, err := GetTnsnames(tnsFile)
	length := len(tnsEntries1)
	if err != nil {
		return nil,err
	}
	if length == 0 {
		return nil, errors.New("Got 0 parsed entry from the tnsnames.ora file: "+tnsFile)
	}
	return tnsEntries1,nil
}

// Parse tnsnames.ora file and extract tns entries information
func GetTnsnames(filename string) (TNSEntries,error) {
	var tnsEntries = make(TNSEntries)
	var err error
	var content []string
	var findEntry = regexp.MustCompile(`(?im)^([\w.,\s]+)\s*=(.*)`)
	var tnsAlias = ""
	var desc = ""
	// get the file content
	content, _ = ReadFileByLine(filename)
	// loop through lines
	for _, line := range content {
		if SkipUselessLine(line) {
			continue
		}
		// find a new entry
		newEntry := findEntry.FindStringSubmatch(line)
		i := len(newEntry)
		if i > 0 {
			// save previous entry
			tnsEntries[tnsAlias] = BuildTnsEntry(desc, tnsAlias)
			// new entry
			tnsAlias = newEntry[1]
			if i > 2 {
				desc = newEntry[2] + "\n"
			}
		} else {
			desc += line
		}
	}
	// save last entry
	if len(tnsAlias) > 0 && len(desc) > 0 {
		tnsEntries[tnsAlias] = BuildTnsEntry(desc, tnsAlias)
	}
	return tnsEntries, err
}


// Skip the useless line
func SkipUselessLine(line string) (skip bool) {
	skip = true
	found := false
	reEmpty := regexp.MustCompile(`\S`)
	reComment := regexp.MustCompile(`^#`)
	found = reEmpty.MatchString(line)
	if !found {
		return
	}
	found = reComment.MatchString(line)
	if found {
		return
	}
	skip = false
	return
}

// BuildTnsEntry build map for entry
func BuildTnsEntry(desc string, tnsAlias string) TNSEntry {
	var service = ""
	reService := regexp.MustCompile(`(?mi)(?:SERVICE_NAME|SID)\s*=\s*([\w.]+)`)
	s := reService.FindStringSubmatch(desc)
	if len(s) > 1 {
		service = s[1]
	}
	servers := GetServers(desc)
	entry := TNSEntry{Name: tnsAlias, Desc: desc, Service: service, Servers: servers}
	return entry
}

// GetServers extract TNSAddress part
func GetServers(tnsDesc string) (servers []TNSAddress) {
	re := regexp.MustCompile(`(?m)HOST\s*=\s*([\w\-_.]+)\s*\)\s*\(\s*PORT\s*=\s*(\d+)`)
	match := re.FindAllStringSubmatch(tnsDesc, -1)
	for _, a := range match {
		if len(a) > 1 {
			host := a[1]
			port := a[2]
			servers = append(servers, TNSAddress{
				Host: host, Port: port,
			})
		}
	}
	return
}

// ReadFileByLine read a file and return string array of lines
func ReadFileByLine(filename string) ([]string, error) {
	var lines []string
	filename = filepath.Clean(filename)
	if _, err := os.Stat(filename); err != nil {
		return lines, fmt.Errorf("file %s not found", filename)
	}
	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return lines, err
	}
	var line string
	reader := bufio.NewReader(f)
	for {
		line, err = reader.ReadString('\n')
		lines = append(lines, line)
		if err != nil {
			break
		}
	}
	if err == io.EOF {
		err = nil
	}
	return lines, err
}