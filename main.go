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
	a "./analizador"
	e "./ejecutor"
	//str "./structs"
)

//

func main() {

	fmt.Print("Hola mundo ")
	fmt.Println(a.NoMolestar)

	//a.Leer("/home/helmut/Escritorio/prueba.txt")
	//fun()
	//e.MontarMBR(contenido)
	e.MontarParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", "hola Mundo")
	ruta := []string{"/", "home", "user"}
	e.CrearAVDInicio("vda1", ruta, uint32(5), uint32(5), uint16(777), "123456789123456789123456789123456789123456789123456789123456789123456789123456789123456789123456789123456789")
	//e.FormatearParticion("vda1", "full")
	//e.CrearRoot("vda1", 1, 1, 777)
	/*contenido, err := e.LeerDisco("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk")
	if err != nil {
		panic(err)
	}
	fmt.Println(contenido)
	*/
	/*
		diciembre
			fmt.Println("")
			fmt.Println("PARTICIONES:")
			for i := 0; i < len(str.ParticionesMontadas); i++ {
				fmt.Println(str.ParticionesMontadas[i])
			}

			e.DesmontarParticion("vda1")
			fmt.Println("")
			fmt.Println("PARTICIONES:")
			for i := 0; i < len(str.ParticionesMontadas); i++ {
				fmt.Println(str.ParticionesMontadas[i])
			}
	*/

	//97-122
}

func fun() {

	var tipo byte = 'P'
	var ajuste byte = 'B'
	var tamanio uint32 = uint32(6000)
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
	nombre[9] = 'o'

	e.CrearParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", tipo, ajuste, tamanio, nombre)
	contenido, err := e.LeerDisco("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk")
	if err != nil {
		panic(err)
	}

	e.MontarMBR(contenido)

}
