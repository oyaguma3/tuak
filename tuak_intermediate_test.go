package tuak

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"tuak/testvectors"
)

func TestF2345IntermediateVectors(t *testing.T) {
	for testSet := 1; testSet <= 6; testSet++ {
		t.Run(fmt.Sprintf("set%d", testSet), func(t *testing.T) {
			wantIn, wantOut, err := loadF2345IntermediateTestSet(testSet)
			if err != nil {
				t.Fatalf("load intermediate vectors: %v", err)
			}
			runIntermediateTest(t, testSet, "f2345.in", "f2345.out", func(tuak *TUAK) error {
				_, _, _, _, err := tuak.F2345()
				return err
			}, wantIn, wantOut)
		})
	}
}

func TestF5StarIntermediateVectors(t *testing.T) {
	for testSet := 1; testSet <= 6; testSet++ {
		t.Run(fmt.Sprintf("set%d", testSet), func(t *testing.T) {
			wantIn, wantOut, err := loadF5StarIntermediateTestSet(testSet)
			if err != nil {
				t.Fatalf("load intermediate vectors: %v", err)
			}
			runIntermediateTest(t, testSet, "f5star.in", "f5star.out", func(tuak *TUAK) error {
				_, err := tuak.F5Star()
				return err
			}, wantIn, wantOut)
		})
	}
}

func TestF1IntermediateVectors(t *testing.T) {
	for testSet := 1; testSet <= 6; testSet++ {
		t.Run(fmt.Sprintf("set%d", testSet), func(t *testing.T) {
			wantIn, wantOut, err := loadF1IntermediateTestSet(testSet)
			if err != nil {
				t.Fatalf("load intermediate vectors: %v", err)
			}
			runIntermediateTest(t, testSet, "f1.in", "f1.out", func(tuak *TUAK) error {
				_, err := tuak.F1()
				return err
			}, wantIn, wantOut)
		})
	}
}

func TestF1StarIntermediateVectors(t *testing.T) {
	for testSet := 1; testSet <= 6; testSet++ {
		t.Run(fmt.Sprintf("set%d", testSet), func(t *testing.T) {
			wantIn, wantOut, err := loadF1StarIntermediateTestSet(testSet)
			if err != nil {
				t.Fatalf("load intermediate vectors: %v", err)
			}
			runIntermediateTest(t, testSet, "f1star.in", "f1star.out", func(tuak *TUAK) error {
				_, err := tuak.F1Star()
				return err
			}, wantIn, wantOut)
		})
	}
}

func TestTOPcIntermediateVectors(t *testing.T) {
	for testSet := 1; testSet <= 6; testSet++ {
		t.Run(fmt.Sprintf("set%d", testSet), func(t *testing.T) {
			if !hasTOPcIntermediate(testSet) {
				t.Skip("no TOPc intermediate vectors in reference")
			}
			wantIn, wantOut, err := loadTOPcIntermediateTestSet(testSet)
			if err != nil {
				t.Fatalf("load intermediate vectors: %v", err)
			}
			runTOPcIntermediateTest(t, testSet, wantIn, wantOut)
		})
	}
}

func runIntermediateTest(t *testing.T, testSet int, inLabel, outLabelPrefix string, run func(*TUAK) error, wantIn, wantOut []byte) {
	t.Helper()
	v, err := loadTUAKTestSet(testSet)
	if err != nil {
		t.Fatalf("load test set %d: %v", testSet, err)
	}

	var gotIn, gotOut []byte
	hook := func(label string, data []byte) {
		if label == inLabel {
			gotIn = data
		}
		if strings.HasPrefix(label, outLabelPrefix) {
			gotOut = data
		}
	}

	tuak, err := NewWithTOPc(
		decodeHex(t, v.K),
		decodeHex(t, v.Topc),
		decodeHex(t, v.Rand),
		decodeHex(t, v.SQN),
		decodeHex(t, v.AMF),
		append(optionsFromVector(v), WithDebugHook(hook))...,
	)
	if err != nil {
		t.Fatalf("NewWithTOPc: %v", err)
	}

	if err := run(tuak); err != nil {
		t.Fatalf("tuak run: %v", err)
	}
	if gotIn == nil || gotOut == nil {
		t.Fatalf("missing debug captures: in=%v out=%v", gotIn != nil, gotOut != nil)
	}
	if !bytes.Equal(gotIn, wantIn) {
		t.Fatalf("IN mismatch")
	}
	if !bytes.Equal(gotOut, wantOut) {
		t.Fatalf("OUT mismatch")
	}
}

func runTOPcIntermediateTest(t *testing.T, testSet int, wantIn, wantOut []byte) {
	t.Helper()
	v, err := loadTUAKTestSet(testSet)
	if err != nil {
		t.Fatalf("load test set %d: %v", testSet, err)
	}

	var gotIn, gotOut []byte
	hook := func(label string, data []byte) {
		switch label {
		case "topc.in":
			gotIn = data
		case "topc.out", "topc.out.1", "topc.out.2":
			gotOut = data
		}
	}

	k := decodeHex(t, v.K)
	top := decodeHex(t, v.Top)
	_, err = ComputeTOPc(k, top, append(optionsFromVector(v), WithDebugHook(hook))...)
	if err != nil {
		t.Fatalf("ComputeTOPc: %v", err)
	}
	if gotIn == nil || gotOut == nil {
		t.Fatalf("missing debug captures: in=%v out=%v", gotIn != nil, gotOut != nil)
	}
	if !bytes.Equal(gotIn, wantIn) {
		t.Fatalf("TOPc IN mismatch")
	}
	if !bytes.Equal(gotOut, wantOut) {
		t.Fatalf("TOPc OUT mismatch")
	}
}

func loadTUAKTestSet(id int) (testvectors.TUAKVector, error) {
	data, err := testvectors.LoadTUAKVectors()
	if err != nil {
		return testvectors.TUAKVector{}, err
	}
	for _, v := range data.Tests {
		if v.ID == id {
			return v, nil
		}
	}
	return testvectors.TUAKVector{}, errNotFound(fmt.Sprintf("test set %d", id))
}

func loadF2345IntermediateTestSet(testSet int) ([]byte, []byte, error) {
	startPrefix := sectionPrefix(7, testSet)
	endPrefix := nextSectionPrefix(7, testSet, 6, "")
	return loadIntermediatePair(
		"reference/35232-i00_chapter4-7.txt",
		startPrefix, fmt.Sprintf("Test set %d", testSet),
		[]string{"IN when computing f2-f5:"},
		[]string{
			"OUT when computing f2-f5:",
			"OUT/IN after one Keccak iteration, when computing f2-f5:",
			"OUT after second Keccak iteration, when computing f2-f5:",
		},
		endPrefix,
		[]string{"IN when computing f5*:", "As for Test Set"},
	)
}

func loadF5StarIntermediateTestSet(testSet int) ([]byte, []byte, error) {
	startPrefix := sectionPrefix(7, testSet)
	endPrefix := nextSectionPrefix(7, testSet, 6, "")
	in, out, err := loadIntermediatePair(
		"reference/35232-i00_chapter4-7.txt",
		startPrefix, fmt.Sprintf("Test set %d", testSet),
		[]string{"IN when computing f5*:"},
		[]string{
			"OUT when computing f5*:",
			"OUT/IN after one Keccak iteration, when computing f5*:",
			"OUT after second Keccak iteration, when computing f5*:",
		},
		endPrefix,
		nil,
	)
	if errNotFoundValue(err) {
		refSet, ok, refErr := findReferencedF5StarTestSet(testSet)
		if refErr != nil {
			return nil, nil, refErr
		}
		if ok && refSet != testSet {
			return loadF5StarIntermediateTestSet(refSet)
		}
	}
	return in, out, err
}

func loadF1IntermediateTestSet(testSet int) ([]byte, []byte, error) {
	startPrefix := sectionPrefix(6, testSet)
	endPrefix := nextSectionPrefix(6, testSet, 6, "7")
	return loadIntermediatePair(
		"reference/35232-i00_chapter4-7.txt",
		startPrefix, fmt.Sprintf("Test set %d", testSet),
		[]string{"IN when computing f1:"},
		[]string{
			"OUT when computing f1:",
			"OUT/IN after one Keccak iteration, when computing f1:",
			"OUT after second Keccak iteration, when computing f1:",
		},
		endPrefix,
		[]string{"IN when computing f1*:"},
	)
}

func loadF1StarIntermediateTestSet(testSet int) ([]byte, []byte, error) {
	startPrefix := sectionPrefix(6, testSet)
	endPrefix := nextSectionPrefix(6, testSet, 6, "7")
	return loadIntermediatePair(
		"reference/35232-i00_chapter4-7.txt",
		startPrefix, fmt.Sprintf("Test set %d", testSet),
		[]string{"IN when computing f1*:"},
		[]string{
			"OUT when computing f1*:",
			"OUT/IN after one Keccak iteration, when computing f1*:",
			"OUT after second Keccak iteration, when computing f1*:",
		},
		endPrefix,
		nil,
	)
}

func loadTOPcIntermediateTestSet(testSet int) ([]byte, []byte, error) {
	startPrefix := sectionPrefix(6, testSet)
	endPrefix := nextSectionPrefix(6, testSet, 6, "7")
	return loadIntermediatePair(
		"reference/35232-i00_chapter4-7.txt",
		startPrefix, fmt.Sprintf("Test set %d", testSet),
		[]string{"IN when computing TOPc:"},
		[]string{
			"OUT when computing TOPc:",
			"OUT/IN after one Keccak iteration, when computing TOPc:",
			"OUT after second Keccak iteration, when computing TOPc:",
		},
		endPrefix,
		[]string{"IN when computing f1:"},
	)
}

func loadIntermediatePair(path, startPrefix, startNeedle string, inPrefixes, outPrefixes []string, endPrefix string, endCapturePrefixes []string) ([]byte, []byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(string(content), "\n")
	inTest := false
	inCapture := false
	outCapture := false
	var inTokens, outTokens []string

	for _, line := range lines {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, startPrefix) && strings.Contains(s, startNeedle) {
			inTest = true
			continue
		}
		if inTest && endPrefix != "" && strings.HasPrefix(s, endPrefix) {
			break
		}
		if !inTest {
			continue
		}
		if hasPrefixAny(s, inPrefixes) {
			inCapture = true
			outCapture = false
			inTokens = nil
			continue
		}
		if hasPrefixAny(s, outPrefixes) {
			outCapture = true
			inCapture = false
			outTokens = nil
			continue
		}
		if len(endCapturePrefixes) > 0 && hasPrefixAny(s, endCapturePrefixes) {
			inCapture = false
			outCapture = false
			continue
		}

		if inCapture {
			inTokens = append(inTokens, hexByteRE.FindAllString(s, -1)...)
		}
		if outCapture {
			outTokens = append(outTokens, hexByteRE.FindAllString(s, -1)...)
		}
	}

	inBytes, err := decodeHexTokens(inTokens)
	if err != nil {
		return nil, nil, err
	}
	outBytes, err := decodeHexTokens(outTokens)
	if err != nil {
		return nil, nil, err
	}
	if len(inBytes) != 200 || len(outBytes) != 200 {
		return nil, nil, errInvalidVectorSize{inLen: len(inBytes), outLen: len(outBytes)}
	}
	return inBytes, outBytes, nil
}

func hasTOPcIntermediate(testSet int) bool {
	content, err := os.ReadFile("reference/35232-i00_chapter4-7.txt")
	if err != nil {
		return false
	}
	lines := strings.Split(string(content), "\n")
	startPrefix := sectionPrefix(6, testSet)
	endPrefix := nextSectionPrefix(6, testSet, 6, "7")
	inTest := false
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, startPrefix) && strings.Contains(s, fmt.Sprintf("Test set %d", testSet)) {
			inTest = true
			continue
		}
		if inTest && endPrefix != "" && strings.HasPrefix(s, endPrefix) {
			break
		}
		if !inTest {
			continue
		}
		if strings.HasPrefix(s, "IN when computing TOPc:") {
			return true
		}
	}
	return false
}

func findReferencedF5StarTestSet(testSet int) (int, bool, error) {
	content, err := os.ReadFile("reference/35232-i00_chapter4-7.txt")
	if err != nil {
		return 0, false, err
	}
	lines := strings.Split(string(content), "\n")
	startPrefix := sectionPrefix(7, testSet)
	endPrefix := nextSectionPrefix(7, testSet, 6, "")
	inTest := false
	re := regexp.MustCompile(`As for Test Set\s+(\d+)`)
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if strings.HasPrefix(s, startPrefix) && strings.Contains(s, fmt.Sprintf("Test set %d", testSet)) {
			inTest = true
			continue
		}
		if inTest && endPrefix != "" && strings.HasPrefix(s, endPrefix) {
			break
		}
		if !inTest {
			continue
		}
		if strings.Contains(s, "computing f5*") {
			m := re.FindStringSubmatch(s)
			if len(m) == 2 {
				return atoi(m[1]), true, nil
			}
		}
	}
	return 0, false, nil
}

func sectionPrefix(section, testSet int) string {
	return fmt.Sprintf("%d.%d", section, testSet+2)
}

func nextSectionPrefix(section, testSet, max int, fallback string) string {
	if testSet < max {
		return fmt.Sprintf("%d.%d", section, testSet+3)
	}
	return fallback
}

func hasPrefixAny(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

var hexByteRE = regexp.MustCompile(`(?i)\b[0-9a-f]{2}\b`)

func decodeHexTokens(tokens []string) ([]byte, error) {
	if len(tokens) == 0 {
		return nil, errNotFound("hex tokens")
	}
	joined := strings.ToLower(strings.Join(tokens, ""))
	return hex.DecodeString(joined)
}

type errNotFound string

func (e errNotFound) Error() string {
	return "tuak test: not found: " + string(e)
}

func errNotFoundValue(err error) bool {
	_, ok := err.(errNotFound)
	return ok
}

type errInvalidVectorSize struct {
	inLen  int
	outLen int
}

func (e errInvalidVectorSize) Error() string {
	return fmt.Sprintf("tuak test: invalid vector size: in=%d out=%d", e.inLen, e.outLen)
}

func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}
