package couchdb

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
)

// Get mime type from file name.
func mimeType(name string) string {
	ext := filepath.Ext(name)
	return mime.TypeByExtension(ext)
}

// Write JSON to multipart/related.
func writeJSON(document *Document, writer *multipart.Writer, file *os.File) error {
	partHeaders := textproto.MIMEHeader{}
	partHeaders.Set("Content-Type", "application/json")
	part, err := writer.CreatePart(partHeaders)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	path := file.Name()

	// make empty map
	document.Attachments = make(map[string]Attachment)
	attachment := Attachment{
		Follows:     true,
		ContentType: mimeType(path),
		Length:      stat.Size(),
	}
	// add attachment to map
	filename := filepath.Base(path)
	document.Attachments[filename] = attachment

	bytes, err := json.Marshal(document)
	if err != nil {
		return err
	}

	_, err = part.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// Write actual file content to multipart/related.
func writeMultipart(writer *multipart.Writer, file io.Reader) error {
	part, err := writer.CreatePart(textproto.MIMEHeader{})
	if err != nil {
		return err
	}

	// copy file content into multipart message
	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	return nil
}
