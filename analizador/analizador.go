package analizador

import (
	"bufio"
	"fmt"
	"os"
	s "strings"
)

var estado int = 0
var comandosLeidos = make([]comando, 0)

type comando struct {
	nombre string
	valor  string
}

func analizarComando(lineaComnados string) {

	comandos := s.Split(lineaComnados, " ")

	var inicial string = comandos[0]

	for i := 0; i < len(comandos); i++ {

		/*if comandos[i] == "pausa" {
			lector := bufio.NewReader(os.Stdin)
			fmt.Println("En pausa ...")
			input, _:=lector.ReadString('\n')
			_ = input
		}*/

		if inicial == "pausa" {
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		if inicial == "exec" {
			path := s.TrimPrefix(comandos[i], "-path->")
			fmt.Println("PATH ENCONTRADO: ", path)
		}

	}

}

func analizarParametros(comadno string) {

	fmt.Println("Hola")

}

func ejecutar() {

}

//
func Leer(url string) {

	file, err := os.Open(url)
	check(err)
	fileInfo, err := os.Lstat(url)
	check(err)

	cadenaBytes := make([]byte, fileInfo.Size()) //OBTIENE LA CADENA DE BYTES DEL ARCHIVO
	check(err)

	n, err := file.Read(cadenaBytes) //SE LEE EL ARCHIVO, SE PASA COMO PARAMETRO EL TAMAÑO EN BYTES DEL ARCHIVO
	check(err)

	fmt.Println("BYTES LEIDOS: ", n)

	cadena := s.Split(string(cadenaBytes), "\n") //SEPARA EL ARCHIVO POR LINEAS

	//LEE EL ARCHIVO LINEA POR LINEA
	for i := 0; i < len(cadena)-1; i++ {
		fmt.Println("Linea ", i+1, ": ", cadena[i])
		analizarComando(cadena[i])
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
