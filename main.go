package main

import (
	"fmt"
	//a "./analizador"
	e "./ejecutor"
)


type ejem struct {

	nombre string
	valor int
}

func main() {

	fmt.Println("Hola mundo")
	//a.Leer("/home/helmut/Escritorio/prueba.txt")

	//	b := []byte{'g', 'o', 'l', 'a', 'n', 'g'}

	
	e.CrearArchivoBinario()

	e.EditarArchivo("/home/helmut/Escritorio/disco.dsk", "HOLA QUE TAL COMO ESTAS?", 0)

}
