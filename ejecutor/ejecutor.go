package ejecutor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"unsafe"
	"bufio"
	str "../structs"
)

var comandos []str.Comando

//
func CrearDisco() {

	fmt.Println("Disco creado")

	file, err := os.Create("/home/helmut/Escritorio/disco.bin")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var otro int8 = 0

	s := &otro

	fmt.Println(unsafe.Sizeof(otro))
	//Escribimos un 0 en el inicio del archivo.
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, s)
	escribirBytes(file, binario.Bytes())
	//Nos posicionamos en el byte 1023 (primera posicion es 0)
	file.Seek(1023, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

	//Escribimos un 0 al final del archivo.
	var binario2 bytes.Buffer
	binary.Write(&binario2, binary.BigEndian, s)
	escribirBytes(file, binario2.Bytes())

	cadena, err :=retrieveROM("/home/helmut/Escritorio/disco.bin")
	_=err

	fmt.Println("NO SE", cadena)

}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func retrieveROM(filename string) ([]byte, error) {
    file, err := os.Open(filename)

    if err != nil {
        return nil, err
    }
    defer file.Close()

    stats, err := file.Stat()
    if err != nil {
        return nil, err
    }

    var size int64 = stats.Size()
    bytes := make([]byte, size)

    bufr := bufio.NewReader(file)
    _,err = bufr.Read(bytes)

    return bytes, err
}