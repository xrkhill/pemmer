# Pemmer

Pemmer is a command line utility for formatting base64 encoded PEM blob strings
that are missing line breaks and `-----BEGIN XXX-----...-----END XXX-----`
preamble and postamble into into correctly formatted PEM files.

This can be useful when copy pasting PEM content from word processing documents,
or when you need to make API requests with base64 encoded PKIX file contents
embedded in JSON.

## Usage

To convert from unformatted blob, run pemmer like this:

	$ ./pemmer -to pem < unformatted-blob.txt > formatted.pem

To convert PEM content to a blob string, run pemmer like this:

	$ ./pemmer -to blob < formatted.pem > unformatted-blob.txt

## References

[RFC 1421](https://tools.ietf.org/html/rfc1421) Privacy Enhancement for Internet
Electronic Mail: Part I: Message Encryption and Authentication Procedures

[RFC 7468](https://tools.ietf.org/html/rfc7468) Textual Encodings of PKIX, PKCS,
and CMS Structures

[Package pem](https://golang.org/pkg/encoding/pem/) Golang.org
