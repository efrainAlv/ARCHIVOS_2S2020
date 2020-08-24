package main

//CADA STRUCT PARTICION PESA 27 BYTES
//EL MBR PESA 138 BYTES

import (
	"bytes"
	"encoding/binary"
	"fmt"

	//a "./analizador"
	e "./ejecutor"
)

type ejem struct {
	nombre string
	valor  int
}

func main() {

	fmt.Println("Hola mundo")
	//a.Leer("/home/helmut/Escritorio/prueba.txt")

	buf := make([]byte, 8)

	buf[0] = 'h'
	buf[1] = 'o'
	buf[2] = 'l'
	buf[3] = 'a'
	buf[4] = 'h'
	buf[5] = 'o'
	buf[6] = 'l'
	buf[7] = 'a'
	n := binary.PutUvarint(buf, 4)
	fmt.Println(n)

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, buf)

	fmt.Println(buffer.Bytes())

	e.EditarArchivo("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", buffer.Bytes(), 138)

}
