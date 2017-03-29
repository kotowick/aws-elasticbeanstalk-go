// Package utils provides a client for AWS S3.
//
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/smallfish/simpleyaml"
)

// recmap provides a struct of string: interfaces{}
// It is used for Unmarshalling
type recmap map[string]interface{}

// ShellOutParams provides a struct contains:
// 	- command name
//	- command args
//	- the first key in the response to iterate over
//	- the keys to map
//
type ShellOutParams struct {
	CmdName string
	CmdArgs []string
	Key     string
	Keys    []string
}

// PrintSeparator displays a separator line
//
func PrintSeparator() {
	fmt.Printf("\n--------------------------\n")
}

// PrintSeparatorWithMessage prints a separator with the specified message
// in-between
//
func PrintSeparatorWithMessage(message string) {
	PrintSeparator()
	fmt.Printf("\n%s\n", message)
	PrintSeparator()
}

// GetDefault returns a default value when two strings are passed in
// The first param (a) will always be the prefered option if set.
//
func GetDefault(a, b string) string {
	if a != "" {
		return a
	} else if b != "" {
		return b
	} else {
		return ""
	}
}

// GetConfig parses a YML file and returns a key:value map
//
func GetConfig(filePath string) map[string]string {
	source, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	yaml, err := simpleyaml.NewYaml(source)
	if err != nil {
		panic(err)
	}

	m, err := yaml.Map()

	if err != nil {
		fmt.Println(err)
	}

	values := make(map[string]string)

	for k, v := range m {
		if s, ok := k.(string); ok {
			values[s] = v.(string)
		}
	}

	return values
}

// CombineConfigOptions combines a given key:value map with another key:value map
func CombineConfigOptions(configOptions map[string]string, additionalConfigOptions map[string]string) map[string]string {
	for k, v := range additionalConfigOptions {
		configOptions[k] = v
	}
	return configOptions
}

// Keys prints out all keys in a string array that are within the []interface{}
//
func Keys(objects []interface{}, keyNames []string) {
	fmt.Println("")
	fmt.Println("Before J")
	for j := range objects {
		fmt.Println("Inside J")
		for key := range keyNames {
			object := objects[j].(map[string]interface{})
			fmt.Printf("|_____ %q\n", object[keyNames[key]])
		}
	}
}

// ShellOut executes a given command and displays specified key values
//
func ShellOut(message string, p ShellOutParams) {
	if message != "" {
		PrintSeparatorWithMessage(message)
	}

	stdout, err := exec.Command(p.CmdName, p.CmdArgs...).Output()

	if err != nil {
		fmt.Printf("%s\n", p.CmdArgs)
		fmt.Printf("%s\n", stdout)
		log.Fatal(err)
	}

	if p.Key != "" {
		m := make(recmap)
		if err := json.Unmarshal(stdout, &m); err != nil {
			log.Fatal(err)
		}

		objects := m[p.Key].([]interface{})
		Keys(objects, p.Keys)
	}
}

// VerifyParamatersWithOr returns true if any of the paramaters in args[] are empty.
// This checks to see if required variables are set or not
func VerifyParamatersWithOr(args map[string]string) bool {
	for k, v := range args {
		if v == "" {
			log.Fatal(fmt.Sprintf("%s is not set.", k))
			return false
		}
	}
	return true
}

// VerifyParamatersWithAnd returns true if any of the paramaters in args[] are not empty.
// This checks to see if required variables are set or not
func VerifyParamatersWithAnd(args map[string]string) bool {
	for _, v := range args {
		if v != "" {
			return true
		}
	}
	log.Fatal("VerifyParamatersWithOr failed: all of the paramaters returns empty")
	return false
}
