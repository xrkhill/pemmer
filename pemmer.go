package main

import (
	"bufio"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	dashes = "-----"
	begin  = "-----BEGIN "
	end    = "-----END "
)

var labels = map[string]string{
	"cert":       "CERTIFICATE",
	"crl":        "X509 CRL",
	"csr":        "CERTIFICATE REQUEST",
	"pkcs7":      "PKCS7",
	"cms":        "CMS",
	"privkey":    "PRIVATE KEY",
	"encprivkey": "ENCRYPTED PRIVATE KEY",
	"attrcert":   "ATTRIBUTE CERTIFICATE",
	"pubkey":     "PUBLIC KEY",
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	to := flag.String("to", "", "convert to [blob | pem]")

	labelKeys := make([]string, len(labels))
	i := 0
	for k := range labels {
		labelKeys[i] = k
		i++
	}

	labelUsage := fmt.Sprintf("PEM type label, one of: %s", strings.Join(labelKeys, ", "))
	label := flag.String("label", "cert", labelUsage)

	flag.Parse()

	var labelType string
	var ok bool
	if *to == "pem" {
		labelType, ok = labels[*label]
		if !ok {
			log.Fatalf("invalid label type: %s", *label)
		}
	}

	switch *to {
	case "blob":
		toBlob(reader, writer)
	case "pem":
		toPEM(reader, writer, labelType)
	default:
		log.Fatalf("invalid to, must be one of 'blob' or 'pem'")
	}

	err := writer.Flush()
	if err != nil {
		log.Fatalf("unable to flush contents to buffer: %v", err)
	}
}

func parseFlags() {
}

func toPEM(reader *bufio.Reader, writer *bufio.Writer, labelType string) {
	preamble := begin + labelType + dashes + "\n"
	_, err := writer.WriteString(preamble)
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

	postamble := "\n" + end + labelType + dashes + "\n"
	_, err = writer.WriteString(postamble)
	if err != nil {
		log.Fatalf("unable to write postamble to buffer: %v", err)
	}
}

func toBlob(reader *bufio.Reader, writer *bufio.Writer) {
	// check for properly PEM encoded content and fail if not parsable
	pemBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatal("unable to read bytes from buffer: %v", err)
	}
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock == nil {
		log.Fatal("unable to decode PEM data")
	}

	// rewind cursor to beginning of buffer
	os.Stdin.Seek(0, io.SeekStart)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if !strings.HasPrefix(text, begin) && !strings.HasPrefix(text, end) {
			newText := strings.ReplaceAll(text, "\n", "")

			_, err = writer.WriteString(newText)
			if err != nil {
				log.Fatalf("unable to write string to buffer: %v", err)
			}
		}
	}
}
