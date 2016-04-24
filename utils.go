package main

import (
	"encoding/binary"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var animals []string
var adjectives []string
var aboutPageText string

type versionsInfo struct {
	VersionDate string
	VersionNum  int
}

func init() {
	rand.Seed(time.Now().Unix())
	animalsText, _ := ioutil.ReadFile(path.Join(RuntimeArgs.SourcePath, "static/text/animals"))
	animals = strings.Split(string(animalsText), ",")
	adjectivesText, _ := ioutil.ReadFile(path.Join(RuntimeArgs.SourcePath, "static/text/adjectives"))
	adjectives = strings.Split(string(adjectivesText), "\n")
}

func randomAnimal() string {
	return strings.Replace(strings.Title(animals[rand.Intn(len(animals)-1)]), " ", "", -1)
}

func randomAdjective() string {
	return strings.Replace(strings.Title(adjectives[rand.Intn(len(adjectives)-1)]), " ", "", -1)
}

func randomAlliterateCombo() (combo string) {
	combo = ""
	// // first determine which names are taken from program data
	// takenNames := []string{}
	// err := db.View(func(tx *bolt.Tx) error {
	// 	// Assume bucket exists and has keys
	// 	b := tx.Bucket([]byte("programdata"))
	// 	c := b.Cursor()
	// 	for k, v := c.First(); k != nil; k, v = c.Next() {
	// 		takenNames = append(takenNames, strings.ToLower(string(v)))
	// 	}
	// 	return nil
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(takenNames)
	// generate random alliteration thats not been used
	for {
		animal := randomAnimal()
		adjective := randomAdjective()
		if animal[0] == adjective[0] { //&& stringInSlice(strings.ToLower(adjective+animal), takenNames) == false {
			combo = adjective + animal
			break
		}
	}
	return
}

// is there a string in a slice?
func stringInSlice(s string, strings []string) bool {
	for _, k := range strings {
		if s == k {
			return true
		}
	}
	return false
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func contentType(filename string) string {
	switch {
	case strings.Contains(filename, ".css"):
		return "text/css"
	case strings.Contains(filename, ".jpg"):
		return "image/jpeg"
	case strings.Contains(filename, ".png"):
		return "image/png"
	case strings.Contains(filename, ".js"):
		return "application/javascript"
	}
	return "text/html"
}

func diffRebuildtexts(diffs []diffmatchpatch.Diff) []string {
	text := []string{"", ""}
	for _, myDiff := range diffs {
		if myDiff.Type != diffmatchpatch.DiffInsert {
			text[0] += myDiff.Text
		}
		if myDiff.Type != diffmatchpatch.DiffDelete {
			text[1] += myDiff.Text
		}
	}
	return text
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func getImportantVersions(p WikiData) ([]versionsInfo, time.Duration) {
	// defer timeTrack(time.Now(), "getImportantVersions")
	m := map[int]int{}
	lastTime := time.Now().AddDate(0, -1, 0)
	totalTime := time.Now().Sub(time.Now())
	for i := range p.Diffs {
		parsedTime, _ := time.Parse(time.ANSIC, p.Timestamps[i])
		duration := parsedTime.Sub(lastTime)
		if duration.Minutes() < 3 {
			totalTime += duration
		}
		m[i] = int(duration.Seconds())
		if i > 0 {
			m[i-1] = m[i]
		}
		// On to the next one
		lastTime = parsedTime
	}

	// Sort in order of decreasing diff times
	n := map[int][]int{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))

	// Get the top 4 biggest diff times
	var importantVersions []int
	var r []versionsInfo
	for _, k := range a {
		for _, s := range n[k] {
			if s != 0 && s != len(n) {
				// fmt.Printf("%d, %d\n", s, k)
				importantVersions = append(importantVersions, s)
				if len(importantVersions) > 10 {
					sort.Ints(importantVersions)
					for _, nn := range importantVersions {
						r = append(r, versionsInfo{p.Timestamps[nn], nn})
					}
					return r, totalTime
				}
			}
		}
	}
	sort.Ints(importantVersions)
	for _, nn := range importantVersions {
		r = append(r, versionsInfo{p.Timestamps[nn], nn})
	}
	return r, totalTime
}

func rebuildTextsToDiffN(p WikiData, n int) string {
	dmp := diffmatchpatch.New()
	lastText := ""
	for i, diff := range p.Diffs {
		seq1, _ := dmp.DiffFromDelta(lastText, diff)
		textsLinemode := diffRebuildtexts(seq1)
		rebuilt := textsLinemode[len(textsLinemode)-1]
		if i == n {
			return rebuilt
		}
		lastText = rebuilt
	}
	return "ERROR"
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc prints a random string
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GetLocalIP returns the local ip address
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	bestIP := ""
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && (strings.Contains(ipnet.IP.String(), "192.168.1") || strings.Contains(ipnet.IP.String(), "192.168")) {
				return ipnet.IP.String()
			}
		}
	}
	return bestIP
}
