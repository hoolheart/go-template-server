package common

import (
	"crypto/rand"
	"crypto/sha1"
	"log"
	"fmt"
	"io"
	"encoding/json"
)

// CreateUUID a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func CreateUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Encrypt plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}

// ExportJSON exports a json data to an output
func ExportJSON(d interface{}, w io.Writer) (err error) {
	encoder := json.NewEncoder(w)//prepare encoder
	return encoder.Encode(d)//encode json
}

// LoadJSON parses an input into specific data by using JSON format
func LoadJSON(r io.Reader, d interface{}) (err error) {
	decoder := json.NewDecoder(r)//prepare decoder
	return decoder.Decode(d)//decode json
}
