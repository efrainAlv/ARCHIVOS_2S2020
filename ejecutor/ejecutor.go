package ejecutor

/*

138 bytes del MBR
27 bytes de cada particion
30 bytes de la infor del MBR

estado 0 	-> 30
tipo 1 		-> 31
ajuste 2
inicio 3
tamanio 7
nombre 23


*/

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	rand "math/rand"
	"os"
	"strconv"
	s "strings"
	"time"

	str "../structs"
	//a "../analizador"
)

var comandos []str.Comando

const tamanioMBR = 138
const tamanioPart = 27
const tamanioInfoMBR = 30

//
func CrearDisco(tamanio int64, url string, nombre string) {

crearArchivo:
	file, err := os.Create(url + nombre)
	defer file.Close()
	if err != nil {
		crearDirectorioSiNoExiste(url)
		goto crearArchivo
	}

	//creamos una variable de tipo int 8 con valor de cero
	cero := int8(0)

	//Nos posicionamos en el byte 1023 que seria el ultimo byte del archivo (el 1023 se puede cambiar dependiendo del tamaño que se quiera el archivo en bytes)
	//seek nos posiciona en el byte 1023 al momento que queramos leer o escribir en el archivo
	file.Seek(tamanio-1, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

	//creamos un buffer para pode leer y escribir en archivos
	buffer := new(bytes.Buffer)

	//el metodo write escribe la representacion binaria de la variable cero en la variable buffer
	if err := binary.Write(buffer, binary.BigEndian, cero); err != nil {
		fmt.Println(err)
	}

	//se crea el archivo con el mismo tamaño en bytes que la variable buffer, y el archivo se mantiene en la variable file
	escribirArchivo(file, buffer.Bytes())

}

//crea un archivo de la longitud del parametro bytes
func escribirArchivo(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

//
func LeerDisco(url string) ([]byte, error) {

	//se obtiene el archivo de la direccion filename, y se guarda en la variable file
	file, err := os.Open(url)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	//se obtienen las propiedades del archivo y se guarda en props
	props, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var size int64 = props.Size()
	contenido := make([]byte, size)

	//se crea un tipo Read con el tamaño del archivo a leer
	bufr := bufio.NewReader(file)

	//con el metodo Read de bufr el arreglo de bytes, "contenido", obtiene el contenido del archivo
	_, err = bufr.Read(contenido)

	return contenido, err
}

//
func editarArchivo(url string, cadena []byte, inicio int64) {
	// Read Write Mode
	file, err := os.OpenFile(url, os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.WriteAt(cadena, inicio) // Write at 0 beginning
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	_ = len
}

func crearDirectorioSiNoExiste(directorio string) {

	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err = os.Mkdir(directorio, 0777)
		if err != nil {
			// Aquí puedes manejar mejor el error, es un ejemplo
			panic(err)
		} else {
			fmt.Println("DIRECTORIO CREADO")
		}
	}
}

//
func CrearMBR(tamanio int64, url string) {
	fechaPart := generarFecha()

	var p1 str.Particion
	var p2 str.Particion
	var p3 str.Particion
	var p4 str.Particion

	p1.Inicio = uint32(tamanioMBR)
	p1.Estado = 1
	p2.Estado = 1
	p3.Estado = 1
	p4.Estado = 1

	mbr := str.MBR{Tamanio: uint32(tamanio), Fecha: fechaPart, Firma: rand.Uint32(), Part1: p1, Part2: p2, Part3: p3, Part4: p4}

	//fmt.Println("VERIFICANDO TAMAÑO DEL DISCO: ", mbr.Tamanio)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, mbr)

	editarArchivo(url, buffer.Bytes(), 0)

}

//
func MontarMBR(contenidoDisco []byte) (MBR str.MBR) {

	mbr := str.MBR{}

	var tamanio = contenidoDisco[:4]
	mbr.Tamanio = convertirBinario(tamanio)

	var fecha = contenidoDisco[4:26]

	for i := 0; i < 22; i++ {
		mbr.Fecha[i] = fecha[i]
	}

	var firma = contenidoDisco[26:30]
	mbr.Firma = convertirBinario(firma)

	//fmt.Printf("%c", mbr)

	mbr = montarParticionesAlMBR(contenidoDisco, mbr)

	fmt.Println(mbr)
	return mbr
}

func montarParticionesAlMBR(contenidoDisco []byte, mbr str.MBR) (mBR str.MBR) {

	n := 0
	for i := 0; i < 4; i++ {

		part := contenidoDisco[30+n : 57+n]
		partN := str.Particion{}

		estadoP := part[0:1]
		partN.Estado = estadoP[0]

		tipoP := part[1:2]
		partN.Tipo = tipoP[0]

		ajusteP := part[2:3]
		partN.Ajuste = ajusteP[0]

		inicioP := part[3:7]
		partN.Inicio = convertirBinario(inicioP)

		tamanioP := part[7:11]
		partN.Tamanio = convertirBinario(tamanioP)

		var nombre = part[11:27]
		for k := 0; k < 16; k++ {
			partN.Nombre[k] = nombre[k]
		}

		switch i {
		case 0:
			mbr.Part1 = partN
			break

		case 1:
			mbr.Part2 = partN
			break

		case 2:
			mbr.Part3 = partN
			break

		case 3:
			mbr.Part4 = partN
			break

		}
		n += 27
	}

	return mbr

}

//
func CrearParticion(url string, tipo byte, ajuste byte, tamanio uint32, nombre [16]byte) (paso bool) {

	contenidoDisco, err := LeerDisco(url)

	if err != nil {
		panic(err)
	}

	mbr := MontarMBR(contenidoDisco)

	var tipoDisponible = false
	var nombreDisponible = true
	var espacioDisp uint32 = 0
	inicioPart := uint32(138)
	numeroPart := 0

	n := 0
	for i := 0; i < 4; i++ {
		if getParticion(mbr, i).Tipo == tipo {
			n++
		}
		if getParticion(mbr, i).Nombre == nombre {
			nombreDisponible = false
		}
	}

	if nombreDisponible {
		if tipo == 'E' && n < 1 {
			tipoDisponible = true

		} else if tipo == 'P' && n < 3 {
			tipoDisponible = true

		} else {
			fmt.Println("*************************************************************")
			fmt.Println("*                          ALERTA                           *")
			fmt.Println("*************************************************************")
			fmt.Println("*      NO HAY ESPACIO PARA UNA PARTICION DEL TIPO ", fmt.Sprintf("%c", tipo), "       *")
			fmt.Println("*************************************************************")
		}
	} else {

		fmt.Println("*************************************************************")
		fmt.Println("*                          ALERTA                           *")
		fmt.Println("*************************************************************")
		fmt.Println("*  EL NOMBRE DE LA PARTICION YA ESXISTE, INTENTE CON OTRO   *")
		fmt.Println("*************************************************************")
		goto final
	}

	if tipoDisponible {

		for i := 0; i < 4; i++ {
			if getParticion(mbr, i).Estado == 1 {
				espacioDisp, inicioPart = consultarEspacioDisponible(mbr, i)
				if espacioDisp >= tamanio {
					numeroPart = i
					break
				} else {
					espacioDisp = 0
				}

			}
		}

	}

	if espacioDisp > 0 {
		particion := getParticion(mbr, numeroPart)

		particion.Estado = 0
		particion.Tipo = tipo
		particion.Ajuste = ajuste
		particion.Inicio = inicioPart
		particion.Tamanio = tamanio
		particion.Nombre = nombre

		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, particion)
		posicion := tamanioInfoMBR + tamanioPart*(numeroPart)
		editarArchivo(url, buffer.Bytes(), int64(posicion))
		goto final
	}

	fmt.Println("")
	fmt.Println("*************************************************************")
	fmt.Println("*                          ALERTA                           *")
	fmt.Println("*************************************************************")
	fmt.Println("*       NO HAY ESPACIO SUFICIENTE PARA LA PARTICION         *")
	fmt.Println("*************************************************************")
	fmt.Println("")

final:
	return false
}

//
func MontarParticion(url string, nombre string) {

	contenidoDisco, err := LeerDisco(url)
	if err != nil {
		panic(err)
	}

	mbr := MontarMBR(contenidoDisco)
	var buffer bytes.Buffer
	buffer.WriteString(nombre)
	var nombrePart [16]byte
	var nombreTemp = buffer.Bytes()
	for i := 0; i < len(nombre); i++ {
		nombrePart[i] = nombreTemp[i]
	}

	part := str.Particion{}
	nombreExiste := false
	for i := 0; i < 4; i++ {
		if part = getParticion(mbr, i); part.Nombre == nombrePart {
			nombreExiste = true
			break
		}
	}

	if nombreExiste {

		if len(str.ParticionesMontadas) == 0 {
			str.ParticionesMontadas = append(str.ParticionesMontadas, str.ParticionMontada{Particion: part, Letra: 97, Numero: 1, Ruta: url})
		} else {

			n := 0
			codigo := byte(97)
			for i := 0; i < len(str.ParticionesMontadas); i++ {
				if str.ParticionesMontadas[i].Ruta == url {
					n++
				}
			}

			for i := 0; i < len(str.ParticionesMontadas); i++ {
				if str.ParticionesMontadas[i].Letra == codigo && str.ParticionesMontadas[i].Ruta == url {
					codigo += byte(1)
					i = 0
				}
			}

			str.ParticionesMontadas = append(str.ParticionesMontadas, str.ParticionMontada{Particion: part, Letra: codigo, Numero: uint16(n + 1), Ruta: url})
		}

	} else {
		fmt.Println("*************************************************************")
		fmt.Println("*                          ALERTA                           *")
		fmt.Println("*************************************************************")
		fmt.Println("*         LA PARTICION NO EXISTE DENTRO DEL DISCO           *")
		fmt.Println("*************************************************************")
	}

}

//
func DesmontarParticion(idParticion string) {

	letra := idParticion[2]
	numero := idParticion[3:len(idParticion)]
	codigoLetra, err := strconv.ParseUint(numero, 10, 16)
	if err != nil {
		panic(err)
	}

	desmontado := false
	for i := 0; i < len(str.ParticionesMontadas); i++ {
		if str.ParticionesMontadas[i].Letra == byte(letra) && str.ParticionesMontadas[i].Numero == uint16(codigoLetra) {
			if i == len(str.ParticionesMontadas) {
				copy(str.ParticionesMontadas, str.ParticionesMontadas[:i-1])
				desmontado = true
				break
			} else {
				str.ParticionesMontadas = append(str.ParticionesMontadas[:i], str.ParticionesMontadas[i+1:]...)
				desmontado = true
				break
			}
		}
	}

	if desmontado {
		fmt.Println("PARTICION DESMONTADA")
	} else {
		fmt.Println("PARTICION NO ENCONTRADA")
	}

}

//
func FormatearParticion(idPart string, tipo string) {

	part := getParticionByID(idPart)
	if tipo == "fast" {

	} else if tipo == "full" {
		var cadena []byte
		for i := 0; i < int(part.Particion.Tamanio); i++ {
			cadena = append(cadena, 0)
		}

		rutaDisco := s.Split(part.Ruta, "/")
		editarArchivo(part.Ruta, cadena, int64(part.Particion.Inicio))
		superBoot := CrearSuperBoot(part.Particion.Tamanio, rutaDisco[len(rutaDisco)])
		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, superBoot)
		editarArchivo(part.Ruta, buffer.Bytes(), int64(part.Particion.Inicio))
	}

}

//
func CrearSuperBoot(tamanioPart uint32, nombre string) str.SuperBoot {

	nEstructuras := calcularNumeroDeEstructuras(tamanioPart)

	var nombreDisco [20]byte
	for i := 0; i < len(nombreDisco); i++ {
		nombreDisco[i] = nombre[i]
	}
	cantidadAVD := nEstructuras
	cantidadDetalleDirect := nEstructuras
	cantidadInodos := nEstructuras * 5
	cantidadBloques := nEstructuras * 20
	cantidadAVDLibres := cantidadAVD
	cantidadDetalleDirecttLibres := cantidadDetalleDirect
	cantidadInodosLibres := cantidadInodos
	cantidadBloquesLibres := cantidadBloques

	fechaCreacion := generarFecha()
	fechaUltimoMontaje := generarFecha()
	numeroMontajes := uint16(1)

	apuntadorBitMapAVD := uint32(str.TamSuperBoot)
	apuntadorAVD := apuntadorBitMapAVD + cantidadAVD
	apuntadorBitMapDetalleDirect := apuntadorAVD + cantidadAVD*uint32(str.TamAVD)
	apuntadorDetalleDirect := apuntadorBitMapDetalleDirect + cantidadDetalleDirect
	apuntadorBitMapInodos := apuntadorDetalleDirect + cantidadDetalleDirect*uint32(str.TamDetalleDirect)
	apuntadorInodos := apuntadorBitMapInodos + cantidadDetalleDirect
	apuntadorBitMapBloques := apuntadorInodos + cantidadInodos*uint32(str.TamInodo)
	apuntadorBloques := apuntadorBitMapBloques + cantidadBloques
	apuntadorBitacora := apuntadorBloques + cantidadBloques*uint32(str.TamBloque)

	tamanioAVD := uint32(str.TamAVD)
	tamanioDetalleDirect := uint32(str.TamDetalleDirect)
	tamanioInodo := uint32(str.TamInodo)
	tamanioBloque := uint32(str.TamBloque)
	tamanioBitacora := uint32(str.TamBitacora)

	primerAVDLibre := apuntadorAVD
	primerDetalleDirectLibre := apuntadorDetalleDirect
	primerInodoLibre := apuntadorInodos
	primerBloqueLibre := apuntadorBloques

	superboot := str.SuperBoot{
		NombreDisco:                  nombreDisco,
		CantidadAVD:                  cantidadAVD,
		CantidadDetalleDirect:        cantidadDetalleDirect,
		CantidadInodos:               cantidadInodos,
		CantidadBloques:              cantidadBloques,
		CantidadAVDLibres:            cantidadAVDLibres,
		CantidadDetalleDirectLibres:  cantidadDetalleDirecttLibres,
		CantidadInodosLibres:         cantidadInodosLibres,
		CantidadBloquesLibres:        cantidadBloquesLibres,
		FechaCreacion:                fechaCreacion,
		FechaUltimoMontaje:           fechaUltimoMontaje,
		NumeroMontajes:               numeroMontajes,
		ApuntadorAVD:                 apuntadorAVD,
		ApuntadorBitMapAVD:           apuntadorBitMapAVD,
		ApuntadorDetalleDirect:       apuntadorDetalleDirect,
		ApuntadorBitMapDetalleDirect: apuntadorBitMapDetalleDirect,
		ApuntadorInodos:              apuntadorInodos,
		ApuntadorBitMapInodos:        apuntadorBitMapInodos,
		ApuntadorBloques:             apuntadorBloques,
		ApuntadorBitMapBloques:       apuntadorBitMapBloques,
		ApuntadorBitacora:            apuntadorBitacora,
		TamanioAVD:                   tamanioAVD,
		TamanioDetalleDirect:         tamanioDetalleDirect,
		TamanioInodo:                 tamanioInodo,
		TamanioBloque:                tamanioBloque,
		TamanioBitacora:              tamanioBitacora,
		PrimerAVDLibre:               primerAVDLibre,
		PrimerDetalleDirectLibre:     primerDetalleDirectLibre,
		PrimerInodoLibre:             primerInodoLibre,
		PrimerBloqueLibre:            primerBloqueLibre,
		NumeroMagico:                 uint32(20171350)}

	return superboot
}

func generarFecha() (fechaReturn [22]byte) {

	anyo, mes, dia := time.Now().Date()
	hora, min, sec := time.Now().Clock()

	fecha := strconv.Itoa(anyo) + "-" + mes.String() + "-" + strconv.Itoa(dia)
	horaFecha := strconv.Itoa(hora) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)

	var buffer bytes.Buffer
	fecha = fecha + " " + horaFecha

	//fmt.Println("FECHA: ", fecha)

	buffer.Reset()
	buffer.WriteString(fecha)
	cadena2 := buffer.Bytes()

	n := 0
	if len(fechaReturn) < len(cadena2) {
		n = len(fechaReturn)
	} else {
		n = len(cadena2)
	}
	for i := 0; i < n; i++ {
		fechaReturn[i] = cadena2[i]
	}

	return fechaReturn
}

//
func calcularNumeroDeEstructuras(tamanioPart uint32) (cantidad uint32) {

	dos := str.TamAVD
	tres := str.TamDetalleDirect
	cuatro := 5 * str.TamInodo
	cinco := 20 * str.TamBloque
	seis := str.TamBitacora

	numerador := tamanioPart - uint32(2*str.TamSuperBoot)
	denominador := 27 + dos + tres + (cuatro + (cinco) + seis)

	fmt.Println("numerador", numerador)
	fmt.Println("denominador", denominador)

	return numerador / uint32(denominador)
}

//Retorna una particion montada
func getParticionByID(idParticion string) (part str.ParticionMontada) {
	letra := idParticion[2]
	numero := idParticion[3:len(idParticion)]
	codigoLetra, err := strconv.ParseUint(numero, 10, 16)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(str.ParticionesMontadas); i++ {
		if str.ParticionesMontadas[i].Letra == byte(letra) && str.ParticionesMontadas[i].Numero == uint16(codigoLetra) {
			part = str.ParticionesMontadas[i]
		}
	}
	return part
}

func getParticion(mbr str.MBR, n int) (selecPart str.Particion) {

	n++
	switch n {
	case 1:
		return mbr.Part1

	case 2:
		return mbr.Part2

	case 3:
		return mbr.Part3

	case 4:
		return mbr.Part4
	}

	return
}

func consultarEspacioDisponible(mbr str.MBR, ini int) (espacioDisp uint32, inicioPart uint32) {

	limite := mbr.Tamanio
	inicio := uint32(138)

	for i := ini; i > -1; i-- {
		if part := getParticion(mbr, i); part.Estado == 0 {
			inicio = part.Inicio + part.Tamanio
			break
		}
	}

	for i := ini + 1; i < 4; i++ {
		if part := getParticion(mbr, i); part.Estado == 0 {
			limite = part.Inicio
			break
		}
	}

	espacioDisp = limite - inicio

	return espacioDisp, inicio
}

func convertirBinario(numeros []byte) (numerosTraducidos uint32) {

	n := len(numeros)
	numeroBinario := ""

	for i := 0; i < n; i++ {
		temp := fmt.Sprintf("%b", numeros[i])

		nTemp := len(temp)
		for j := 0; j < 8-nTemp; j++ {
			temp = "0" + temp
			//fmt.Println(i,".- ", temp)
		}
		numeroBinario = numeroBinario + temp
		//fmt.Println(numeroBinario)
	}

	numeroTraducido, _ := strconv.ParseUint(numeroBinario, 2, 32)

	//fmt.Println("BINARIO", numeroBinario, "DECIMAL", numeroTraducido)
	return uint32(numeroTraducido)
}

//
func CrearGraficaMBR(contenidoDisco []byte) {

	//mbr := MontarMBR(contenidoDisco)

	file, err := os.Create("/home/helmut/Escritorio/graficaMBR.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	graph := ""

	errr := ioutil.WriteFile("/home/helmut/Escritorio/graficaMBR.txt", []byte(graph), 0666)
	if errr != nil {
		panic(errr)
	}

}
