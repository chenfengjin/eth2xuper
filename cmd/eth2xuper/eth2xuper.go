package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var output string
	flag.StringVar(&output, "output", "eth2xuper", "output dir")
	flag.Parse()
	args := flag.Args()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	for i := 0; i < flag.NArg(); i++ {
		if !strings.HasSuffix(args[i], "json") {
			log.Printf("\t\t %s is not a valid OpenZeppelin output", args[i])
			continue
		}
		inputFileName := args[i]
		if err := os.MkdirAll(output, 0755); err != nil {
			log.Fatal(err)
		}

		f, err := os.Open(inputFileName)
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadAll(f);
		if err != nil {
			log.Fatal(err)
		}

		code := struct {
			ABI     []map[string]interface{} `json:"abi"`
			EvmCode string                   `json:"bytecode"`
		}{}

		err = json.Unmarshal(data, &code)
		if err != nil {
			log.Fatal(err)
		}

		code.EvmCode = strings.TrimLeft(code.EvmCode, "0x") // Remove 0x prefix

		if err != nil {
			log.Fatal(err)
		}
		d, err := json.Marshal(code.ABI)

		_, solFilename := filepath.Split(inputFileName)
		solFilename = strings.TrimSuffix(solFilename,filepath.Ext(solFilename))
		abiFileName := solFilename + ".abi"
		binFileName := solFilename + ".bin"
		if err := ioutil.WriteFile(filepath.Join(output,abiFileName), d, 0644); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(filepath.Join(output,binFileName), []byte(code.EvmCode), 0644); err != nil {
			log.Fatal("err")
		}
		//输入可能带目录
		log.Printf("convert succeed to directory %s\n", output)
	}
}
