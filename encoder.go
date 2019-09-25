package apng

import (
	"fmt"
	"os"
	"bytes"
	"io"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/png"
)


/******************************************************
					apng structure
*******************************************************/

type pngData struct{
	ihdr []byte
	idat []byte
}

type APNGModel struct{
	images []image.Image
	chunks []pngData
	delays []int
}

func (ap APNGModel) PrintPNGChunks(){
	for _, png := range ap.chunks  {
		fmt.Println("IHDR")
		fmt.Println(png.ihdr)
		fmt.Println("IDAT")
		fmt.Println(png.idat)
	}	
}

func (ap APNGModel) LogPNGChunks(){

	f, _ := os.Create("log.txt")

	for _, png := range ap.chunks {
		f.Write([]byte("IHDR\n"))
		f.Write(png.ihdr)
		f.Write([]byte("\n"))
		f.Write([]byte("IDAT\n"))
		f.Write(png.idat)
		f.Write([]byte("\n"))
	}	
}

func (ap *APNGModel) AppendImage(r io.Reader){
	curPng, err := png.Decode(r)
	if err != nil{
		fmt.Println(err)
		return
	}
	ap.images = append(ap.images, curPng)
}

func (ap *APNGModel) AppendDelay(delay int){
	ap.delays = append(ap.delays, delay)
}

/******************************************************
				png chunk manipulation
*******************************************************/

func GetPNGChunk(imgBuffer *bytes.Buffer) (pngData){
	chunk := pngData{}

	//skip png header
	imgBuffer.Next(8)

	for {
		tmp := make([]byte, 8)
		_, err := io.ReadFull(imgBuffer, tmp[:8])
		if err != nil{
			fmt.Println("Error : ", err)
			break
		}

		length := binary.BigEndian.Uint32(tmp[:4])
		

		tmpVal := make([]byte, length)
		io.ReadFull(imgBuffer, tmpVal)

		switch string(tmp[4:8]){
		case "IHDR":
			chunk.ihdr = make([]byte, length)
			copy(chunk.ihdr, tmpVal)
		case "IDAT":
			chunk.idat = append(chunk.idat, tmpVal...)
			tmpVal = nil
		default:
			fmt.Println("Found ", string(tmp[4:8]))
		}

		//skip crc
		imgBuffer.Next(4)
	}

	return chunk
}

func writeUint32(b []uint8, u uint32) {
	b[0] = uint8(u >> 24)
	b[1] = uint8(u >> 16)
	b[2] = uint8(u >> 8)
	b[3] = uint8(u)
}

func writePNGHeader(b *[]uint8){
	*b = append(*b, 0x89,0x50,0x4E,0x47,0x0D,0x0A,0x1A,0x0A)
}

func writeIENDHeader(b *[]uint8){
	*b = append(*b,[]byte("IEND")...)
}

func writeCRC32(data *[]byte){
		crcBytes := make([]byte, 4)
		crc := crc32.NewIEEE()
		crc.Write(*data)
		writeUint32(crcBytes, crc.Sum32())
		*data = append(*data, crcBytes...)
}

func writeChunk(chunk []byte, header string, toChunk *[]byte){
	chunkLe := make([]byte, 4)
	chunkTagVal := make([]byte,0, len(chunk) + 8)

	writeUint32(chunkLe, uint32(len(chunk)))

	chunkTagVal = append(chunkTagVal, []byte(header)...)
	chunkTagVal = append(chunkTagVal, chunk...)

	writeCRC32(&chunkTagVal)
	*toChunk = append(*toChunk, chunkLe...)
	*toChunk = append(*toChunk, chunkTagVal...)
}

func writeacTL(b *[]byte, numFrames int, numPlays int){
	*b = append(*b, []byte("acTL")...)
}

func writefcTL(b *[]byte){
	*b = append(*b, []byte("fcTL")...)
}

func savePNGChunk(chunk pngData){

	var pngData []byte

	//write PNG header
	var pngHeader []uint8
	writePNGHeader(&pngHeader)
	pngData = append(pngData, pngHeader...)

	//Write IHDR chunk
	writeChunk(chunk.ihdr,"IHDR",&pngData)

	//Write IDAT chunk
	writeChunk(chunk.idat,"IDAT",&pngData)
	
	//Write IEND chunk
	writeChunk(nil,"IEND",&pngData)

	f, _ := os.Create("test4.png")

	_, err := f.Write(pngData)
	if err != nil {
		fmt.Println(err)
	}
}
