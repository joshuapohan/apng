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
	buffer []byte
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

func (ap *APNGModel) GetPNGChunk(imgBuffer *bytes.Buffer) (pngData){
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

func (ap *APNGModel) writeChunk(chunk []byte, header string, toChunk *[]byte){
	chunkLe := make([]byte, 4)
	chunkTagVal := make([]byte,0, len(chunk) + 8)

	writeUint32(chunkLe, uint32(len(chunk)))

	chunkTagVal = append(chunkTagVal, []byte(header)...)
	chunkTagVal = append(chunkTagVal, chunk...)

	writeCRC32(&chunkTagVal)
	*toChunk = append(*toChunk, chunkLe...)
	*toChunk = append(*toChunk, chunkTagVal...)
}

func (ap *APNGModel) writePNGHeader(){
	ap.buffer = append(ap.buffer, 0x89,0x50,0x4E,0x47,0x0D,0x0A,0x1A,0x0A)
}

func (ap *APNGModel) WriteIHDR(chunk pngData){
	ap.writeChunk(chunk.ihdr, "IHDR", &ap.buffer)
}

func (ap *APNGModel) WriteacTL(img image.Image){
	tmpBuffer := []byte{}

	nbFrames := make([]byte, 4)
	writeUint32(nbFrames, uint32(len(ap.images)))
	tmpBuffer = append(tmpBuffer, nbFrames...)

	nbLoop := make([]byte, 4)
	writeUint32(nbLoop, 0)
	tmpBuffer = append(tmpBuffer, nbLoop...)

	ap.writeChunk(tmpBuffer, "acTL", &ap.buffer)
}

func (ap *APNGModel) WritefcTL(seqNb int, img image.Image, delay int){
	fcTLValue := make([]byte, 4)
	writeUint32(fcTLValue, uint32(seqNb))
	//width
	appendUint32(&fcTLValue, uint32(img.Bounds().Max.X - img.Bounds().Min.X))
	//height
	appendUint32(&fcTLValue, uint32(img.Bounds().Max.Y - img.Bounds().Min.Y))
	//x_offset
	appendUint32(&fcTLValue, uint32(img.Bounds().Min.X))
	//y_offset
	appendUint32(&fcTLValue, uint32(img.Bounds().Min.Y))
	//delay_num
	appendUint16(&fcTLValue, uint16(delay))
	//delay_den
	appendUint16(&fcTLValue, uint16(100))
	//dispose_op
	appendUint8(&fcTLValue, uint8(0))
	//blend_op
	appendUint8(&fcTLValue, uint8(0))
	
	ap.writeChunk(fcTLValue, "fcTL", &ap.buffer)
}

func (ap *APNGModel) WriteIDAT(chunk pngData){
	ap.writeChunk(chunk.idat, "IDAT", &ap.buffer)
}

func (ap *APNGModel) WritefDAT(seqNb int, chunk pngData){
	fDatValue := make([]byte, 4)
	writeUint32(fDatValue, uint32(seqNb))
	fDatValue = append(fDatValue, chunk.idat...)
	ap.writeChunk(fDatValue, "fdAT", &ap.buffer)
}

func (ap *APNGModel) WriteIENDHeader(){
	empty := make([]byte,0)
	ap.writeChunk(empty, "IEND", &ap.buffer)
}

func (ap *APNGModel) Encode(){
	seqNb := 0
	for index, img := range ap.images{
		curImgBuffer := new(bytes.Buffer)
		if err := png.Encode(curImgBuffer, img); err != nil{
			fmt.Println(err)
			return
		}
		curPngChunk := ap.GetPNGChunk(curImgBuffer)
		if(index == 0){
			ap.writePNGHeader()
			ap.WriteIHDR(curPngChunk)
			ap.WriteacTL(img)
			ap.WritefcTL(seqNb, img, ap.delays[index])
			seqNb++
			ap.WriteIDAT(curPngChunk)
		}else{
			ap.WritefcTL(seqNb, img, ap.delays[index])
			seqNb++
			ap.WritefDAT(seqNb, curPngChunk)
			seqNb++
		}
	}
	ap.WriteIENDHeader()
}

func (ap *APNGModel) SavePNGData(path string){

	f, _ := os.Create(path)

	_, err := f.Write(ap.buffer)
	if err != nil {
		fmt.Println(err)
	}
}

/******************************************************
                    byte manipulation
*******************************************************/
func appendUint8(b *[]uint8, u uint8){
	tmp := make([]byte, 1)
	writeUint8(tmp, u)
	*b = append(*b, tmp...)
}

func appendUint16(b *[]uint8, u uint16){
	tmp := make([]byte, 2)
	writeUint16(tmp, u)
	*b = append(*b, tmp...)
}

func appendUint32(b *[]uint8, u uint32){
	tmp := make([]byte, 4)
	writeUint32(tmp, u)
	*b = append(*b, tmp...)
}

func writeUint8(b []uint8, u uint8) {
	b[0] = uint8(u)
}

func writeUint16(b []uint8, u uint16) {
	b[0] = uint8(u >> 8)
	b[1] = uint8(u)
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

func writeCRC32(data *[]byte){
		crcBytes := make([]byte, 4)
		crc := crc32.NewIEEE()
		crc.Write(*data)
		writeUint32(crcBytes, crc.Sum32())
		*data = append(*data, crcBytes...)
}
