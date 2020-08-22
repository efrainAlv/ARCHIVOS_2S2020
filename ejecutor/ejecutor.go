package ejecutor

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	str "../structs"
)

type ejemplo struct {
	letra  int8
	numero int16
}

var comandos []str.Comando

//
func CrearArchivoBinario() {

	fmt.Println("Disco creado")

	file, err := os.Create("/home/helmut/Escritorio/disco.dsk")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	//creamos una variable de tipo int 8 con valor de cero
	cero := int8(0)

	//Nos posicionamos en el byte 1023 que seria el ultimo byte del archivo (el 1023 se puede cambiar dependiendo del tamaño que se quiera el archivo en bytes)
	//seek nos posiciona en el byte 1023 al momento que queramos leer o escribir en el archivo
	file.Seek(1023, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

	//creamos un buffer para pode leer y escribir en archivos
	var buffer bytes.Buffer

	//el metodo write escribe la representacion binaria de la variable cero en la variable buffer
	binary.Write(&buffer, binary.BigEndian, cero)

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
func LeerArchivoBinario(url string) ([]byte, error) {

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
func EditarArchivo(url string, contenido string, inicio int64) {
    // Read Write Mode
    file, err := os.OpenFile(url, os.O_RDWR, 0644)
     
    if err != nil {
        log.Fatalf("failed opening file: %s", err)
    }
    defer file.Close()
	 
	var buffer bytes.Buffer
	buffer.WriteString(contenido)

	cadena := buffer.Bytes()

    len, err := file.WriteAt(cadena, inicio) // Write at 0 beginning
    if err != nil {
        log.Fatalf("failed writing to file: %s", err)
	}
	_=len
}
 