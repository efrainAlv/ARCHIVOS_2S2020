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
	"bufio"
	"os"
	"strings"

	a "./analizador"
	e "./ejecutor"
	//str "./structs"
)

//

func main() {

	fmt.Print("Hola mundo ")
	fmt.Println(a.NoMolestar)

	salida := false
	for i := 0; !salida; i++ {

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("")
		fmt.Println("*************************************************************")
		fmt.Println("*                        BIENVENIDO                         *")
		fmt.Println("*************************************************************")
		fmt.Println("*   INGRESE EL COMANDO INICIAL, ESCRIBA EXIT PARA SALIR     *")
		fmt.Println("*************************************************************")
		fmt.Println("")
		
		entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
		eleccion := strings.TrimRight(entrada, "\r\n") // Remover el salto de línea de la entrada del usuario
		fmt.Println("")
		
		if strings.Contains(eleccion, "exec") {
			if strings.Contains(eleccion, "-path->") {
				path := strings.Split(eleccion, "->")[1]

				a.Leer(path)
			}
		} else {
			if eleccion == "exit" || eleccion == "EXIT" {
				salida = true
			}
		}

	}

	//fun()
	//e.MontarParticion("/home/helmut/Escritorio/Mis Discos/Disco_3.dsk", "hola Mundo")
	//e.FormatearParticion("vda1", "full")
	//e.CrearRoot("vda1", 1, 1, 777)
	/*ruta := []string{"/", "home", "user"}
	var cont string = "123456789123456789123456789123456789123456789123456789123456789123456789123456789"
	paso := e.CrearArchivo("vda1", cont, ruta, "ejemplo.txt", 1, 1, 777)
	*/
	//fmt.Println(paso)

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
