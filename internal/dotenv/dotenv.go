package dotenv

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func Load() {

	var mode string
	flag.StringVar(&mode, "mode", "development", "Define the go environment mode")
	flag.Parse()
	os.Setenv("mode", mode)

	f, err := os.Open(fmt.Sprintf("./%s.env", mode))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		kv := scanner.Text()
		separator := strings.Index(kv, "=")
		key := kv[:separator]
		value := kv[separator+1:]
		os.Setenv(key, value)
	}

}
