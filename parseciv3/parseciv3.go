package parseciv3

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"

	"github.com/myjimnelson/c3sat/civ3decompress"
)

type baseClass struct {
	name   string
	length uint32
	buffer bytes.Buffer
}

type Civ3Data struct{}

// ReadFile takes a filename and returns the decompressed file data or the raw data if it's not compressed. Also returns true if compressed.
func ReadFile(path string) ([]byte, bool, error) {
	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		return nil, false, ReadError{err}
	}
	defer file.Close()

	var compressed bool
	var data []byte
	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		return nil, false, ReadError{err}
	}
	// reset pointer to parse from beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, false, ReadError{err}
	}
	switch {
	case header[0] == 0x00 && (header[1] == 0x04 || header[1] == 0x05 || header[1] == 0x06):
		compressed = true
		data, err = civ3decompress.Decompress(file)
		if err != nil {
			return nil, false, ReadError{err}
		}
	default:
		// log.Println("Not a compressed file. Proceeding with uncompressed stream.")
		// TODO: I'm sure I'm doing this in a terribly inefficient way. Need to refactor everything to pass around file pointers I think
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, false, ReadError{err}
		}
	}
	return data, compressed, error(nil)

}

// Parseciv3 ...
func Parseciv3(path string) error {
	civdata, compressed, err := ReadFile(path)
	_ = compressed
	r := bytes.NewReader(civdata)
	// get the first four bytes to determine file type
	header, err := readBytes(r, 4)
	if err != nil {
		return ReadError{err}
	}
	// reset pointer to parse from beginning
	r.Seek(0, 0)
	switch string(header) {
	case "CIV3":
		// log.Println("Civ3 save file detected")
		// readcivheader(r)
		// readbic(r)
	case "BIC ", "BICX":
		// log.Fatal("Civ3 BIC file detected. Currently not parsing these directly.")
		// readbic(r)
	default:
		log.Fatalf("Civ3 file not detected. First four bytes:\n%s", hex.Dump(header))
	}
	return error(nil)
}

// func check(e error) {
// 	if e != nil {
// 		return ReadError{err}
// 	}
// }

// readBytes repeatedly calls bytes.Reader.ReadByte()
func readBytes(r *bytes.Reader, n int) ([]byte, error) {
	var out bytes.Buffer
	for i := 0; i < n; i++ {
		byt, err := r.ReadByte()
		if err != nil {
			return []byte(nil), ReadError{err}
		}
		out.WriteByte(byt)
	}
	return out.Bytes(), error(nil)
}

/*
// oops, I spent time refactoring this for error handlnig but I'm not going to keep it
func readBase(r *bytes.Reader) (baseClass, error) {
	var c baseClass

	buffer, err := readBytes(r, 8)
	if err != nil {
		return c, err
	}
	c.name = string(buffer[:4])
	c.length = binary.LittleEndian.Uint32(buffer[4:4])

	buffer, err = readBytes(r, int(c.length))
	if err != nil {
		return c, err
	}
	c.buffer.Write(buffer)

	return c, error(nil)
}

func somethingsdifferent(s string, r *bytes.Reader) {
	// seeking backwards is causing EOF when I pass
	// r.Seek(-4, 2)
	log.Fatalf("%s\n%s\n", s, hex.Dump(readBytes(r, 256)))
}

func readcivheader(r *bytes.Reader) {
	var civ3header civ3
	err := binary.Read(r, binary.LittleEndian, &civ3header)
	if err != nil {
		return ReadError{err}
	}
	log.Println(civ3header)
	// civ3header := readBytes(r, 30)
	// _ = civ3header

	// Seems to be four or five ints followed by 2 strings, relative paths to bic resources and bic file
	bicheader := readBase(r)
	_ = bicheader
}

func readbic(r *bytes.Reader) {
	bicqvernum := readBytes(r, 8)
	switch string(bicqvernum) {
	case "BICQVER#", "BIC VER#", "BICXVER#":
	default:
		log.Fatal(string(bicqvernum))
		somethingsdifferent(string(bicqvernum), r)
	}
	// always seems to be 1
	bicone := readBytes(r, 4)
	if binary.LittleEndian.Uint32(bicone) != 1 {
		somethingsdifferent(string(binary.LittleEndian.Uint32(bicone)), r)
	}
	bicdescriptionlength := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
	if bicdescriptionlength != 720 {
		somethingsdifferent(string(binary.LittleEndian.Uint32(bicone)), r)
	}
	bicdescription := readBytes(r, bicdescriptionlength)
	_ = bicdescription
	// log.Println(hex.Dump(bicdescription))
	// log.Println(hex.Dump(bicdescription[:16]))
	// log.Println(hex.Dump(bicdescription[8:12]))

	// At this point, my epic SAVs have GAME, downloaded scenario-based games have BLDG

	var buffernext []byte
	bicnext := string(readBytes(r, 4))
	switch bicnext {
	case "BLDG":
		log.Println(bicnext)
		numbuildings := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		for i := 0; i < numbuildings; i++ {
			buffernext = readBytes(r, 0x110)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
		// CTZN
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numcitizentypes := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		for i := 0; i < numcitizentypes; i++ {
			buffernext := readBytes(r, 0x80)
			_ = buffernext
		}
		// CULT
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numcult := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numcult
		for i := 0; i < numcult; i++ {
			buffernext := readBytes(r, 0x5c)
			_ = buffernext
		}
		// DIFF
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numdiff := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numdiff
		for i := 0; i < numdiff; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x7c)
			_ = buffernext
		}
		// ERAS
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numeras := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numeras
		for i := 0; i < numeras; i++ {
			buffernext := readBytes(r, 0x10c)
			_ = buffernext
		}
		// ESPN
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numespn := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numespn
		for i := 0; i < numespn; i++ {
			buffernext := readBytes(r, 0xec)
			_ = buffernext
		}
		// EXPR
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numexpr := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numexpr
		for i := 0; i < numexpr; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x2c)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
		// FLAV
		bicnext = string(readBytes(r, 4))
		if (bicnext) == "FLAV" {
			log.Println(bicnext)
			numflavgroup := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
			for i := 0; i < numflavgroup; i++ {
				numflavors := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
				for i := 0; i < numflavors; i++ {

					// for i := 0; i < 1; i++ {
					buffernext := readBytes(r, 0x124)
					_ = buffernext
				}
			}
			bicnext = string(readBytes(r, 4))
		}
		// GOOD
		// bicnext already read because of optional FLAV
		log.Println(bicnext)
		numgood := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numgood
		for i := 0; i < numgood; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x5c)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
		// GOVT
		// apparently GOVT length changed during C3C as Mesopotamia and Mesoamerica scenario BIQs get misaligned in this parser.
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numgovt := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numgovt
		for i := 0; i < numgovt; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x23c)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
		// RULE
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numrule := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numrule
		for i := 0; i < numrule; i++ {
			// for i := 0; i < 1; i++ {
			buffernext := readBytes(r, 0x2d4)
			_ = buffernext
			// log.Println(hex.Dump(buffernext))
		}
		// PRTO
		bicnext = string(readBytes(r, 4))
		log.Println(bicnext)
		numprto := int(binary.LittleEndian.Uint32(readBytes(r, 4)))
		_ = numprto
		// for i := 0; i < numprto; i++ {
		for i := 0; i < 3; i++ {
			buffernext := readBytes(r, 0x103)
			_ = buffernext
			log.Println(hex.Dump(buffernext))
		}
	default:
		log.Println("Unexpected class name: ", bicnext)

	}

	bicnext = string(readBytes(r, 4))
	log.Println(bicnext)
	log.Println(binary.LittleEndian.Uint32(readBytes(r, 4)))

	// var bicgame baseClass
	// bicgame.length = 1
	// for bicgame.name != "GAME" && 0 < bicgame.length && bicgame.length < 500 {
	// 	bicgame = readBase(r)
	// 	log.Println(bicgame.name, bicgame.length)
	// 	// log.Println(bicgame.name, hex.Dump(bicgame.buffer.Bytes()))
	// }

	log.Println(hex.Dump(readBytes(r, 0x40)))
	log.Println("")

}
*/
