package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var ConfigFile = flag.String("config", "./config.ini", "your config file")
var config *Config

func setpbcopy(content []byte) error {
	_, fErr := os.Stat(config.CPFile)
	if fErr != nil {
		if os.IsNotExist(fErr) {
			_, err := os.Create(config.CPFile)
			if err != nil {
				return err
			}
			return nil
		}

		return fErr
	}

	err := ioutil.WriteFile(config.CPFile, content, 0644)
	if err != nil {
		fmt.Println(1)
		return err
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf(`cat %s | pbcopy`, config.CPFile))

	if _, err = cmd.Output(); err != nil {
		return err
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	//buffer := bytes.NewBuffer(make([]byte, 2 * 1024 * 1024))
	var builder strings.Builder
	_, err := io.Copy(&builder, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		return
	}

	setErr := setpbcopy([]byte(builder.String()))
	if setErr != nil {
		log.Println(setErr)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server Error"))
		return
	}

	fmt.Fprintf(w, "Done!")
}

func main() {
	flag.Parse()

	var errStr string
	config, errStr = initConfig(*ConfigFile)
	if errStr != "" {
		log.Fatal(errStr)
	}

	// start Server
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatal("spbcopy server: ", err)
	}
}
