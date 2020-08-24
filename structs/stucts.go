package structs

//
type Comando struct {
	Nombre string
	Valor  string
}

//
type MBR struct {
	Tamanio uint32
	Fecha   [22]byte
	Firma   uint32
	Part1   Particion
	Part2   Particion
	Part3   Particion
	Part4   Particion
}

//
type Particion struct {
	Estado  byte
	Tipo    byte
	Ajuste  byte
	Inicio  uint32
	Tamanio uint32
	Nombre  [16]byte
}
