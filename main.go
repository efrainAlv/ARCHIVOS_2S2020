package main

//CADA STRUCT PARTICION PESA 27 BYTES
//EL MBR PESA 138 BYTES

import (
	"fmt"
	/*
		"os"
		"log"
		"bytes"
	*/
	//a "./analizador"
	e "./ejecutor"
	str "./structs"
)

//
var ParticionesMontadas []str.ParticionMontada

func main() {

	fmt.Println("Hola mundo")
	//a.Leer("/home/helmut/Escritorio/prueba.txt")
	//fun()
	e.MontarParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", "hola Mundaa")

}

func fun() {

	var tipo byte = 'P'
	var ajuste byte = 'B'
	var tamanio uint32 = uint32(50)
	var nombre [16]byte
	nombre[0] = 'h'
	nombre[1] = 'o'
	nombre[2] = 'l'
	nombre[3] = 'a'
	nombre[4] = ' '
	nombre[5] = 'M'
	nombre[6] = 'u'
	nombre[7] = 'n'
	nombre[8] = 'd'
	nombre[9] = 'a'

	fmt.Println(uint16(9999))

	e.CrearParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", tipo, ajuste, tamanio, nombre)
	contenido, err := e.LeerDisco("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk")
	if err != nil {
		panic(err)
	}

	e.MontarMBR(contenido)

}
