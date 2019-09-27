package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"log"
	"apng"
)

func logError(e error){
	if e != nil{
		log.Fatal(e)
	}
}

func readInputFiles(){
	files, err := ioutil.ReadDir("../input")

	logError(err)
	test := &apng.APNGModel{}

	for _, fileInfo := range files{
		f, err := os.Open("../input/" + fileInfo.Name())
		logError(err)
		test.AppendImage(f)
		test.AppendDelay(64)
		logError(err)
	}
	test.Encode()
	test.SavePNGData("logAPNG9.png")
}

func ProcessInputFiles(){
	fmt.Println("Writing input files")
	readInputFiles()
}

func main(){
	ProcessInputFiles()
}
