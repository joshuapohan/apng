package main

import (
	"fmt"
	"os"
	"image/png"
	"bytes"
	"io/ioutil"
	"log"
)

func logError(e error){
	if e != nil{
		log.Fatal(e)
	}
}

func readInputFiles(ap *apng){
	files, err := ioutil.ReadDir("./input")
	logError(err)

	for _, fileInfo := range files{
		f, err := os.Open("./input/" + fileInfo.Name())
		logError(err)
		curPng, err := png.Decode(f)
		logError(err)

		curImgBuffer := new(bytes.Buffer)

		if err := png.Encode(curImgBuffer, curPng); err != nil{
			fmt.Println(err)
			return
		}

		
		curPngChunk := getPNGChunk(curImgBuffer)
		(*ap).chunks = append((*ap).chunks, curPngChunk)
	}
}

func ProcessInputFiles(){
	var pngs apng
	readInputFiles(&pngs)
	pngs.LogPNGChunks()
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