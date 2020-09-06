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

func main() {

	fmt.Println("Hola mundo")
	//a.Leer("/home/helmut/Escritorio/prueba.txt")
	//fun()
	//ParticionesMontadas = append(ParticionesMontadas)
	//e.MontarParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", "hola Munda")

	str.ParticionesMontadas = append(str.ParticionesMontadas, str.ParticionMontada{Particion: str.Particion{Inicio: 500, Tamanio:8000}, Letra: 97, Numero: 1, Ruta: "/home/helmut/Escritorio/Mis Discos/Disco_3.dsk"})
	str.ParticionesMontadas = append(str.ParticionesMontadas, str.ParticionMontada{Particion: str.Particion{}, Letra: 100, Numero: 2, Ruta: "/home/helmut/Escritorio/Mis Discos/Disco_3.dsk"})
	str.ParticionesMontadas = append(str.ParticionesMontadas, str.ParticionMontada{Particion: str.Particion{}, Letra: 99, Numero: 3, Ruta: "/home/helmut/Escritorio/Mis Discos/Disco_3.dsk"})
	str.ParticionesMontadas = append(str.ParticionesMontadas, str.ParticionMontada{Particion: str.Particion{}, Letra: 98, Numero: 4, Ruta: "/home/helmut/Escritorio/Mis Discos/Disco_3.dsk"})
	//e.MontarParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", "hola Munda")
	e.FormatearParticion("dva1", "full")

	/*
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

	fmt.Println("CANTIDAD DE ESTRUCTURAS", e.CalcularNumeroDeEstructuras(2000000))
	fmt.Println(str.TamSuperBoot)
	fmt.Println(str.TamAVD)
	fmt.Println(str.TamDetalleDirect)
	fmt.Println(str.TamInodo)
	fmt.Println(str.TamBloque)
	fmt.Println(str.TamBitacora)

	//97-122
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
