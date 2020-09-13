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

func editarArchivoVariasVeces(url string, cadena [][]byte, inicio int64) {
	// Read Write Mode
	file, err := os.OpenFile(url, os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	for i := 0; i < len(cadena); i++ {
		_, err := file.WriteAt(cadena[i], inicio) // Write at 0 beginning
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
	}

}

func editarArchivoVariasCadenas(url string, cadena []byte, inicio int64, n int) {
	// Read Write Mode
	file, err := os.OpenFile(url, os.O_RDWR, 0644)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	for i := 0; i < n; i++ {
		len, err := file.WriteAt(cadena, inicio) // Write at 0 beginning
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
		_ = len
	}
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

		nombre := part[11:27]
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
			goto final
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

		contenidoParticion := contenidoDisco[part.Inicio : part.Inicio+part.Tamanio]

		if len(str.ParticionesMontadas) == 0 {
			inicio := int(part.Inicio)
			fin := int(part.Inicio + part.Tamanio)
			superBoot := MontarSuperBoot(contenidoDisco[inicio:fin])
			avd, bitmap := montarAVD(superBoot, contenidoParticion)
			detalleD, bitMapDD := montarDetalleD(superBoot, contenidoParticion)
			bloquesM, bitMapBloque := montarBloques(superBoot, contenidoParticion)
			inodos, bitMapInodo := montarInodo(superBoot, contenidoParticion)

			fmt.Println("BLOQUE DE CADA AVD:")
			fmt.Println(avd)
			fmt.Println("BIT MAP DE LOS BLOQUES AVD")
			fmt.Println(bitmap)

			bloques := str.Bloques{BitMapAVD: bitmap, AVD: avd, BitMapDetalleD: bitMapDD, DetalleD: detalleD, Bloque: bloquesM, BitMapBloques: bitMapBloque, Inodo: inodos, BitMapInodo: bitMapInodo}
			partMont := str.ParticionMontada{Particion: part, ContenidoParticion: contenidoParticion, Letra: 97, Numero: 1, Ruta: url, Superboot: superBoot, Bloques: bloques}
			fmt.Println("****************************************************PARTICION MONTADA *****************************************************")
			fmt.Println(partMont)
			fmt.Println("***************************************************************************************************************************")
			str.ParticionesMontadas = append(str.ParticionesMontadas, partMont)

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

			inicio := int(part.Inicio)
			fin := int(part.Inicio + part.Tamanio)
			superBoot := MontarSuperBoot(contenidoDisco[inicio:fin])
			partMont := str.ParticionMontada{Particion: part, ContenidoParticion: contenidoParticion, Letra: codigo, Numero: uint16(n + 1), Ruta: url, Superboot: superBoot}
			fmt.Println("****************************************************PARTICION MONTADA *****************************************************")
			fmt.Println(partMont)
			fmt.Println("***************************************************************************************************************************")
			str.ParticionesMontadas = append(str.ParticionesMontadas, partMont)
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

	if part.Ruta != "" {

		if tipo == "fast" {

		} else if tipo == "full" {
			var cadena []byte
			for i := 0; i < int(part.Particion.Tamanio); i++ {
				cadena = append(cadena, 0)
			}

			rutaDisco := s.Split(part.Ruta, "/")
			editarArchivo(part.Ruta, cadena, int64(part.Particion.Inicio))
			superBoot := CrearSuperBoot(part.Particion.Tamanio, rutaDisco[len(rutaDisco)-1])
			part.Superboot = superBoot
			var buffer bytes.Buffer
			binary.Write(&buffer, binary.BigEndian, superBoot)
			editarArchivo(part.Ruta, buffer.Bytes(), int64(part.Particion.Inicio))
			editarArchivo(part.Ruta, buffer.Bytes(), int64(superBoot.ApuntadorBitacora+superBoot.CantidadAVD*superBoot.TamanioBitacora))
		}

	} else {

		fmt.Println("")
		fmt.Println("*************************************************************")
		fmt.Println("*                          ALERTA                           *")
		fmt.Println("*************************************************************")
		fmt.Println("*         NO HAY UNA PARTICION MONTADA CON ESE ID           *")
		fmt.Println("*************************************************************")
		fmt.Println("")

	}

}

//
func CrearSuperBoot(tamanioPart uint32, nombre string) str.SuperBoot {

	nEstructuras := calcularNumeroDeEstructuras(tamanioPart)

	var nombreDisco [20]byte
	for i := 0; i < len(nombre); i++ {
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
		ApuntadorBitMapAVD:           apuntadorBitMapAVD,
		ApuntadorAVD:                 apuntadorAVD,
		ApuntadorBitMapDetalleDirect: apuntadorBitMapDetalleDirect,
		ApuntadorDetalleDirect:       apuntadorDetalleDirect,
		ApuntadorBitMapInodos:        apuntadorBitMapInodos,
		ApuntadorInodos:              apuntadorInodos,
		ApuntadorBitMapBloques:       apuntadorBitMapBloques,
		ApuntadorBloques:             apuntadorBloques,
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

//
func MontarSuperBoot(contenidoParticion []byte) str.SuperBoot {

	nombre := contenidoParticion[0:20]
	var nombreDisco [20]byte
	for i := 0; i < len(nombre); i++ {
		nombreDisco[i] = nombre[i]
	}
	cantidadAVD := convertirBinario(contenidoParticion[20:24])
	cantidadDetalleDirect := convertirBinario(contenidoParticion[24:28])
	cantidadInodos := convertirBinario(contenidoParticion[28:32])
	cantidadBloques := convertirBinario(contenidoParticion[32:36])
	cantidadAVDLibres := convertirBinario(contenidoParticion[36:40])
	cantidadDetalleDirecttLibres := convertirBinario(contenidoParticion[40:44])
	cantidadInodosLibres := convertirBinario(contenidoParticion[44:48])
	cantidadBloquesLibres := convertirBinario(contenidoParticion[48:52])

	fechaC := contenidoParticion[52:74]
	var fechaCreacion [22]byte
	for i := 0; i < len(fechaC); i++ {
		fechaCreacion[i] = fechaC[i]
	}
	fechaUM := contenidoParticion[74:96]
	var fechaUltimoMontaje [22]byte
	for i := 0; i < len(fechaUM); i++ {
		fechaUltimoMontaje[i] = fechaUM[i]
	}
	//16
	numeroMontajes := uint16(convertirBinario(contenidoParticion[96:98]))

	apuntadorBitMapAVD := convertirBinario(contenidoParticion[98:102])
	apuntadorAVD := convertirBinario(contenidoParticion[102:106])
	apuntadorBitMapDetalleDirect := convertirBinario(contenidoParticion[106:110])
	apuntadorDetalleDirect := convertirBinario(contenidoParticion[110:114])
	apuntadorBitMapInodos := convertirBinario(contenidoParticion[114:118])
	apuntadorInodos := convertirBinario(contenidoParticion[118:122])
	apuntadorBitMapBloques := convertirBinario(contenidoParticion[122:126])
	apuntadorBloques := convertirBinario(contenidoParticion[126:130])
	apuntadorBitacora := convertirBinario(contenidoParticion[130:134])

	tamanioAVD := convertirBinario(contenidoParticion[134:138])
	tamanioDetalleDirect := convertirBinario(contenidoParticion[138:142])
	tamanioInodo := convertirBinario(contenidoParticion[142:146])
	tamanioBloque := convertirBinario(contenidoParticion[146:150])
	tamanioBitacora := convertirBinario(contenidoParticion[150:154])

	primerAVDLibre := convertirBinario(contenidoParticion[154:158])
	primerDetalleDirectLibre := convertirBinario(contenidoParticion[158:162])
	primerInodoLibre := convertirBinario(contenidoParticion[162:166])
	primerBloqueLibre := convertirBinario(contenidoParticion[166:170])
	numeroMagico := convertirBinario(contenidoParticion[170:174])

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
		ApuntadorBitMapAVD:           apuntadorBitMapAVD,
		ApuntadorAVD:                 apuntadorAVD,
		ApuntadorBitMapDetalleDirect: apuntadorBitMapDetalleDirect,
		ApuntadorDetalleDirect:       apuntadorDetalleDirect,
		ApuntadorBitMapInodos:        apuntadorBitMapInodos,
		ApuntadorInodos:              apuntadorInodos,
		ApuntadorBitMapBloques:       apuntadorBitMapBloques,
		ApuntadorBloques:             apuntadorBloques,
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
		NumeroMagico:                 numeroMagico}

	return superboot

}

//
func CrearRoot(idPart string, idProp uint32, idGrupo uint32, permisos uint16) {

	paso := true

	part := getParticionByID(idPart)

	fechaCreacion := generarFecha()
	var nombreDirectorio [20]byte
	nombreDirectorio[0] = '/'
	subAVD := [6]uint32{0, 0, 0, 0, 0, 0}
	apuntadorDetalleDirect := uint32(0)
	apuntadorIndirecto := uint32(0)

	inicioParticion := part.Particion.Inicio
	inicioAVD := part.Superboot.ApuntadorAVD
	inicioBitMap := part.Superboot.ApuntadorBitMapAVD
	finalBloquesAVD := inicioAVD + uint32(str.TamAVD)*part.Superboot.CantidadAVD

	for i := inicioAVD + 22; i < finalBloquesAVD; i += uint32(str.TamAVD) {
		var nombreTemp [20]byte
		for j := 0; j < 20; j++ {
			nombreTemp[j] = part.ContenidoParticion[i+uint32(j)]
		}
		if nombreTemp == nombreDirectorio {
			paso = false
			break
		}
	}

	if paso {
		avd := str.AVD{FechaCreacion: fechaCreacion,
			NombreDirectorio:       nombreDirectorio,
			SubAVD:                 subAVD,
			ApuntadorDetalleDirect: apuntadorDetalleDirect,
			ApuntadorIndirecto:     apuntadorIndirecto,
			IDPropietario:          idProp,
			IDGrupo:                idGrupo,
			Permisos:               permisos}

		var buffer bytes.Buffer
		binary.Write(&buffer, binary.BigEndian, avd)
		cadenaAVD := buffer.Bytes()

		cadena, bitLibre := editarBitMapYBloque(part.ContenidoParticion[inicioBitMap:finalBloquesAVD], part.Superboot.CantidadAVD, part.Superboot.TamanioAVD, cadenaAVD)
		part.Superboot.CantidadAVDLibres--
		part.Superboot.PrimerAVDLibre = bitLibre + part.Superboot.ApuntadorBitMapAVD

		buffer.Reset()
		binary.Write(&buffer, binary.BigEndian, part.Superboot)
		//EDITANDO EL SUPERBOOT
		editarArchivo(part.Ruta, buffer.Bytes(), int64(inicioParticion))
		//EDITANDO BITMAP Y BLOQUE DONDE 138 ES EL FINAL DEL MBR
		editarArchivo(part.Ruta, cadena, int64(inicioBitMap)+138)

		fmt.Println("*************************************************************")
		fmt.Println("*                 CARPETA ROOT CREADA                       *")
		fmt.Println("*************************************************************")

	} else {

		fmt.Println("*************************************************************")
		fmt.Println("*                        ALERTA                             *")
		fmt.Println("*************************************************************")
		fmt.Println("*                EL DIRECCTORIO YA EXISTE                   *")
		fmt.Println("*************************************************************")

	}
}

//
func CrearAVDInicio(idPart string, ruta []string, idProp uint32, idGrupo uint32, permisos uint16, contenido string) {

	nombresDirect := [][20]byte{}

	for i := 0; i < len(ruta); i++ {
		nombresDirect = append(nombresDirect, convertirNombreASlice(convertirStringASlice(ruta[i])))
	}

	part := getParticionByID(idPart)

	apuntador, ultimo, indice := buscarAVD(part.Bloques.AVD, nombresDirect, 0, part.Superboot.ApuntadorAVD)

	CrearArchivo(part, contenido)

	avdNuevos := []uint32{}
	part, avdNuevos = crearAVD(part, ultimo, indice, nombresDirect, idProp, idGrupo, permisos, avdNuevos)

	fmt.Println("APUNTADOR:", apuntador)
	fmt.Println("ULTIMO ENCONTRADO:", ultimo)
	fmt.Println("INDICE:", indice)

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, part.Superboot)
	editarArchivo(part.Ruta, buffer.Bytes(), int64(part.Particion.Inicio))
	buffer.Reset()

	binary.Write(&buffer, binary.BigEndian, part.Bloques.BitMapAVD)
	editarArchivo(part.Ruta, buffer.Bytes(), int64(part.Superboot.ApuntadorBitMapAVD+part.Particion.Inicio))
	buffer.Reset()

	for i := 0; i < len(part.Bloques.AVD); i++ {
		for j := 0; j < len(avdNuevos); j++ {

			if avdNuevos[j] == part.Bloques.AVD[i].Apuntador {
				binary.Write(&buffer, binary.BigEndian, part.Bloques.AVD[i].Avd)
				editarArchivo(part.Ruta, buffer.Bytes(), int64(part.Bloques.AVD[i].Apuntador+part.Particion.Inicio))
				buffer.Reset()
			}

		}

	}

}

func crearAVD(part str.ParticionMontada, ultimo uint32, indice int, nombres [][20]byte, idProp uint32, idGrupo uint32, permisos uint16, avdM []uint32) (str.ParticionMontada, []uint32) {

	if indice+1 <= len(nombres) {

		n := 0
		for i := 0; i < len(part.Bloques.BitMapAVD); i++ {
			if part.Bloques.BitMapAVD[i] == 0 {
				part.Bloques.BitMapAVD[i] = 1
				part.Superboot.CantidadAVDLibres--
				for j := i; j < len(part.Bloques.BitMapAVD); j++ {
					if part.Bloques.BitMapAVD[j] == 0 {
						part.Superboot.PrimerAVDLibre = part.Superboot.ApuntadorBitMapAVD + uint32(j)
						break
					}
				}
				break
			}
			n++
		}

		lleno := true
		for i := 0; i < len(part.Bloques.AVD); i++ {
			if part.Bloques.AVD[i].Apuntador == ultimo {
				for j := 0; j < len(part.Bloques.AVD[i].Avd.SubAVD); j++ {
					if part.Bloques.AVD[i].Avd.SubAVD[j] == 0 {
						nuevoApuntador := part.Superboot.ApuntadorAVD + uint32(n)*uint32(str.TamAVD)
						part.Bloques.AVD[i].Avd.SubAVD[j] = nuevoApuntador
						lleno = false

						fechaCreacion := generarFecha()
						subAVD := [6]uint32{0, 0, 0, 0, 0, 0}

						avd := str.AVD{FechaCreacion: fechaCreacion,
							NombreDirectorio:       nombres[indice],
							SubAVD:                 subAVD,
							ApuntadorDetalleDirect: 0,
							ApuntadorIndirecto:     0,
							IDPropietario:          idProp,
							IDGrupo:                idGrupo,
							Permisos:               permisos}

						avdMontada := str.AVD_Montada{Apuntador: nuevoApuntador, Avd: avd}
						for k := 0; k < len(part.Bloques.AVD); k++ {
							if part.Bloques.AVD[k].Apuntador == nuevoApuntador {
								part.Bloques.AVD[k] = avdMontada
								break
							}
						}

						avdM = append(avdM, ultimo)
						avdM = append(avdM, nuevoApuntador)
						part, avdM = crearAVD(part, nuevoApuntador, indice+1, nombres, idProp, idGrupo, permisos, avdM)
						goto fin
					}
				}
				break
			}

		}

		if lleno {

			nuevoApuntador := part.Superboot.ApuntadorBitMapAVD + uint32(n)*uint32(str.TamAVD)
			ultimo = nuevoApuntador
			apuntadorD := uint32(0)
			idPropietario := uint32(0)
			idGrupoo := uint32(0)
			permisoss := uint16(0)

			for i := 0; i < len(part.Bloques.AVD); i++ {
				if part.Bloques.AVD[i].Apuntador == ultimo {
					part.Bloques.AVD[i].Avd.ApuntadorIndirecto = nuevoApuntador
					apuntadorD = part.Bloques.AVD[i].Avd.ApuntadorDetalleDirect
					idPropietario = part.Bloques.AVD[i].Avd.IDPropietario
					idGrupoo = part.Bloques.AVD[i].Avd.IDGrupo
					permisoss = part.Bloques.AVD[i].Avd.Permisos

					break
				}
			}

			fechaCreacion := generarFecha()
			subAVD := [6]uint32{0, 0, 0, 0, 0, 0}

			avd := str.AVD{FechaCreacion: fechaCreacion,
				NombreDirectorio:       nombres[indice-1],
				SubAVD:                 subAVD,
				ApuntadorDetalleDirect: apuntadorD,
				ApuntadorIndirecto:     0,
				IDPropietario:          idPropietario,
				IDGrupo:                idGrupoo,
				Permisos:               permisoss}

			avdMontada := str.AVD_Montada{Apuntador: nuevoApuntador, Avd: avd}
			for i := 0; i < len(part.Bloques.AVD); i++ {
				if part.Bloques.AVD[i].Apuntador == nuevoApuntador {
					part.Bloques.AVD[i] = avdMontada
					break
				}
			}

			avdM = append(avdM, ultimo)
			avdM = append(avdM, nuevoApuntador)
			part, avdM = crearAVD(part, ultimo, indice+1, nombres, idProp, idGrupo, permisos, avdM)
			goto fin

		}

	}

fin:
	return part, avdM
}

//
func CrearDetalleD() {


	
}

/*
func crearDetalle(part str.ParticionMontada, ultimo uint32, avdAptr uint32, nombre [20]byte, idProp uint32, idGrupo uint32, permisos uint16, detalleD []uint32) {

	n := 0
	for i := 0; i < len(part.Bloques.BitMapDetalleD); i++ {
		if part.Bloques.BitMapDetalleD[i] == 0 {
			part.Bloques.BitMapDetalleD[i] = 1
			part.Superboot.CantidadDetalleDirectLibres--
			for j := i; j < len(part.Bloques.BitMapDetalleD); j++ {
				if part.Bloques.BitMapDetalleD[j] == 0 {
					part.Superboot.PrimerDetalleDirectLibre = part.Superboot.ApuntadorBitMapDetalleDirect + uint32(j)
					break
				}
			}
			break
		}
		n++
	}

	lleno := true

	for i := 0; i < len(part.Bloques.AVD); i++ {

		if part.Bloques.AVD[i].Apuntador == avdAptr {

			nuevoApuntador := part.Superboot.ApuntadorDetalleDirect + uint32(n)*uint32(str.TamDetalleDirect)
			for j := 0; j < len(part.Bloques.DetalleD); j++ {
				if part.Bloques.DetalleD[j].Apuntador == nuevoApuntador {

					dD := part.Bloques.DetalleD[j].DetalleD.Archivos
					for k := 0; k < len(dD); k++ {

						temp := [20]byte{}
						if dD[k].Nombre == temp {

							lleno = false
						}

					}
					break
				}
			}

			break
		}

	}

	for i := 0; i < len(part.Bloques.DetalleD); i++ {
		if part.Bloques.DetalleD[i].Apuntador == ultimo {

			for j := 0; j < len(part.Bloques.AVD[i].Avd.SubAVD); j++ {
				if part.Bloques.AVD[i].Avd.SubAVD[j] == 0 {
					nuevoApuntador := part.Superboot.ApuntadorAVD + uint32(n)*uint32(str.TamAVD)
					part.Bloques.AVD[i].Avd.SubAVD[j] = nuevoApuntador
					lleno = false

					fechaCreacion := generarFecha()
					subAVD := [6]uint32{0, 0, 0, 0, 0, 0}

					avd := str.AVD{FechaCreacion: fechaCreacion,
						NombreDirectorio:       nombres[indice],
						SubAVD:                 subAVD,
						ApuntadorDetalleDirect: 0,
						ApuntadorIndirecto:     0,
						IDPropietario:          idProp,
						IDGrupo:                idGrupo,
						Permisos:               permisos}

					avdMontada := str.AVD_Montada{Apuntador: nuevoApuntador, Avd: avd}
					for k := 0; k < len(part.Bloques.AVD); k++ {
						if part.Bloques.AVD[k].Apuntador == nuevoApuntador {
							part.Bloques.AVD[k] = avdMontada
							break
						}
					}

					avdM = append(avdM, ultimo)
					avdM = append(avdM, nuevoApuntador)
					part, avdM = crearAVD(part, nuevoApuntador, indice+1, nombres, idProp, idGrupo, permisos, avdM)
					goto fin
				}
			}
			break
		}

	}

	if lleno {

		nuevoApuntador := part.Superboot.ApuntadorBitMapAVD + uint32(n)*uint32(str.TamAVD)
		ultimo = nuevoApuntador
		apuntadorD := uint32(0)
		idPropietario := uint32(0)
		idGrupoo := uint32(0)
		permisoss := uint16(0)

		for i := 0; i < len(part.Bloques.AVD); i++ {
			if part.Bloques.AVD[i].Apuntador == ultimo {
				part.Bloques.AVD[i].Avd.ApuntadorIndirecto = nuevoApuntador
				apuntadorD = part.Bloques.AVD[i].Avd.ApuntadorDetalleDirect
				idPropietario = part.Bloques.AVD[i].Avd.IDPropietario
				idGrupoo = part.Bloques.AVD[i].Avd.IDGrupo
				permisoss = part.Bloques.AVD[i].Avd.Permisos

				break
			}
		}

		fechaCreacion := generarFecha()
		subAVD := [6]uint32{0, 0, 0, 0, 0, 0}

		avd := str.AVD{FechaCreacion: fechaCreacion,
			NombreDirectorio:       nombres[indice-1],
			SubAVD:                 subAVD,
			ApuntadorDetalleDirect: apuntadorD,
			ApuntadorIndirecto:     0,
			IDPropietario:          idPropietario,
			IDGrupo:                idGrupoo,
			Permisos:               permisoss}

		avdMontada := str.AVD_Montada{Apuntador: nuevoApuntador, Avd: avd}
		for i := 0; i < len(part.Bloques.AVD); i++ {
			if part.Bloques.AVD[i].Apuntador == nuevoApuntador {
				part.Bloques.AVD[i] = avdMontada
				break
			}
		}

		avdM = append(avdM, ultimo)
		avdM = append(avdM, nuevoApuntador)
		part, avdM = crearAVD(part, ultimo, indice+1, nombres, idProp, idGrupo, permisos, avdM)
		goto fin

	}

fin:
}

*/

func convertirStringADatos(datos []byte) (salida [25]byte) {

	for i := 0; i < len(datos) && i < 25; i++ {
		salida[i] = datos[i]
	}

	return salida
}

//
func CrearArchivo(part str.ParticionMontada, contenido string) {

	cadena := convertirStringASlice(contenido)
	bloquesDatos := [][25]byte{}

	for i := 0; i < len(cadena); i++ {

		if len(cadena) >= 25 {
			bloquesDatos = append(bloquesDatos, convertirStringADatos(cadena[0:25]))
			cadena = cadena[25:]

		} else {
			bloquesDatos = append(bloquesDatos, convertirStringADatos(cadena[0:]))
			break
		}
	}

	//CREAR BLOQUES

	apuntadoresBloques := []uint32{}
	temp := uint32(0)
	for i := 0; i < len(bloquesDatos); i++ {
		part, temp = crearBloques(part, bloquesDatos[i])
		apuntadoresBloques = append(apuntadoresBloques, temp)
	}

	//CREAR INODO

	InodoNuevos := []uint32{}
	part, InodoNuevos = crearInodo(part, 0, apuntadoresBloques, 0, uint16(len(apuntadoresBloques)), uint32(len(contenido)), 2, 3, 77, InodoNuevos)

	//CREAR DETALLE DE DIRECTORIO
}

func crearInodo(part str.ParticionMontada, ultimo uint32, bloques []uint32, indice int, nBloques uint16, tamanio uint32, idP uint32, idG uint32, permisos uint16, inodos []uint32) (str.ParticionMontada, []uint32) {

	n := 0
	for i := 0; i < len(part.Bloques.BitMapInodo); i++ {
		if part.Bloques.BitMapInodo[i] == 0 {
			part.Bloques.BitMapInodo[i] = 1
			part.Superboot.CantidadInodosLibres--
			for j := i; j < len(part.Bloques.BitMapInodo); j++ {
				if part.Bloques.BitMapInodo[j] == 0 {
					part.Superboot.PrimerInodoLibre = part.Superboot.ApuntadorInodos + uint32(j)
					break
				}
			}
			break
		}
		n++
	}

	lleno := true
	for i := 0; i < len(part.Bloques.Inodo); i++ {
		if part.Bloques.Inodo[i].Apuntador == part.Superboot.ApuntadorInodos+uint32(n)*uint32(str.TamInodo) {

			bloquesInodo := [4]uint32{}
			for j := 0; j < len(bloquesInodo) && indice < len(bloques); j++ {
				if bloquesInodo[j] == 0 {
					bloquesInodo[j] = bloques[indice]
					indice++
				}
			}
			if ultimo!=0 {
				for j := 0; j < len(part.Bloques.Inodo); j++ {
					if part.Bloques.Inodo[j].Apuntador == ultimo {
						part.Bloques.Inodo[j].Inodo.ApuntadorIndirecto = part.Bloques.Inodo[i].Apuntador
						break
					}
				}
			}
			if indice < len(bloques) {
				lleno = true
				ultimo = part.Bloques.Inodo[i].Apuntador
			} else {
				lleno = false
			}
			
			inodo := str.Inodo{
				Numero:             uint32(n),
				TamanioArchivo:     tamanio,
				NumeroBloques:      nBloques,
				Bloques:            bloquesInodo,
				ApuntadorIndirecto: 0,
				IDPropietario:      idP,
				IDGrupo:            idG,
				Permisos:           permisos}

			part.Bloques.Inodo[i].Inodo = inodo
			inodos = append(inodos, part.Bloques.Inodo[i].Apuntador)
			break
		}

	}

	if lleno {
		part, inodos = crearInodo(part, ultimo, bloques, indice, nBloques, tamanio, idP, idG, permisos, inodos)
	}

	return part, inodos
}

func crearBloques(part str.ParticionMontada, contenido [25]byte) (str.ParticionMontada, uint32) {

	n := 0
	for i := 0; i < len(part.Bloques.BitMapBloques); i++ {
		if part.Bloques.BitMapBloques[i] == 0 {
			part.Bloques.BitMapBloques[i] = 1
			part.Superboot.CantidadBloquesLibres--
			for j := i; j < len(part.Bloques.BitMapBloques); j++ {
				if part.Bloques.BitMapBloques[j] == 0 {
					part.Superboot.PrimerBloqueLibre = part.Superboot.ApuntadorBloques + uint32(j)
					break
				}
			}
			break
		}
		n++
	}

	nuevoApuntador := part.Superboot.ApuntadorBloques + uint32(n)*uint32(str.TamBloque)
	for i := 0; i < len(part.Bloques.Bloque); i++ {
		if part.Bloques.Bloque[i].Apuntador == nuevoApuntador {
			bTemp := str.Bloque{Datos: contenido}
			part.Bloques.Bloque[i].Bloque = bTemp
			break
		}

	}

	return part, nuevoApuntador
}

func buscarAVD(avd []str.AVD_Montada, nombres [][20]byte, indice int, apuntador uint32) (apuntadorAVD uint32, ultimoReconocido uint32, indiceNombre int) {

	indiceNombre = indice
	ultimoReconocido = apuntador
	apuntadorAVD = 0

	if len(nombres) >= indice+1 {

		for i := 0; i < len(avd); i++ {

			if avd[i].Apuntador == apuntador {
				if nombres[indice] == avd[i].Avd.NombreDirectorio {

					if len(nombres) == indice+1 {
						apuntadorAVD = apuntador
						goto fin

					} else {
						apuntadorAVD, ultimoReconocido, indiceNombre = buscarAVD(avd, nombres, indice+1, apuntador)
						goto fin
					}

				}
				break
			}

		}

		for i := 0; i < len(avd); i++ {

			if avd[i].Apuntador == apuntador {

				for j := 0; j < len(avd[i].Avd.SubAVD); j++ {

					if avd[i].Avd.SubAVD[j] != 0 {
						if temp := getNombreAVDByApuntador(avd, avd[i].Avd.SubAVD[j]); temp == nombres[indice] {
							apuntadorAVD, ultimoReconocido, indiceNombre = buscarAVD(avd, nombres, indice+1, avd[i].Avd.SubAVD[j])
							break
						}
					}
				}
				break
			}
		}

		if apuntadorAVD == 0 {

			for i := 0; i < len(avd); i++ {
				if avd[i].Apuntador == apuntador {
					if avd[i].Avd.ApuntadorIndirecto != 0 {
						apuntadorAVD, ultimoReconocido, indiceNombre = buscarAVD(avd, nombres, indice+1, apuntador)
					}
					break
				}
			}

		}

	}

fin:
	return apuntadorAVD, ultimoReconocido, indiceNombre

}

func getNombreAVDByApuntador(avd []str.AVD_Montada, apuntador uint32) (nombreDirectorio [20]byte) {

	for i := 0; i < len(avd); i++ {
		if avd[i].Apuntador == apuntador {
			nombreDirectorio = avd[i].Avd.NombreDirectorio
			break
		}
	}

	return nombreDirectorio
}

func montarAVD(sb str.SuperBoot, contenidoParticion []byte) (bloqueAVD []str.AVD_Montada, BitMapAVD []byte) {

	inicio := sb.ApuntadorAVD
	fin := inicio + sb.CantidadAVD*uint32(str.TamAVD)
	contenidoBloque := contenidoParticion[inicio:fin]

	BitMapAVD = contenidoParticion[sb.ApuntadorBitMapAVD : sb.ApuntadorBitMapAVD+sb.CantidadAVD]

	for i := 0; i < len(contenidoBloque); i += int(str.TamAVD) {

		fecha := convertirFechaASlice(contenidoBloque[i : i+22])
		nombre := convertirNombreASlice(contenidoBloque[i+22 : i+42])
		subAVD := [6]uint32{}
		subAVD[0] = convertirBinario(contenidoBloque[i+42 : i+46])
		subAVD[1] = convertirBinario(contenidoBloque[i+46 : i+50])
		subAVD[2] = convertirBinario(contenidoBloque[i+50 : i+54])
		subAVD[3] = convertirBinario(contenidoBloque[i+54 : i+58])
		subAVD[4] = convertirBinario(contenidoBloque[i+58 : i+62])
		subAVD[5] = convertirBinario(contenidoBloque[i+62 : i+66])
		apuntadorDetalleD := convertirBinario(contenidoBloque[i+66 : i+70])
		apuntadorIndirecto := convertirBinario(contenidoBloque[i+70 : i+74])
		idProp := convertirBinario(contenidoBloque[i+74 : i+78])
		idGrupo := convertirBinario(contenidoBloque[i+78 : i+82])
		permisos := uint16(convertirBinario(contenidoBloque[i+82 : i+84]))

		avd := str.AVD{FechaCreacion: fecha,
			NombreDirectorio:       nombre,
			SubAVD:                 subAVD,
			ApuntadorDetalleDirect: apuntadorDetalleD,
			ApuntadorIndirecto:     apuntadorIndirecto,
			IDPropietario:          idProp,
			IDGrupo:                idGrupo,
			Permisos:               permisos}

		avdMontada := str.AVD_Montada{Apuntador: uint32(i) + inicio, Avd: avd}
		bloqueAVD = append(bloqueAVD, avdMontada)
	}

	return bloqueAVD, BitMapAVD
}

func montarDetalleD(sb str.SuperBoot, contenidoParticion []byte) (bloqueDD []str.DetalleDirectorio_Montado, BitMapDetalleD []byte) {

	inicio := sb.ApuntadorDetalleDirect
	fin := inicio + sb.CantidadDetalleDirect*uint32(str.TamDetalleDirect)
	contenidoBloque := contenidoParticion[inicio:fin]

	BitMapDetalleD = contenidoParticion[sb.ApuntadorBitMapDetalleDirect : sb.ApuntadorBitMapDetalleDirect+sb.CantidadDetalleDirect]

	for i := 0; i < len(contenidoBloque); i += int(str.TamDetalleDirect) {

		temp := []str.InfoArchivo{}
		for j := uint32(i); j < uint32(5)*uint32(str.TamInfoArchivo)+uint32(i); j += uint32(str.TamInfoArchivo) {
			nombre := convertirNombreASlice(contenidoBloque[j : j+20])
			apuntadorInodo := convertirBinario(contenidoBloque[j+20 : j+24])
			fechaC := convertirFechaASlice(contenidoBloque[j+24 : j+46])
			fechaM := convertirFechaASlice(contenidoBloque[j+46 : j+68])

			archivo := str.InfoArchivo{Nombre: nombre, ApuntadorInodo: apuntadorInodo, FechaCreacion: fechaC, FechaModificacion: fechaM}
			temp = append(temp, archivo)
		}

		archivos := [5]str.InfoArchivo{}
		for j := 0; j < 5; j++ {
			archivos[j] = temp[j]
		}

		numero := 5*uint32(str.TamInfoArchivo) + uint32(i)
		apuntadorIndirecto := convertirBinario(contenidoBloque[numero : numero+4])

		detalleD := str.DetalleDirectorio{Archivos: archivos, ApuntadorIndirecto: apuntadorIndirecto}

		dDMontado := str.DetalleDirectorio_Montado{Apuntador: uint32(i) + inicio, DetalleD: detalleD}
		bloqueDD = append(bloqueDD, dDMontado)
	}

	return bloqueDD, BitMapDetalleD
}

func montarInodo(sb str.SuperBoot, contenidoParticion []byte) (bloqueInodo []str.Inodo_Montado, BitMapInodo []byte) {

	inicio := sb.ApuntadorInodos
	fin := inicio + sb.CantidadInodos*uint32(str.TamInodo)
	contenidoBloque := contenidoParticion[inicio:fin]

	BitMapInodo = contenidoParticion[sb.ApuntadorBitMapInodos : sb.ApuntadorBitMapInodos+sb.CantidadInodos]

	for i := 0; i < len(contenidoBloque); i += int(str.TamInodo) {

		numero := convertirBinario(contenidoBloque[i : i+4])
		tamanio := convertirBinario(contenidoBloque[i+4 : i+8])
		numeroBloques := uint16(convertirBinario(contenidoBloque[i+8 : i+10]))
		bloques := [4]uint32{}
		bloques[0] = convertirBinario(contenidoBloque[i+10 : i+14])
		bloques[1] = convertirBinario(contenidoBloque[i+14 : i+18])
		bloques[2] = convertirBinario(contenidoBloque[i+18 : i+22])
		bloques[3] = convertirBinario(contenidoBloque[i+22 : i+26])
		indirecto := convertirBinario(contenidoBloque[i+26 : i+30])
		idProp := convertirBinario(contenidoBloque[i+30 : i+34])
		idGrupo := convertirBinario(contenidoBloque[i+34 : i+38])
		permisos := uint16(convertirBinario(contenidoBloque[i+38 : i+40]))

		inodo := str.Inodo{Numero: numero,
			TamanioArchivo:     tamanio,
			NumeroBloques:      numeroBloques,
			Bloques:            bloques,
			ApuntadorIndirecto: indirecto,
			IDPropietario:      idProp,
			IDGrupo:            idGrupo,
			Permisos:           permisos}

		inodoMontado := str.Inodo_Montado{Apuntador: uint32(i) + inicio, Inodo: inodo}
		bloqueInodo = append(bloqueInodo, inodoMontado)
	}

	return bloqueInodo, BitMapInodo
}

func montarBloques(sb str.SuperBoot, contenidoParticion []byte) (bloques []str.Bloque_Montado, BitMapBloque []byte) {

	inicio := sb.ApuntadorBloques
	fin := inicio + sb.CantidadBloques*uint32(str.TamBloque)
	contenidoBloque := contenidoParticion[inicio:fin]

	BitMapBloque = contenidoParticion[sb.ApuntadorBitMapBloques : sb.ApuntadorBitMapBloques+sb.CantidadBloques]

	for i := 0; i < len(contenidoBloque); i += int(str.TamBloque) {

		datos := convertirStringADatos(contenidoBloque[i : i+25])

		bloque := str.Bloque{Datos: datos}

		bloqueM := str.Bloque_Montado{Apuntador: uint32(i) + inicio, Bloque: bloque}
		bloques = append(bloques, bloqueM)
	}

	return bloques, BitMapBloque
}

func convertirNombreASlice(nombre []byte) (nombreSalida [20]byte) {

	for i := 0; i < len(nombre) && i < 20; i++ {
		nombreSalida[i] = nombre[i]
	}

	return nombreSalida
}

func convertirFechaASlice(fecha []byte) (fechaSalida [22]byte) {

	for i := 0; i < len(fecha) && i < 22; i++ {
		fechaSalida[i] = fecha[i]
	}

	return fechaSalida
}

func convertirStringASlice(nombre string) (nombreAVD []byte) {
	for i := 0; i < len(nombre); i++ {
		nombreAVD = append(nombreAVD, nombre[i])
	}

	return nombreAVD
}

func editarBitMapYBloque(sectorPart []byte, tamanioBitMap uint32, tamanioBloque uint32, estructura []byte) (cadenaSalida []byte, primerBitLibre uint32) {

	indiceBitMap := 0
	for i := 0; i < int(tamanioBitMap); i++ {
		if sectorPart[i] == 0 {
			sectorPart[i] = 1
			indiceBitMap++
			break
		}
	}

	inicioBloque := tamanioBitMap + (uint32(indiceBitMap)-1)*tamanioBloque
	inidiceEstructura := 0
	for i := inicioBloque; i < uint32(len(sectorPart)); i++ {
		if inidiceEstructura < len(estructura) {
			sectorPart[i] = estructura[inidiceEstructura]
			inidiceEstructura++
		} else {
			break
		}
	}

	return sectorPart, uint32(indiceBitMap)
}

//
func MostrarAVD(contenidoPart []byte, superBoot str.SuperBoot) {

	inicioBitMapAVD := superBoot.ApuntadorBitMapAVD

	n := 0
	for i := inicioBitMapAVD; i < superBoot.ApuntadorAVD; i++ {
		if contenidoPart[i] == 1 {
			a := superBoot.ApuntadorAVD + superBoot.TamanioAVD*uint32(n)
			b := a + superBoot.TamanioAVD
			inicioAVD := contenidoPart[a:b]
			fmt.Println("")
			fmt.Println("NOMBRE DE LA CARPETA *******************************")
			fmt.Println(fmt.Sprintf("%c", inicioAVD))
			fmt.Println(inicioAVD)

		}
		n++
	}

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
	return uint32(numeroTraducido)
}

//
func crearGraficaMBR(contenidoDisco []byte) {

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
