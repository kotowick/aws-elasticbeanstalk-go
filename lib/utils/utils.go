package utils

import (
    "fmt"
    "io/ioutil"
    "github.com/smallfish/simpleyaml"
  	"os/exec"
  	"log"
  	"encoding/json"
)

type Interface interface {}

type recmap map[string]interface{}

type ShellOutParams struct {
  CmdName string
	CmdArgs []string
	Key string
	Keys []string
}

/**
  Method: Print a separator string
*/
func PrintSeparator(){
  fmt.Printf("\n--------------------------\n")
}

/**
  Method: Print a separator string with a message in-between
*/
func PrintSeparatorWithMessage(message string){
  PrintSeparator()
  fmt.Printf("\n%s\n", message)
  PrintSeparator()
}

/**
  Method: Get the default value between two strings. The first param (a) will
          always be the prefered option if set.
*/
func GetDefault(a, b string) string {
  if a != "" {
      return a
  } else if b != "" {
    return b
  } else {
    return ""
  }
}

func GetConfig(file_path string) map[string]string {
  source, err := ioutil.ReadFile(file_path)
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

func CombineConfigOptions(configOptions map[string]string, additionalConfigOptions map[string]string) map[string]string{
  for k, v := range additionalConfigOptions {
    configOptions[k] = v
  }
  return configOptions
}

/**
	Method: Iterate through an interface of keys. Print out key names that is
					specified by the key_names input.
*/
func Keys(objects []interface{}, key_names []string, ) {
	fmt.Println("\n")
  fmt.Println("Before J")
	for j := range objects{
    fmt.Println("Inside J")
		for key := range key_names{
			object := objects[j].(map[string]interface{})
			fmt.Printf("|_____ %q\n", object[key_names[key]])
		}
	}
}

/**
	Method: execute a shell command and output specified key values
*/
func ShellOut(message string, p ShellOutParams){
	if message != ""{
    PrintSeparatorWithMessage(message)
  }

	stdout, err := exec.Command(p.CmdName, p.CmdArgs...).Output()

	if err != nil {
		fmt.Printf("%s\n", p.CmdArgs)
		fmt.Printf("%s\n", stdout)
    log.Fatal(err)
	}

	if p.Key != ""{
		m := make(recmap)
		err = json.Unmarshal(stdout, &m)

	  objects := m[p.Key].([]interface{})
		Keys(objects, p.Keys)
	}
}
