package analizador

import(

	"fmt"
	"os"
	s "strings"
	"bufio"
)

var estado int =0;
var comandosLeidos = make([]comando, 0)


type comando struct{
	nombre string
	valor string	
}


func analizarComando(lineaComnados string){
	
	comandos:= s.Split(lineaComnados, " ")

	for i := 0; i < len(comandos); i++ {
		
		/*if comandos[i] == "pausa" {
			lector := bufio.NewReader(os.Stdin)
			fmt.Println("En pausa ...")
			input, _:=lector.ReadString('\n')
			_ = input
		}*/

		if comandos[i] == "pausa" {
			bufio.NewReader(os.Stdin).ReadBytes('\n') 
		
		}else if comandos[i] == "exec" {
			comandosLeidos = append(comandosLeidos, comando{"exec", "inicial"})
		}else if comandos[i] == "exec" {
			comandosLeidos = append(comandosLeidos, comando{"exec", "inicial"})
		}else if s.ContainsAny("-path" ,comandos[i]) {
			comandosLeidos = append(comandosLeidos, comando{"-path", comandos[i]})
		}

	}

}


func analizarParametros(comadno string, ){



}


func ejecutar(){


}


//
func Leer(url string){

	file, err:=os.Open(url)
	check(err)
	fileInfo, err:= os.Lstat(url)
	check(err)

	cadenaBytes := make([]byte, fileInfo.Size()) //OBTIENE LA CADENA DE BYTES DEL ARCHIVO
	check(err)

	n, err:=file.Read(cadenaBytes) //SE LEE EL ARCHIVO, SE PASA COMO PARAMETRO EL TAMAÑO EN BYTES DEL ARCHIVO
	check(err)

	fmt.Println("BYTES LEIDOS: ", n)

	cadena := s.Split(string(cadenaBytes), "\n") //SEPARA EL ARCHIVO POR LINEAS

	//LEE EL ARCHIVO LINEA POR LINEA
	for i := 0; i < len(cadena)-1; i++ {
		fmt.Println("Linea ", i+1, ": ",cadena[i])
		analizarComando(cadena[i])
	}

}


func check(e error){
	if e!= nil{
		panic(e)
	}
}
