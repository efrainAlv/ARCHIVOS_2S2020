package disk

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"unsafe"
	"fmt"
)

type mbr struct {
}

func leerArchivo(direccion string) {

	file, err := os.Open(direccion)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

}

func leerBytes(file *os.File, numero int) []byte {

	bytes := make([]byte, numero)

	_, err := file.Read(bytes)

	if err != nil {
		log.Fatal(err)
	}

	return bytes
}


func escribirBytes(file *os.File, bytes []byte){

	_, err := file.Write(bytes)

	if err != nil{
		log.Fatal(err)
	}

}


func escribirArchivo(direccion string, tamanioDisco int64) {

	file, err := os.Create(direccion)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	var tamanio int8 = 0;

	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, tamanio)

	escribirBytes(file, binario.Bytes())

	file.Seek(tamanioDisco, 0)

	


}
