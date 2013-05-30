//-tags:libgadget
package ReadGadget

import "testing"
import "fmt"

func TestReadGadget(t *testing.T) {
	res := NewParticulesFromFile("test/test.dat")
	if len(res) > 5 {
		for i := 0; i < 5; i++ {
			fmt.Println("\t", i, "\t", res[i].Id)
		}
	}
}

func TestWriteGadget(t *testing.T) {
	res := NewParticulesFromFile("test/test.dat")

	tmp := NewSnapshot("test/test.dat")
	head := tmp.Head
	tmp.Close()

	snap := OpenSnapshot("test/test2.dat")
	defer snap.Close()
	snap.Head = head
	fmt.Println(snap.Head, "\n", head)
	if err := snap.DoWrite(res, 1.0, 1.0); err != nil {
		t.Error(err)
	}
}
