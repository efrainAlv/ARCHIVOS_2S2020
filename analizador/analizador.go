package analizador

import (
	"bufio"
	"fmt"
	"os"
	s "strings"
)

var estado int = 0
var comandosLeidos = make([]comando, 0)

var comandoExtendido bool

type comando struct {
	nombre string
	valor  string
}

func analizarComando(lineaComandos string, inicial string) {

	comandos := s.Split(lineaComandos, " ")

	if inicial == "vacio"{
		inicial  = comandos[0]

	}

	for i := 0; i < len(comandos); i++ {

		/*if comandos[i] == "pausa" {
			lector := bufio.NewReader(os.Stdin)
			fmt.Println("En pausa ...")
			input, _:=lector.ReadString('\n')
			_ = input
		}*/

		comandos[i] = s.ToLower(comandos[i])

		if inicial == "pausa" {
			bufio.NewReader(os.Stdin).ReadBytes('\n')

		}
		
		if comandos[i]=="\\*"{
			comandoExtendido = true
			break
		}else if comandos[i]=="\\n"{
			analizarParametros(comandosLeidos)
			comandosLeidos = make([]comando, 0)
		}
		switch inicial {
		
		case "exec":
			if s.Contains(comandos[i], "-path->") {
				param := s.TrimPrefix(comandos[i], "-path->")

				fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := comando{"-path", param}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			}
			break

		case "Mkdisk":
			if s.Contains(comandos[i], "-size->") {
				param := s.TrimPrefix(comandos[i], "-size->")

				fmt.Println("SIZE ENCONTRADO: ", param)

				comandoLeido := comando{"-size", param}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			} else if s.Contains(comandos[i], "-path->") {
				param := s.TrimPrefix(comandos[i], "-path->")

				fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := comando{"-path", param}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(comandos[i], "-name->") {
				param := s.TrimPrefix(comandos[i], "-name->")

				fmt.Println("NAME ENCONTRADO: ", param)

				comandoLeido := comando{"-name", param}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(comandos[i], "-unit->") {
				param := s.TrimPrefix(comandos[i], "-unit->")

				fmt.Println("UNIT ENCONTRADO: ", param)

				comandoLeido := comando{"-unit", param}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}
			break

		}

	}

}

func analizarParametros([]comando) {

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
		
		if !comandoExtendido {
			analizarComando(cadena[i], "vacio")
		}else{
			analizarComando(cadena[i], comandosLeidos[0].nombre)
		}
		
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
