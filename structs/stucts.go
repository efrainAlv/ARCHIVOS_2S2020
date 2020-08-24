package structs

//
type Comando struct {
	Nombre string
	Valor  string
}

//
type MBR struct {

	Tamanio uint32
	Fecha [19]byte
	Firma uint32
}

//
type Particion struct {

	Estado []byte
	Tipo []byte
	Ajuste[]byte
	Inicio int32
	Tamanio int32
	Nombre [16]byte

}


func prueba(){

}