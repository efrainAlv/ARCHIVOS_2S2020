package structs

import (
	"unsafe"
)

//
var TamSuperBoot = unsafe.Sizeof(SuperBoot{})

//
var TamAVD = unsafe.Sizeof(AVD{})

//
var TamDetalleDirect = unsafe.Sizeof(DetalleDirectorio{})

//
var TamInodo = unsafe.Sizeof(Inodo{})

//
var TamBloque = unsafe.Sizeof(Bloque{})

//
var TamBitacora = unsafe.Sizeof(Bitacora{})

//
type Comando struct {
	Nombre string
	Valor  string
}

/*
	Tamaño real = 138bytes
*/
type MBR struct {
	Tamanio uint32
	Fecha   [22]byte
	Firma   uint32
	Part1   Particion
	Part2   Particion
	Part3   Particion
	Part4   Particion
}

/*
	Tamaño real = 138bytes
*/
type EBR struct {
	Estado       byte
	Ajuste       byte
	Inicio       uint32
	Tamanio      uint32
	SiguienteEBR uint32
	Nombre       [16]byte
}

/*
	Tamaño real = 27bytes
*/
type Particion struct {
	Estado  byte
	Tipo    byte
	Ajuste  byte
	Inicio  uint32
	Tamanio uint32
	Nombre  [16]byte
}

//
var ParticionesMontadas []ParticionMontada

//
type ParticionMontada struct {
	Particion Particion
	Letra     byte
	Numero    uint16
	Ruta      string
}

/*
	Tamaño real = 172bytes
*/
type SuperBoot struct {
	NombreDisco                  [20]byte
	CantidadAVD                  uint32
	CantidadDetalleDirect        uint32
	CantidadInodos               uint32
	CantidadBloques              uint32
	CantidadAVDLibres            uint32
	CantidadDetalleDirectLibres uint32
	CantidadInodosLibres         uint32
	CantidadBloquesLibres        uint32
	FechaCreacion                [22]byte
	FechaUltimoMontaje           [22]byte
	NumeroMontajes               uint16
	ApuntadorBitMapAVD             uint32
	ApuntadorAVD                 uint32
	ApuntadorBitMapDetalleDirect   uint32
	ApuntadorDetalleDirect       uint32
	ApuntadorBitMapInodos          uint32
	ApuntadorInodos              uint32
	ApuntadorBitMapBloques         uint32
	ApuntadorBloques             uint32
	ApuntadorBitacora             uint32
	TamanioAVD                   uint32
	TamanioDetalleDirect         uint32
	TamanioInodo                 uint32
	TamanioBloque                uint32
	TamanioBitacora              uint32
	PrimerAVDLibre               uint32
	PrimerDetalleDirectLibre     uint32
	PrimerInodoLibre             uint32
	PrimerBloqueLibre            uint32
	NumeroMagico                 uint32
}

/*
	Tamaño real = 84bytes
*/
type AVD struct {
	FechaCreacion          [22]byte
	NombreDirectorio       [20]byte
	SubAVD                 [6]uint32
	ApuntadorDetalleDirect uint32
	ApuntadorIndirecto     uint32
	IDPropietario          uint32
	IDGrupo                uint32
	Permisos               uint16
}

/*
	Tamaño real = 344bytes
*/
type DetalleDirectorio struct {
	Archivos           [5]InforArchivo
	ApuntadorIndirecto uint32
}

/*
	Tamaño real = 68bytes
*/
type InforArchivo struct {
	Nombre            [20]byte
	ApuntarInodo      uint32
	FechaCreacion     [22]byte
	FechaModificacion [22]byte
}

/*
	Tamaño real = 40bytes
*/
type Inodo struct {
	Numero             uint32
	TamanioArchivo     uint32
	NumeroBloques      uint16
	Bloques            [4]uint32
	ApuntadorIndirecto uint32
	IDPropietario      uint32
	IDGrupo            uint32
	Permisos           uint16
}

/*
	Tamaño real = 25bytes
*/
type Bloque struct {
	Datos [25]byte
}

/*
	Tamaño real = 64bytes
*/
type Bitacora struct {
	TipoOperacion [20]byte
	Tipo          byte
	Nombre        [20]byte
	Contenido     byte
	Fecha         [22]byte
}
