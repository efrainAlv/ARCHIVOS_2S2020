package analizador

import (
	"bufio"
	"fmt"
	"os"
	conv "strconv"
	s "strings"

	e "../ejecutor"
	str "../structs"
	main "../main"
)

var comandosLeidos = make([]str.Comando, 0)

var palabraInicial string = ""

var comandoExtendido bool //si vienen el token "\*" es verdadero

func analizarcomando(lineacomandos string, inicial string) {

	var comandos []string
	if s.Contains(lineacomandos, "\"") {
		comandos = regresarSinComillas(lineacomandos)
	} else {
		comandos = s.Split(lineacomandos, " ")
	}

	if inicial == "vacio" {
		comandosLeidos = make([]str.Comando, 0)
		inicial, palabraInicial = s.ToLower(comandos[0]), s.ToLower(comandos[0])
		comandoLeido := str.Comando{Nombre: inicial, Valor: lineacomandos}
		comandosLeidos = append(comandosLeidos, comandoLeido)
	}

	for i := 0; i < len(comandos); i++ {

		//comandos[i] = s.ToLower(comandos[i])

		fmt.Println("ANALIZANDO ===========================================", comandos[i])

		/*
			if n == 0 && s.Contains(comandos[i], "\"") {
				n = i
			} else if i > n && n > 0 && !s.Contains(comandos[i], "\"") { //fmt.Println("Valor Lista", comandosLeidos[n])
				comandosLeidos[n].Valor = comandosLeidos[n].Valor + " " + comandos[i]
				goto finSwitch

			} else if i > n && s.Contains(comandos[i], "\"") {
				comandosLeidos[n].Valor = comandosLeidos[n].Valor + " " + comandos[i]
				//fmt.Println("N ES DIFERENTE DE I, Y CONTIENE COMILLAS")
				n = 0
			}
		*/

		if inicial == "pause" {
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
			if s.Contains(s.ToLower(comandos[i]), "-path->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}
			break

		case "mkdisk":

			if s.Contains(s.ToLower(comandos[i]), "-size->") {

				arr := s.Split(comandos[i], "->")

				//param := s.TrimPrefix(comandos[i], "-size->")

				//fmt.Println("SIZE ENCONTRADO: ", param)
				comandoLeido := str.Comando{Nombre: "-size", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-path->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-name->") {
				//param := s.TrimPrefix(comandos[i], "-name->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("NAME ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-name", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-unit->") {
				//param := s.TrimPrefix(comandos[i], "-unit->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("UNIT ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-unit", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}

			break

		case "rmdisk":
			if s.Contains(s.ToLower(comandos[i]), "-path->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("SIZE ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}
			break

		case "fdisk":

			if s.Contains(s.ToLower(comandos[i]), "-size->") {

				arr := s.Split(comandos[i], "->")

				//param := s.TrimPrefix(comandos[i], "-size->")

				//fmt.Println("SIZE ENCONTRADO: ", param)
				comandoLeido := str.Comando{Nombre: "-size", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-unit->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-unit", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-path->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-path", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-type->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-type", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-fit->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-fit", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-name->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-name", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-delete->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-delete", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)

			} else if s.Contains(s.ToLower(comandos[i]), "-add->") {
				//param := s.TrimPrefix(comandos[i], "-path->")
				arr := s.Split(comandos[i], "->")
				//fmt.Println("PATH ENCONTRADO: ", param)

				comandoLeido := str.Comando{Nombre: "-add", Valor: arr[1]}
				comandosLeidos = append(comandosLeidos, comandoLeido)
			}

		}

	}

finEstados:
	//fmt.Println("----------------------")
	for i := 0; i < len(comandosLeidos); i++ {
		fmt.Println("comando: ", comandosLeidos[i].Nombre, "VALOR", comandosLeidos[i].Valor)
	}
	//fmt.Println("----------------------")

}

func analizarParametros(comms []str.Comando) {

	comandoInicial := comms[0].Nombre

	switch comandoInicial {

	//*****************************************************	 MDISK *********************************************************************
	case "mkdisk":

		var tamanioDisco int64 = 0
		var ruta string = ""
		var nombre string = ""
		var unidad int64 = 1048576

		for i := 1; i < len(comms); i++ {
			if comms[i].Nombre == "-size" {
				tamanio, err := conv.ParseInt(comms[i].Valor, 10, 64)
				if err != nil {
					panic(err)
				}
				tamanioDisco = tamanio
			} else if comms[i].Nombre == "-path" {
				if s.Contains(comms[i].Valor, "\"") {
					nuevoS := s.Replace(comms[i].Valor, "\"", "", 2)
					ruta = nuevoS
				} else {
					ruta = comms[i].Valor
				}

			} else if comms[i].Nombre == "-name" {
				nombre = comms[i].Valor
			} else if comms[i].Nombre == "-unit" {
				if comms[i].Valor == "k" {
					unidad = 1024
				} else if comms[i].Valor == "m" {
					unidad = 1024 * 1024
				}
			}
		}
		fmt.Println("TAMAÑO DE DISCO", tamanioDisco)
		fmt.Println("URL DE DISCO", ruta)
		fmt.Println("NOMBRE DE DISCO", nombre)
		fmt.Println("UNIDAD DE DISCO", unidad)

		tamanioTotal := tamanioDisco * unidad

		fmt.Println("TOTAL TAMAÑO DE DISCO", tamanioTotal)

		if tamanioTotal <= int64(0) || tamanioTotal > 2147483647 || ruta == "" || nombre == "" {
			fmt.Println("*************************************************************")
			fmt.Println("*                          ALERTA                           *")
			fmt.Println("*************************************************************")
			fmt.Println("*   EL TAMAÑO DEL DISCO, LA RUTA O EL NOMBRE DEL DISCO NO   *")
			fmt.Println("*                  NOS SON VALIDOS                          *")
			fmt.Println("*************************************************************")

		} else {
			e.CrearDisco(tamanioTotal, ruta, nombre)
			e.CrearMBR(tamanioTotal, ruta+nombre)
			contenido, err := e.LeerDisco(ruta + nombre)
			if err != nil {
				panic(err)
			}
			fmt.Println(contenido)
			e.MontarMBR(contenido)
			/*
				var buffer bytes.Buffer
				buffer.WriteString("AQUI COMIENZA LA PARTICION UNO, ESPERO ESTO NO AUMENTE EL TAMAÑO DEL DISCO")
				e.EditarArchivo(ruta+nombre, buffer.Bytes(), 138)
			*/
			fmt.Println("*************************************************************")
			fmt.Println("*              ¡DISCO CREADO CON EXITO!                     *")
			fmt.Println("*************************************************************")
		}
		break

		//*****************************************************	 FDISK *********************************************************************

	case "fdisk":

		var tamanioPart uint32 = 0
		var unidad uint32 = 1024
		var ruta string = ""
		var tipo byte = 'P'
		var ajuste byte = 'W'
		var nombre string
		var agregar uint32 = 0
		var eliminar string = ""

		for i := 0; i < len(comms); i++ {

			if comms[i].Nombre == "-size" {
				tamanio, err := conv.ParseInt(comms[i].Valor, 10, 64)
				if err != nil {
					panic(err)
				}
				tamanioPart = uint32(tamanio)

			} else if comms[i].Nombre == "-unit" {
				if comms[i].Valor == "k" || comms[i].Valor == "K" {
					unidad = 1024
				} else if comms[i].Valor == "m" || comms[i].Valor == "M" {
					unidad = 1024 * 1024
				} else if comms[i].Valor == "b" || comms[i].Valor == "B" {
					unidad = 1
				}

			} else if comms[i].Nombre == "-path" {

				if s.Contains(comms[i].Valor, "\"") {
					nuevoS := s.Replace(comms[i].Valor, "\"", "", 2)
					ruta = nuevoS
				} else {
					ruta = comms[i].Valor
				}

			} else if comms[i].Nombre == "-type" {

				tipo = comms[i].Valor[0]

			} else if comms[i].Nombre == "-fit" {
				ajuste = comms[i].Valor[0]
			} else if comms[i].Nombre == "-name" {
				if s.Contains(comms[i].Valor, "\"") {
					nuevoS := s.Replace(comms[i].Valor, "\"", "", 2)
					nombre = nuevoS
				} else {
					nombre = comms[i].Valor
				}
			} else if comms[i].Nombre == "-delete" {
				eliminar = comms[i].Valor
			} else if comms[i].Nombre == "-add" {
				add, err := conv.ParseInt(comms[i].Valor, 10, 64)
				if err != nil {
					panic(err)
				}
				agregar = uint32(add)
			}

		}

		var nombrePart [16]byte
		for i := 0; i < 16; i++ {
			nombrePart[i] = nombre[i]
		}

		_ = eliminar
		_ = agregar

		tamanioTotal := tamanioPart * unidad
		if tamanioPart <= 0 || tamanioTotal > 2147483647 || ruta == "" || nombre == "" {
			fmt.Println("*************************************************************")
			fmt.Println("*                          ALERTA                           *")
			fmt.Println("*************************************************************")
			fmt.Println("*   EL TAMAÑO DEL DISCO, LA RUTA O EL NOMBRE DEL DISCO NO   *")
			fmt.Println("*                  NOS SON VALIDOS                          *")
			fmt.Println("*************************************************************")

		} else {

			e.CrearParticion(ruta, tipo, ajuste, tamanioTotal, nombrePart)
			fmt.Println("*************************************************************")
			fmt.Println("*              ¡PARTICION CREADA CON EXITO!                 *")
			fmt.Println("*************************************************************")
		}

		break

	}

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

//SI HAY UNA LINEA DE CODIGO CON COMILLAS, ESTE METODO REGRESA CADA COMANDO SIN PROBLEMAS
func regresarSinComillas(cadena string) (comandos []string) {

	var comandoss []string

	subCadena := s.Split(cadena, " ")
	palabra := ""

	for i := 0; i < len(subCadena); i++ {

		palabra = subCadena[i]
		if s.Contains(subCadena[i], "\"") {

			palabra = subCadena[i] + " "
			for j := i + 1; j < len(subCadena); j++ {
				if !s.Contains(subCadena[j], "\"") {
					palabra += subCadena[j] + " "
				} else {
					palabra += subCadena[j]
					i = j
					break
				}
			}
		}
		comandoss = append(comandoss, palabra)
	}

	return comandoss

}
