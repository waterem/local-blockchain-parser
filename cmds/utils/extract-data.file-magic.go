package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
)

type (
	FileMagicBytesDef struct {
		Filetype  string
		MagicData []byte
	}

	FileMagicBytesResult struct {
		Filetype string
		Reversed bool
	}
)

func (f FileMagicBytesResult) Description() string {
	if f.Reversed {
		return f.Filetype + " (reversed)"
	} else {
		return f.Filetype
	}
}

var fileMagicBytes = []FileMagicBytesDef{
	{"DOC Header", []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}},
	{"DOC Footer", []byte{0x57, 0x6f, 0x72, 0x64, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x2e}},
	{"XLS Header", []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}},
	{"XLS Footer", []byte{0xfe, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x57, 0x00, 0x6f, 0x00, 0x72, 0x00, 0x6b, 0x00, 0x62, 0x00, 0x6f, 0x00, 0x6f, 0x00, 0x6b, 0x00}},
	{"PPT Header", []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}},
	{"PPT Footer", []byte{0xa0, 0x46, 0x1d, 0xf0}},
	{"ZIP Header", []byte{0x50, 0x4b, 0x03, 0x04, 0x14}},
	{"ZIP Footer", []byte{0x50, 0x4b, 0x05, 0x06, 0x00}},
	{"ZIPLock Footer", []byte{0x50, 0x4b, 0x03, 0x04, 0x14, 0x00, 0x01, 0x00, 0x63, 0x00, 0x00, 0x00, 0x00, 0x00}},
	{"JPG Header", []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01}},
	{"GIF Header", []byte{0x47, 0x49, 0x46, 0x38, 0x39, 0x61}},
	{"GIF Footer", []byte{0x21, 0x00, 0x00, 0x3b, 0x00}},
	{"PDF Header", []byte{0x25, 0x50, 0x44, 0x46}},
	{"PDF Header (alternate)", []byte{0x26, 0x23, 0x32, 0x30, 0x35}},
	{"PDF Footer", []byte{0x25, 0x25, 0x45, 0x4f, 0x46}},
	{"Torrent Header", []byte{0x61, 0x6e, 0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65}},
	{"GZ Header", []byte{0x1f, 0x8b, 0x08, 0x08}},
	{"TAR Header", []byte{0x1f, 0x8b, 0x08, 0x00}},
	{"TAR.GZ Header", []byte{0x1f, 0x9d, 0x90, 0x70}},
	{"EPUB Header", []byte{0x50, 0x4b, 0x03, 0x04, 0x0a, 0x00, 0x02, 0x00}},
	{"PNG Header", []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}},
	{"8192 Header", []byte{0x6d, 0x51, 0x51, 0x4e, 0x42}},
	{"4096 Header", []byte{0x6d, 0x51, 0x49, 0x4e, 0x42, 0x46, 0x67, 0x2f}},
	{"2048 Header", []byte{0x95, 0x2e, 0x3e, 0x2e, 0x58, 0x4b, 0x7a}},
	{"Secret Header", []byte{0x52, 0x61, 0x72, 0x21, 0x1a, 0x07, 0x00}},
	{"RAR Header", []byte{0x6d, 0x51, 0x45, 0x4e, 0x42, 0x46, 0x67}},
	{"OGG Header", []byte{0x4f, 0x67, 0x67, 0x53}},
	{"WAV Header", []byte{0x42, 0x49, 0x46, 0x46}},
	{"WAV Header (alternate)", []byte{0x57, 0x41, 0x56, 0x45}},
	{"AVI Header", []byte{0x42, 0x49, 0x46, 0x46}},
	{"AVI Header (alternate)", []byte{0x41, 0x56, 0x49, 0x20}},
	{"MIDI Header", []byte{0x4d, 0x54, 0x68, 0x64}},
	{"7z Header", []byte{0x37, 0x7a, 0xbc, 0xaf, 0x27, 0x1c}},
	{"7z Footer", []byte{0x00, 0x00, 0x00, 0x17, 0x06}},
	{"DMG Header", []byte{0x78, 0x01, 0x73, 0x0d, 0x62, 0x62, 0x60}},
	{"Wikileaks", []byte{0x57, 0x69, 0x6b, 0x69, 0x6c, 0x65, 0x61, 0x6b, 0x73}},
	{"Julian Assange", []byte{0x4a, 0x75, 0x6c, 0x69, 0x61, 0x6e, 0x20, 0x41, 0x73, 0x73, 0x61, 0x6e, 0x67, 0x65}},
	{"Mendax", []byte{0x4d, 0x65, 0x6e, 0x64, 0x61, 0x7}},
}

var wlRipemd160SHA256Hashes = []FileMagicBytesDef{}

func init() {
	data, err := ioutil.ReadFile("./wlhashes/ripemd160-sha256-hashes.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "  ")
		if len(parts) != 2 {
			panic(fmt.Sprintf("len(parts) = %v: %v", len(parts), line))
		}
		digestHex, filename := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		digestBytes, err := hex.DecodeString(digestHex)
		if err != nil {
			panic(err)
		}

		wlRipemd160SHA256Hashes = append(wlRipemd160SHA256Hashes, FileMagicBytesDef{filename + " (ripemd160 + sha256 digest)", digestBytes})
	}
}

func SearchDataForMagicFileBytes(data []byte) []FileMagicBytesResult {
	if data == nil {
		return []FileMagicBytesResult{}
	}

	chMatches := make(chan []FileMagicBytesResult)
	go func() {
		matches := []FileMagicBytesResult{}
		for _, def := range fileMagicBytes {
			if bytes.Contains(data, def.MagicData) {
				matches = append(matches, FileMagicBytesResult{Filetype: def.Filetype, Reversed: false})
			}
		}
		chMatches <- matches
	}()

	chMatchesReversed := make(chan []FileMagicBytesResult)
	go func() {
		matches := []FileMagicBytesResult{}
		for _, def := range fileMagicBytes {
			if bytes.Contains(data, ReverseBytes(def.MagicData)) {
				matches = append(matches, FileMagicBytesResult{Filetype: def.Filetype, Reversed: true})
			}
		}
		chMatchesReversed <- matches
	}()

	chRipemd160SHA256 := make(chan []FileMagicBytesResult)
	go func() {
		matches := []FileMagicBytesResult{}
		for _, def := range wlRipemd160SHA256Hashes {
			if bytes.Contains(data, def.MagicData) {
				matches = append(matches, FileMagicBytesResult{Filetype: def.Filetype, Reversed: false})
			}
		}
		chRipemd160SHA256 <- matches
	}()

	chRipemd160SHA256Reversed := make(chan []FileMagicBytesResult)
	go func() {
		matches := []FileMagicBytesResult{}
		for _, def := range wlRipemd160SHA256Hashes {
			if bytes.Contains(data, ReverseBytes(def.MagicData)) {
				matches = append(matches, FileMagicBytesResult{Filetype: def.Filetype, Reversed: true})
			}
		}
		chRipemd160SHA256Reversed <- matches
	}()

	matches := []FileMagicBytesResult{}
	matches = append(matches, <-chMatches...)
	matches = append(matches, <-chMatchesReversed...)
	matches = append(matches, <-chRipemd160SHA256...)
	matches = append(matches, <-chRipemd160SHA256Reversed...)

	return matches
}
