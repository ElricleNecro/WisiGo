//-tags:libgadget
package ReadGadget

import (
	"log"
	"os"
)

//Header du fichier Gadget : contient les informations générales utiles pour la simulation.
type Header struct {
	Npart    [6]int32   // number of particles of each type in this file.
	Mass     [6]float64 // mass of particles of each type. If 0, then the masses are explicitly stored in the mass-block of the snapshot file, otherwise they are omitted.
	Time     float64    // time of snapshot file.
	Redshift float64    // redshift of snapshot file.

	Flag_sfr      int32     // flags whether the simulation was including star formation.
	Flag_feedback int32     // flags whether feedback was included (obsolete).
	NpartTotal    [6]uint32 // total number of particles of each type in this snapshot. This can be different from npart if one is dealing with a multi-file snapshot.
	Flag_cooling  int32     // flags whether cooling was included.

	Num_files int32 // number of files in multi-file snapshot.

	BoxSize     float64 // box-size of simulation in case periodic boundaries were used.
	Omega0      float64 // matter density in units of critical density.
	OmegaLambda float64 // cosmological constant parameter.
	HubbleParam float64 // Hubble parameter in units of 100 km/sec/Mpc.

	Flag_stellarage int32 // flags whether the file contains formation times of star particles.
	Flag_metals     int32 // flags whether the file contains metallicity values for gas and star particles.

	NpartTotalHighWord     [6]uint32                                                                               // High word of the total number of particles of each type.
	Flag_entropy_instead_u int32                                                                                   // flags that IC-file contains entropy instead of u.
	Fill                   [256 - (6*4 + 6*8 + 8 + 8 + 4 + 4 + 6*4 + 4 + 4 + 8 + 8 + 8 + 8 + 4 + 4 + 6*4 + 4)]byte // Fill the structure to 256 bytes.
}

//Type contenant les informations sur les particules (position, vitesse, masse, type, identité)
type Particule struct {
	Pos  [3]float32 //Positions x, y, z  de la particule.
	Vel  [3]float32 //Vitesses vx, vy, vz de la particule.
	Mass float32    //Masse de la particule.
	Type int64      //Type de la particule.
	Id   int64      //Identité de la particule.
}

//Header et liste de fichier associé à un snapshot.
type Snapshot struct {
	logger        *log.Logger
	f             *os.File //Fichier lu.
	Head          Header   //Header de la simulation.
	to_vit, to_id int64
}

type ByType []Particule

func (e ByType) Len() int {
	return len(e)
}

func (e ByType) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e ByType) Less(i, j int) bool {
	return e[i].Type < e[j].Type
}

type WriteError string

func (e WriteError) Error() string {
	return string(e)
}
