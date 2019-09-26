package main

import (
	"fmt"
	"os"
	"image/png"
	"bytes"
	"io/ioutil"
	"log"
	"apng"
)

func logError(e error){
	if e != nil{
		log.Fatal(e)
	}
}

func readInputFiles(ap *apng.APNGModel){
	files, err := ioutil.ReadDir("./input")

	logError(err)
	test := &apng.APNGModel{}

	for _, fileInfo := range files{
		f, err := os.Open("./input/" + fileInfo.Name())
		logError(err)
		curPng, err := png.Decode(f)
		test.AppendImage(f)

		logError(err)

		curImgBuffer := new(bytes.Buffer)

		if err := png.Encode(curImgBuffer, curPng); err != nil{
			fmt.Println(err)
			return
		}

		

		
		//curPngChunk := apng.GetPNGChunk(curImgBuffer)
		//(*ap).chunks = append((*ap).chunks, curPngChunk)
	}
}

func readInputFiles2(){
	files, err := ioutil.ReadDir("../input")

	logError(err)
	test := &apng.APNGModel{}

	for _, fileInfo := range files{
		f, err := os.Open("../input/" + fileInfo.Name())
		logError(err)
		test.AppendImage(f)
		logError(err)
	}
	test.Encode()
	test.SavePNGData()
}

func ProcessInputFiles(){
	//var pngs apng.APNGModel
	//readInputFiles(&pngs)
	//pngs.LogPNGChunks()
	
	fmt.Println("Writing input files")
	readInputFiles2()
}

func main(){
	/*f, _ := os.Open("./mario.png")
	inPng, err := png.Decode(f)
	if err != nil{
		fmt.Println(err)
		f.Close()
		return
	}

	imgBuffer := new(bytes.Buffer)

	if err := png.Encode(imgBuffer, inPng); err != nil{
		fmt.Println(err)
		return
	}

	chunk := getPNGChunk(imgBuffer)

	savePNGChunk(chunk)*/
	ProcessInputFiles()
}
