package main

import (
	"fmt"

	a "./analizador"

	/*"bytes"
	"encoding/binary"
	"strconv"
	"time"
	"unsafe"
	*/
	//e "./ejecutor"
	//str "./structs"
)

type ejem struct {
	nombre string
	valor  int
}

func main() {

	fmt.Println("Hola mundo")

	a.Leer("/home/helmut/Escritorio/prueba.txt")

	//	b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}
	/*
		anyo, mes, dia := time.Now().Date()
		hora, min, sec := time.Now().Clock()

		fecha := strconv.Itoa(anyo) + "-" + string(mes) + "-" + string(dia)
		horaFecha := string(hora) + "-" + string(min) + "-" + string(sec)

		var nuevo str.MBR

		var buffer bytes.Buffer
		var cadena1 [19]byte
		var cadena2 []byte

		buffer.WriteString("Hola que tal como estas")
		cadena2 = buffer.Bytes()
		buffer.Reset()

		for i := 0; i < len(cadena1); i++ {
			cadena1[i] = cadena2[i]
		}

		nuevo = str.MBR{Tamanio: 1000000, Fecha: cadena1, Firma: 24654}

		binary.Write(&buffer, binary.BigEndian, nuevo)

		fmt.Println("TAMAÑO DEL MBR", unsafe.Sizeof(nuevo))
		fmt.Println("CODIFICADO", buffer)

		fmt.Println(fecha)
		fmt.Println(horaFecha)

		fmt.Println("************************************************************************")

		e.EditarArchivo("/home/helmut/Escritorio/disco.dsk", buffer.Bytes(), 1023)
	*/
}
