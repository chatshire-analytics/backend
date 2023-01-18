package main

import (
	"fmt"
	"os"
	"text/template"
)

type input struct {
	fileName string
	values   map[string]interface{}
}

func main() {
	inputs := make([]input, 0)
	mediumValues := map[string]interface{}{
		"Dependencies": []interface{}{
			map[string]string{
				"groupId":    "org.postgresql",
				"artifactId": "postgresql",
				"version":    "42.5.0",
			},
		},
	}
	githubPages := map[string]interface{}{
		"endpoints": []map[string]string{
			{
				"method": "GET",
				"path":   "/",
				"status": "http.StatusOK",
				"body":   "Hello, World!",
			},
			{
				"method": "POST",
				"path":   "/payment",
				"status": "http.StatusOK",
				"body":   "Payment is created",
			},
		},
	}
	mediumWebServerValues := map[string]interface{}{
		"httpLibrary": "gin",
	}
	iterateMap := map[string]interface{}{
		"endpoints": map[string]string{
			"/":        "Hello, World!",
			"/payment": "Payment is created",
		},
	}
	nodejsValues := map[string]interface{}{
		"Dependencies": []interface{}{
			map[string]interface{}{
				"name":          "jest",
				"version":       "29.3.1",
				"devDependency": true,
			},
		},
	}
	inputs = append(inputs,
		input{fileName: "documents/medium-pom.xml", values: mediumValues},
		input{fileName: "documents/github-pages.txt", values: githubPages},
		input{fileName: "documents/x.txt", values: mediumWebServerValues},
		input{fileName: "documents/map.txt", values: iterateMap},
		input{fileName: "documents/package.json", values: nodejsValues},
	)
	//getwd, _ := os.Getwd()
	//fmt.Println(getwd)
	for _, i := range inputs {
		file, err2 := os.ReadFile(i.fileName)
		if err2 != nil {
			return
		}
		tmpl, err := template.New("test").Parse(string(file))
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(os.Stdout, i.values)
		fmt.Println()
		if err != nil {
			panic(err)
		}
	}
}
