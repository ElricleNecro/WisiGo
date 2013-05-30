//-tags:libgadget
package ReadGadget

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"sort"
)

func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		n := sizeof(t.Elem())
		if n < 0 {
			return -1
		}
		return t.Len() * n

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}
	return -1
}

func (e *Snapshot) SetFlag(Flag_sfr, Flag_feedback, Flag_cooling, Flag_stellarage, Flag_metals, Flag_entropy_instead_u int32) {
	e.Head.Flag_sfr = Flag_sfr
	e.Head.Flag_feedback = Flag_feedback
	e.Head.Flag_cooling = Flag_cooling
	e.Head.Flag_stellarage = Flag_stellarage
	e.Head.Flag_metals = Flag_metals
	e.Head.Flag_entropy_instead_u = Flag_entropy_instead_u
}

func (e *Snapshot) SetBoxSize(Box float64) {
	e.Head.BoxSize = Box
}

func (e *Snapshot) SetCosmology(Omega0, OmegaLambda, HubbleParam float64) {
	e.Head.Omega0 = Omega0
	e.Head.OmegaLambda = OmegaLambda
	e.Head.HubbleParam = HubbleParam
}

func (e *Snapshot) SetTime(time, redshift float64) {
	e.Head.Time = time
	e.Head.Redshift = redshift
}

func (e *Snapshot) SetNumFile(nb int32) {
	e.Head.Num_files = nb
}

func (e *Snapshot) write(don interface{}) (err error) {
	var blksize int32 = 0

	blksize = int32(sizeof(reflect.ValueOf(don).Type()))
	err = binary.Write(e.f, binary.LittleEndian, blksize)
	if err != nil {
		//e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(e.f, binary.LittleEndian, don)
	if err != nil {
		//e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(e.f, binary.LittleEndian, blksize)
	if err != nil {
		//e.logger.Println("binary.Write failed:", err)
		return err
	}

	return nil
}

func (e Snapshot) DoWrite(part []Particule, LongFact, VitFact float32) error {
	for i, _ := range e.Head.Npart {
		e.Head.Npart[i] = 0
		e.Head.NpartTotal[i] = 0
		e.Head.Mass[i] = 0.
	}
	//sort.Sort(ByType(part))
	for _, v := range part {
		e.Head.Npart[v.Type] += 1
		e.Head.NpartTotal[v.Type] += 1
		e.Head.Mass[v.Type] = float64(v.Mass)
	}
	//e.Head.Num_files = 1
	//e.Head.Time = 0.0
	//e.Head.Redshift = 0.0

	if err := e.write(e.Head); err != nil {
		return WriteError(fmt.Sprint("Error While writing Header : ", err))
	}
	tmp := make([]float32, 3*len(part))
	z := 0
	for _, v := range part {
		for j := 0; j < 3; j++ {
			tmp[z] = v.Pos[j] / float32(LongFact)
			z += 1
		}
	}
	if err := e.write(tmp); err != nil {
		return WriteError(fmt.Sprint("Error While writing Positions : ", err))
	}

	z = 0
	for _, v := range part {
		for j := 0; j < 3; j++ {
			tmp[z] = v.Vel[j] / float32(VitFact)
			z += 1
		}
	}
	if err := e.write(tmp); err != nil {
		return WriteError(fmt.Sprint("Error While writing Velocities : ", err))
	}

	id := make([]int32, len(part))
	for i, v := range part {
		id[i] = int32(v.Id)
	}
	if err := e.write(id); err != nil {
		return WriteError(fmt.Sprint("Error While writing Identities : ", err))
	}

	return nil
}

func (e Snapshot) WriteAll(fname string, part []Particule, LongFact, VitFact, MassFact float64) error {
	var header Header
	var blksize int = 0
	tmp := make([]float32, 3*len(part))
	tmp_id := make([]int32, 3*len(part))

	f, err := os.OpenFile(fname, os.O_WRONLY, 0)
	if err != nil {
		e.logger.Printf("Ouverture du fichier raté : %v\n", err)
		return err
	}
	defer f.Close()

	sort.Sort(ByType(part))
	for _, v := range part {
		header.Npart[v.Type] += 1
		header.NpartTotal[v.Type] += 1
		header.Mass[v.Type] = float64(v.Mass)
	}

	header.Time = 0.0
	header.Redshift = 0.0

	header.Num_files = 1
	header.BoxSize = 0.0 //BoxSize / LongFact; // / 3.086e16;
	header.Omega0 = 0.0
	header.OmegaLambda = 0.0
	header.HubbleParam = 0.0

	//---------------------------------------------------------
	//		Écriture du header
	//---------------------------------------------------------
	blksize = sizeof(reflect.ValueOf(header).Type())
	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, header)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	//---------------------------------------------------------
	//		Écriture des Positions
	//---------------------------------------------------------
	z := 0
	for _, v := range part {
		for j := 0; j < 3; j++ {
			tmp[z] = v.Pos[j] / float32(LongFact)
			z += 1
		}
	}
	blksize = sizeof(reflect.ValueOf(tmp).Type()) //3 * len(part) * 4
	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, tmp)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	//---------------------------------------------------------
	//		Écriture des Vitesses
	//---------------------------------------------------------
	z = 0
	for _, v := range part {
		for j := 0; j < 3; j++ {
			tmp[z] = v.Vel[j] / float32(VitFact)
			z += 1
		}
	}

	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, tmp)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	//---------------------------------------------------------
	//		Écriture des Identités
	//---------------------------------------------------------
	for i, v := range part {
		tmp_id[i] = int32(v.Id)
	}

	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, tmp)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	err = binary.Write(f, binary.LittleEndian, blksize)
	if err != nil {
		e.logger.Println("binary.Write failed:", err)
		return err
	}

	return nil
}
