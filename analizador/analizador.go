package analizador

import (
	"bufio"
	"fmt"
	"os"
	s "strings"

	e "../ejecutor"
	str "../structs"
)

var comandosLeidos = make([]str.Comando, 0)

var palabraInicial string = ""

var comandoExtendido bool //si vienen el token "\*" es verdadero

func analizarcomando(lineacomandos string, inicial string) {

	n := 0
	_ = n

	comandos := s.Split(lineacomandos, " ")

	if inicial == "vacio" {
		comandosLeidos = make([]str.Comando, 0)
		inicial, palabraInicial = comandos[0], comandos[0]
		comandoLeido := str.Comando{Nombre: comandos[0], Valor: lineacomandos}
		comandosLeidos = append(comandosLeidos, comandoLeido)
	}

	for i := 0; i < len(comandos); i++ {

		comandos[i] = s.ToLower(comandos[i])

		fmt.Println("----------------------------- ANALIZANDO: ", i+1, "-", len(comandos), ": ", comandos[i], "-----------------------------------")

		if n == 0 && s.Contains(comandos[i], "\"") {
			n = i

		} else if i > n && n > 0 && !s.Contains(comandos[i], "\"") {
			comandosLeidos[n].Valor = comandosLeidos[n].Valor + " " + comandos[i]

		} else if i > n && s.Contains(comandos[i], "\"") {
			comandosLeidos[n].Valor = comandosLeidos[n].Valor + " " + comandos[i]
			n = 0
		}

		if inicial == "pausa" {
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}

		if comandos[i] == "\\*" {
			comandoExtendido = true
			goto finEstados
		}
		if i == len(comandos)-1 {
			comandoExtendido = false
		}

		switch inicial {
		case "exec":
			if s.Contains(comandos[i], "-path->") {
				param := s.TrimPrefix(comandos[i], "-path->")

				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: param}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}
			break

		case "Mkdisk":

			if s.Contains(comandos[i], "-size->") {
				param := s.TrimPrefix(comandos[i], "-size->")

				//fmt.Println("SIZE ENCONTRADO: ", param)
				comandoLeido := str.Comando{Nombre: "-size", Valor: param}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(comandos[i], "-path->") {
				param := s.TrimPrefix(comandos[i], "-path->")

				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: param}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(comandos[i], "-name->") {
				param := s.TrimPrefix(comandos[i], "-name->")

				//fmt.Println("NAME ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-name", Valor: param}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(comandos[i], "-unit->") {
				param := s.TrimPrefix(comandos[i], "-unit->")

				//fmt.Println("UNIT ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-unit", Valor: param}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}

			break

		case "rmDisk":
			if s.Contains(comandos[i], "-path->") {
				param := s.TrimPrefix(comandos[i], "-path->")

				//fmt.Println("SIZE ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: param}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}
			break

		}

	}

finEstados:
	fmt.Println("----------------------")
	for i := 0; i < len(comandosLeidos); i++ {
		fmt.Println("comando: ", comandosLeidos[i].Nombre, "VALOR", comandosLeidos[i].Valor)
	}
	fmt.Println("----------------------")

}

func analizarParametros([]str.Comando) {

	fmt.Println("00000000000000000000000000 comando EJECUTADO 00000000000000000000000000")
	e.CrearDisco()
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
			analizarcomando(cadena[i], "vacio")
		} else {
			analizarcomando(cadena[i], comandosLeidos[0].Nombre)
		}

		if !comandoExtendido {
			analizarParametros(comandosLeidos)
		}
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
