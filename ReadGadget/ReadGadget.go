//-tags:libgadget
//Package permettant la lecture des fichiers snapshot de Gadget 2 :
//	Volker Springel,
//	The cosmological simulation code GADGET-2,
//	Royal Astronomical Society, pages 1105-1134,
//	2005.
//
//Ces fichiers sont organisé de la façon suivante, en respectant la norme Fortran :
//	Taille du Header
//		Header
//	Taille du Header
//	Taille du bloc position
//		Position sous la forme d'un tableau applatit
//	Taille du bloc position
//	Taille du bloc vitesse
//		Vitesse sous la forme d'un tableau applatit
//	Taille du bloc vitesse
//	Taille du bloc identité
//		Identités des particules
//	Taille du bloc identité
//Et ainsi de suite, selon ce qui est demandé à Gadget 2 en sortie.
//Les Positions... sont organisé par type croissant.
//
//Normalement, les fichiers Gadget peuvent être séparé en plusieurs fichiers ; mais cette fonctionnalité n'est pas encore supporté par le code actuel.
package ReadGadget

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

func NewParticulesFromFile(file string) []Particule {
	Tmp := NewSnapshot(file)
	defer Tmp.Close()
	pos := Tmp.GetPositions()
	vit := Tmp.GetVitesses()
	id := Tmp.GetIds()

	fmt.Println("\033[36m", len(id), "\033[00m")
	j := 0
	for i := 0; i < 5; i++ {
		fmt.Println("\033[31m(", pos[j], ", ", pos[j+1], ", ", pos[j+2], ") -- (", vit[j], ", ", vit[j+1], ", ", vit[j+2], ")", " -- ", id[i], "\033[00m\n")
		j += 3
	}

	res := make([]Particule, len(pos)/3)

	for i := 0; i < len(res); i++ {
		for j := 0; j < 3; j++ {
			res[i].Pos[j] = pos[i+j]
			res[i].Vel[j] = vit[i+j]
		}
		res[i].Id = int64(id[i])
		switch {
		case int32(i) >= Tmp.Head.Npart[0]+Tmp.Head.Npart[1]+Tmp.Head.Npart[2]+Tmp.Head.Npart[3]+Tmp.Head.Npart[4]:
			res[i].Type = int64(5)
		case int32(i) >= Tmp.Head.Npart[0]+Tmp.Head.Npart[1]+Tmp.Head.Npart[2]+Tmp.Head.Npart[3]:
			res[i].Type = int64(4)
		case int32(i) >= Tmp.Head.Npart[0]+Tmp.Head.Npart[1]+Tmp.Head.Npart[2]:
			res[i].Type = int64(3)
		case int32(i) >= Tmp.Head.Npart[0]+Tmp.Head.Npart[1]:
			res[i].Type = int64(2)
		case int32(i) >= Tmp.Head.Npart[0]:
			res[i].Type = int64(1)
		default:
			res[i].Type = int64(0)
		}
	}

	return res
}

func NewParticulesTypeFromFile(file string, Type int) []Particule {
	Tmp := NewSnapshot(file)
	defer Tmp.Close()
	pos := Tmp.GetPosition(Type)
	vit := Tmp.GetVitesse(Type)
	id := Tmp.GetId(Type)

	res := make([]Particule, len(pos)/3)

	for i := 0; i < len(res); i++ {
		for j := 0; j < 3; j++ {
			res[i].Pos[j] = pos[i+j]
			res[i].Vel[j] = vit[i+j]
		}
		res[i].Type = int64(Type)
		res[i].Id = id[i]
	}

	return res
}

//Crée un objet Snapshot nous permettant de manipuler le fichier, et charge automatiquement l'entête.
func NewSnapshot(fname string) *Snapshot {
	var e Snapshot
	var err error

	e.f, err = os.Open(fname)
	(&e).GetHeader()

	var move int64 = 0
	for _, v := range e.Head.Npart {
		move += int64(3 * v) //int64(3*e.Head.Npart[i])
	}

	e.to_vit = to_pos + move*4 + to_header + to_header
	e.to_id = e.to_vit + move*4 + to_header + to_header

	fmt.Println(to_header, to_pos, e.to_vit, e.to_id, move)

	if err != nil {
		log.Panic(err)
	}

	return &e
}

func OpenSnapshot(fname string) *Snapshot {
	var e Snapshot
	var err error

	e.f, err = os.Create(fname)
	if err != nil {
		log.Panic(err)
	}
	e.logger = log.New(os.Stderr, "Snapshot : ", log.LstdFlags|log.Lshortfile)

	return &e
}

func (e Snapshot) getVal(gi int64, data interface{}) error {
	if newpos, err := e.f.Seek(int64(gi), os.SEEK_SET); newpos != int64(gi) || err != nil {
		return err
	}

	if err := binary.Read(e.f, binary.LittleEndian, data); err != nil {
		return err
	}

	return nil
}

//Récupére les postions de touteles particules du fichier, indiféremment du type.
func (e *Snapshot) GetPositions() []float32 {
	pos := make([]float32, 3*(e.Head.Npart[0]+e.Head.Npart[1]+e.Head.Npart[2]+e.Head.Npart[3]+e.Head.Npart[4]+e.Head.Npart[5]))

	if err := e.getVal(to_pos, &pos); err != nil {
		log.Panic(err)
	}

	return pos
}

func (e *Snapshot) GetVitesses() []float32 {
	pos := make([]float32, 3*(e.Head.Npart[0]+e.Head.Npart[1]+e.Head.Npart[2]+e.Head.Npart[3]+e.Head.Npart[4]+e.Head.Npart[5]))

	if err := e.getVal(e.to_vit, &pos); err != nil {
		log.Panic(err)
	}

	return pos
}

func (e *Snapshot) GetIds() []int32 {
	pos := make([]int32, e.Head.Npart[0]+e.Head.Npart[1]+e.Head.Npart[2]+e.Head.Npart[3]+e.Head.Npart[4]+e.Head.Npart[5])

	fmt.Println("\033[35m", e.to_id, " -- ", len(pos), " (", e.Head.Npart[0]+e.Head.Npart[1]+e.Head.Npart[2]+e.Head.Npart[3]+e.Head.Npart[4]+e.Head.Npart[5], ")\033[00m")
	if err := e.getVal(e.to_id, &pos); err != nil {
		log.Panic(err)
	}

	return pos
}

//Récupére uniquement les particules du type "Type" voulu.
func (e *Snapshot) GetPosition(Type int) []float32 {
	var move int64 = 0
	for i := 0; i <= Type; i++ {
		move += int64(3 * e.Head.Npart[i])
	}

	pos := make([]float32, 3*e.Head.Npart[Type])

	if err := e.getVal(to_pos+move, &pos); err != nil {
		log.Panic(err)
	}

	return pos
}

func (e *Snapshot) GetVitesse(Type int) []float32 {
	var move int64 = 0
	for i := 0; i <= Type; i++ {
		move += int64(3 * e.Head.Npart[i])
	}

	pos := make([]float32, 3*e.Head.Npart[Type])

	if err := e.getVal(e.to_vit+move, &pos); err != nil {
		log.Panic(err)
	}

	return pos
}

func (e *Snapshot) GetId(Type int) []int64 {
	var move int64 = 0
	for i := 0; i <= Type; i++ {
		move += int64(e.Head.Npart[i])
	}

	pos := make([]int64, e.Head.Npart[Type])

	if err := e.getVal(e.to_id+move, &pos); err != nil {
		log.Panic(err)
	}

	return pos
}

//Récupére le Header du fichier, fait automatiquement, à l'ouverture du fichier.
func (e *Snapshot) GetHeader() {
	if newpos, err := e.f.Seek(to_header, os.SEEK_SET); newpos != to_header || err != nil {
		fmt.Printf("GetHeader():49 ")
		log.Panic(err)
	}

	if err := binary.Read(e.f, binary.LittleEndian, &e.Head); err != nil {
		fmt.Printf("GetHeader():59 ")
		log.Panic(err)
	}
}

//Ferme le fichier.
func (e *Snapshot) Close() {
	if err := e.f.Close(); err != nil {
		log.Panic(err)
	}
}
