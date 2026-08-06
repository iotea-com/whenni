package gruenthttputil

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"strings"
	"text/template"
)

func BuildEmail(from string, to []string, subject, plainBodyTemplate, htmlBodyTemplate string, bodyData map[string]string) ([]byte, error) {
	// Generate HTML body
	htmlTmpl := template.Must(template.New("html").Parse(htmlBodyTemplate))
	htmlBody := bytes.Buffer{}
	htmlTmpl.Execute(&htmlBody, bodyData)

	// Generate plain text body
	plainTmpl := template.Must(template.New("plain").Parse(plainBodyTemplate))
	plainBody := bytes.Buffer{}
	plainTmpl.Execute(&plainBody, bodyData)

	// Build multipart message
	joinedToAddresses := strings.Join(to, ",")
	return buildMultipartMessage(from, joinedToAddresses, subject, plainBody.String(), htmlBody.String())
}

// BuildMultipartMessage builds a multipart message with a plain text and HTML part.
// It returns the message as a byte slice and an error if the message cannot be built.
func buildMultipartMessage(from, to, subject, plainText, htmlText string) ([]byte, error) {
	var buf bytes.Buffer

	// Write email headers
	buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")

	boundary := "BOUNDARY123456"
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\r\n", boundary))
	buf.WriteString("\r\n") // End headers

	// Start multipart sections
	// Plain text part
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	qpWriter := quotedprintable.NewWriter(&buf)
	qpWriter.Write([]byte(plainText))
	qpWriter.Close()

	buf.WriteString("\r\n")

	// HTML part
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/html; charset=utf-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	qpWriter = quotedprintable.NewWriter(&buf)
	qpWriter.Write([]byte(htmlText))
	qpWriter.Close()

	buf.WriteString("\r\n")

	// End boundary
	buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return buf.Bytes(), nil
}
