package certificatesGet

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
)

func respond(request *ioteahttp.Request[Input], output *Output) {
	request.Span.AddEvent("respond")

	// Validate that all of the assets in the bundle are present
	if output.Certificate != nil || output.PrivateKey != nil || output.CaCert != nil {
		errorResponse := ioteahttp.NewErrorResponse([]any{"Missing certificate, private key, or CA certificate for the certificate bundle."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		request.FiberContext.Status(fiber.StatusBadRequest).JSON(errorResponseJson)
	}

	// Create the zip file in memory
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	if err := addFileToZip(zipWriter, "certificate.crt", *output.Certificate); err != nil {
		errorResponse := ioteahttp.NewErrorResponse([]any{"Could not add the certificate to the certificate bundle."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		request.FiberContext.Status(fiber.StatusInternalServerError).JSON(errorResponseJson)
		return
	}

	if err := addFileToZip(zipWriter, "private.key", *output.PrivateKey); err != nil {
		errorResponse := ioteahttp.NewErrorResponse([]any{"Could not add the private key to the certificate bundle."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		request.FiberContext.Status(fiber.StatusInternalServerError).JSON(errorResponseJson)
		return
	}

	if err := addFileToZip(zipWriter, "ca.crt", *output.CaCert); err != nil {
		errorResponse := ioteahttp.NewErrorResponse([]any{"Could not add the CA certificate to the certificate bundle."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		request.FiberContext.Status(fiber.StatusInternalServerError).JSON(errorResponseJson)
		return
	}

	if err := zipWriter.Close(); err != nil {
		errorResponse := ioteahttp.NewErrorResponse([]any{"Could not create the certificate bundle."})
		errorResponseJson, _ := errorResponse.MarshalJson()
		request.FiberContext.Status(fiber.StatusInternalServerError).JSON(errorResponseJson)
		return
	}

	// Set the response headers
	zipFilename := fmt.Sprintf("%s.zip", strings.ReplaceAll(request.Input.CertificateId, ":", ""))
	request.FiberContext.Response().Header.Set(fiber.HeaderContentType, "application/zip")
	request.FiberContext.Response().Header.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", zipFilename))
	request.FiberContext.Response().Header.Set(fiber.HeaderAccessControlExposeHeaders, "Content-Disposition")
	// Respond
	request.FiberContext.Status(fiber.StatusOK).SendStream(bytes.NewReader(buf.Bytes()), len(buf.Bytes()))
}

func addFileToZip(zipWriter *zip.Writer, filename, content string) error {
	// Create a new file in the zip archive
	fileWriter, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	// Write the content to the file in the zip archive
	_, err = io.WriteString(fileWriter, content)
	return err
}
