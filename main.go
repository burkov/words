package main

import (
	"bufio"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)


func readIntParam(name string, deflt int, w http.ResponseWriter, r *http.Request) (int, error) {
	raw := r.URL.Query().Get(name)
	var result = deflt
	var err error
	if raw != "" {
		result, err = strconv.Atoi(raw)
		if err != nil || result < 1 || result > 100 {
			result = 0
			w.WriteHeader(400)
			_, err := fmt.Fprint(w, "wrong param\n")
			if err != nil {
				log.Panicln(err)
			}
		}
	}
	return result, err
}

func Words(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	lines, err := readIntParam("lines", 20, w, r)
	perLine, err := readIntParam("perLine", 10, w, r)
	if lines == 0 || perLine == 0 {
		return
	}
	_, err = fmt.Fprintf(w, "attempt %d\n", trackNumberOfCalls())
	if err != nil {
		log.Panicln(err)
	}
	text := paragraph(lines, perLine)
	_, err = fmt.Fprint(w, text)
	if err != nil {
		log.Panicln(err)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	router := httprouter.New()
	router.GET("/words", Words)
	log.Fatal(http.ListenAndServe(":3989", router))
}

func closeOrPanic(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Panicln(err)
	}
}

func readLines(path string) []string {
	var result []string
	fd, err := os.Open(path)
	if err != nil {
		log.Panicln(err)
	}
	defer closeOrPanic(fd)
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })
	return result
}

var medium = readLines("mounted/google-10000-english/google-10000-english-usa-no-swears-medium.txt")
var long = readLines("mounted/google-10000-english/google-10000-english-usa-no-swears-long.txt")

func paragraph(nLines int, nPerLine int) string {
	total := nLines * nPerLine
	words := randomWords(total)
	var sb strings.Builder
	for i := 0; i < nLines; i++ {
		k := i * nPerLine
		sb.WriteString(strings.Join(words[k:k+nPerLine], " "))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomWords(n int) []string {
	rand.Shuffle(len(medium), func(i, j int) { medium[i], medium[j] = medium[j], medium[i] })
	rand.Shuffle(len(long), func(i, j int) { long[i], long[j] = long[j], long[i] })
	result := append(long[:n/2], medium[:n/2]...)
	return result
}

func trackNumberOfCalls() int {
	fd, err := os.OpenFile("mounted/counter", os.O_RDWR|os.O_CREATE, 0644)
	defer closeOrPanic(fd)
	if err != nil {
		log.Panicln(err)
	}
	reader := bufio.NewReader(fd)
	line, err := reader.ReadString('\n')
	var counter = 0
	if err != io.EOF {
		if err != nil {
			log.Panicln(err)
		}
		counter, err = strconv.Atoi(strings.TrimRight(line, "\n"))
		if err != nil {
			log.Panicln(err)
		}
	}
	counter++
	_, err = fd.Seek(0, io.SeekStart)
	if err != nil {
		log.Panicln(err)
	}
	_, err = fd.WriteString(strconv.Itoa(counter) + "\n")
	if err != nil {
		log.Panicln(err)
	}
	return counter - 1
}
