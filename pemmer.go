package main

import (
	"bufio"
	"encoding/pem"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	preambleDashes = "-----"
	beginCert      = "-----BEGIN CERTIFICATE-----"
	endCert        = "-----END CERTIFICATE-----"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	to := flag.String("to", "", "convert to [blob | pem]")
	flag.Parse()

	switch *to {
	case "blob":
		toBlob(reader, writer)
	case "pem":
		toPEM(reader, writer)
	default:
		log.Fatalf("invalid to, must be one of 'blob' or 'pem'")
	}

	err := writer.Flush()
	if err != nil {
		log.Fatalf("unable to flush contents to buffer: %v", err)
	}
}

func toPEM(reader *bufio.Reader, writer *bufio.Writer) {
	_, err := writer.WriteString(beginCert + "\n")
	if err != nil {
		log.Fatalf("unable to write preamble to buffer: %v", err)
	}

	i := 0
	for {
		b, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		if i > 1 && i%64 == 0 {
			_, err = writer.WriteString("\n")
			if err != nil {
				log.Fatalf("unable to write newline to buffer: %v", err)
			}
		}

		err = writer.WriteByte(b)
		if err != nil {
			log.Fatalf("unable to write content byte to buffer: %v", err)
		}

		i++
	}

	_, err = writer.WriteString("\n" + endCert + "\n")
	if err != nil {
		log.Fatalf("unable to write postamble to buffer: %v", err)
	}
}

func toBlob(reader *bufio.Reader, writer *bufio.Writer) {
	pemBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("unable to read bytes from buffer: %v", err)
	}

	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		log.Fatal("unable to decode PEM data")
	}

	os.Stdin.Seek(0, io.SeekStart)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if text != (beginCert+"\n") && text != (endCert+"\n") {
			newText := strings.ReplaceAll(text, "\n", "")

			_, err = writer.WriteString(newText)
			if err != nil {
				log.Fatalf("unable to write string to buffer: %v", err)
			}
		}
	}
}
