// Package testvectors loads JSON fixtures for TUAK tests.
package testvectors

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// KeccakVector holds a single Keccak-f[1600] test case.
type KeccakVector struct {
	ID     int    `json:"id"`
	In     string `json:"in"`
	Out    string `json:"out"`
	InLen  int    `json:"in_len"`
	OutLen int    `json:"out_len"`
}

// KeccakFile is the JSON container for Keccak vectors.
type KeccakFile struct {
	Source      string         `json:"source"`
	KeccakF1600 []KeccakVector `json:"keccak_f1600"`
}

// TUAKVector holds TUAK inputs and expected outputs.
type TUAKVector struct {
	ID               int    `json:"id"`
	K                string `json:"k"`
	Rand             string `json:"rand"`
	SQN              string `json:"sqn"`
	AMF              string `json:"amf"`
	Top              string `json:"top"`
	Topc             string `json:"topc"`
	F1               string `json:"f1"`
	F1Star           string `json:"f1_star"`
	F2               string `json:"f2"`
	F3               string `json:"f3"`
	F4               string `json:"f4"`
	F5               string `json:"f5"`
	F5Star           string `json:"f5_star"`
	Klength          int    `json:"klength"`
	MAClength        int    `json:"maclength"`
	CKlength         int    `json:"cklength"`
	IKlength         int    `json:"iklength"`
	RESLength        int    `json:"reslength"`
	KeccakIterations int    `json:"keccak_iterations"`
}

// TUAKFile is the JSON container for TUAK vectors.
type TUAKFile struct {
	Source string       `json:"source"`
	Tests  []TUAKVector `json:"tests"`
}

// F2345Vector holds f2-f5 vectors from TS 35.232.
type F2345Vector struct {
	ID       int    `json:"id"`
	K        string `json:"k"`
	Rand     string `json:"rand"`
	Top      string `json:"top"`
	Topc     string `json:"topc"`
	F2       string `json:"f2"`
	F3       string `json:"f3"`
	F4       string `json:"f4"`
	F5       string `json:"f5"`
	F5Star   string `json:"f5_star"`
	Klength  int    `json:"klength"`
	CKlength int    `json:"cklength"`
	IKlength int    `json:"iklength"`
}

// F2345File is the JSON container for f2-f5 vectors.
type F2345File struct {
	Source string        `json:"source"`
	Tests  []F2345Vector `json:"tests"`
}

// LoadKeccakVectors loads Keccak-f[1600] test vectors.
func LoadKeccakVectors() (*KeccakFile, error) {
	var data KeccakFile
	if err := loadJSON("ts35232_keccak.json", &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// LoadTUAKVectors loads TUAK vectors from JSON fixtures.
func LoadTUAKVectors() (*TUAKFile, error) {
	var data TUAKFile
	if err := loadJSON("ts35233_vectors.json", &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// LoadTUAKVectorsText loads TUAK vectors derived from text extracts.
func LoadTUAKVectorsText() (*TUAKFile, error) {
	var data TUAKFile
	if err := loadJSON("ts35233_vectors_text.json", &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// LoadF2345Vectors loads f2-f5 vectors from JSON fixtures.
func LoadF2345Vectors() (*F2345File, error) {
	var data F2345File
	if err := loadJSON("ts35232_f2345.json", &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func loadJSON(filename string, out interface{}) error {
	path, err := testdataPath(filename)
	if err != nil {
		return err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func testdataPath(filename string) (string, error) {
	root, err := repoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "testdata", filename), nil
}

func repoRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("runtime.Caller failed")
	}

	dir := filepath.Dir(file)
	for {
		candidate := filepath.Join(dir, "testdata", "ts35232_keccak.json")
		if _, err := os.Stat(candidate); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("testdata not found from %s", file)
		}
		dir = parent
	}
}
