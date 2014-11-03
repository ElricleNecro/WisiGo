package ReadGadget

import (
	"binary"
	"fmt"
	"log"
	"os"
)

type Gadget struct {
	name   string
	format int
	Part   []Particule
	header Header
	log    *log.Logger
}

func New(name string, format int) *Gadget {
	res := new(Gadget)
	res.name = name
	res.format = format
	return res
}

func (c *Gadget) SetLogger(val *log.Logger) {
	c.log = val
}

func (c *Gadget) GetLogger() *log.Logger {
	return c.log
}

type UnsupportedFormat int

func (c UnsupportedFormat) Error() string {
	return fmt.Sprint("Unsupported format: ", c)
}

type NotImplemented string

func (c NotImplemented) Error() string {
	return fmt.Sprint("You are trying to use a not implemented method: ", c)
}

func (c *Gadget) Read() error {
	switch c.format {
	case 1:
		return c.read_format1()
	case 2:
		return c.read_format2()
	case 3:
		return c.read_format3()
	default:
		return UnsupportedFormat(c.format)
	}

}

func (c *Gadget) read_format1() error {
	var size, nb_files, pc, pc_new int = 0, 1, 0, 0
	var err error = nil
	var file *os.File

	for i := 0; i < nb_files; i++ {
		if file, err = os.Open(c.name); err != nil {
			return err
		}

		pc = pc_new
	}

	return NotImplemented("Read()->read_format1()")
}

func (c *Gadget) read_format2() error {
	return NotImplemented("Read()->read_format2()")
}

func (c *Gadget) read_format3() error {
	return NotImplemented("Read()->read_format3()")
}
