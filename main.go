package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

//===main====
func main() {

	http.HandleFunc("/", welcome)      // welcome page
	http.HandleFunc("/login", login)   //input data page
	http.HandleFunc("/result", result) //output page

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

//===welcome page====

func welcome(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello USER!") // write data to response
}

//===login page==============

func login(w http.ResponseWriter, r *http.Request) {
	//	fmt.Println("method:", r.Method) //get request method

	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		fmt.Println("username:", r.Form["username"])

		getdata := r.Form["username"]

		//===============process data input========

		str, err := json.Marshal(getdata)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}

		userinput := string(str[:])

		re, err := regexp.Compile(`[^\w]`)
		if err != nil {
			log.Fatal(err)
		}
		userinput = re.ReplaceAllString(userinput, " ")
		userinput = strings.ToLower(userinput)
		// fmt.Println(userinput)
		//====================count word and append================
		for word, occur := range countsimilarword(userinput) {
			occurance := strconv.Itoa(occur)
			var userinput []string

			userinput = append(userinput, "  word: ", word, "  occurs: ", occurance, "  times  ")

			file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

			if err != nil {
				log.Fatalf("failed creating file: %s", err)
			}

			datawriter := bufio.NewWriter(file)

			for _, data := range userinput {
				_, _ = datawriter.WriteString(data)
			}

			datawriter.Flush()
			file.Close()
		}

	}

}

func result(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "result.html")
}

// ======================count similar word===========
func countsimilarword(st string) map[string]int {

	input := strings.Fields(st)
	wordcount := make(map[string]int)
	for _, word := range input {
		_, matched := wordcount[word]
		if matched {
			wordcount[word] += 1
		} else {
			wordcount[word] = 1
		}
	}
	return wordcount

}

