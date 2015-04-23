package hmutil

import (
	"log"
	"os/exec"
	"strings"

	"os"
	"fmt"
	"path/filepath"
	"bytes"
	"io/ioutil"
	"io"
	"crypto/aes"
  "crypto/cipher"
	"compress/gzip"
	"archive/tar"
)

func System(cmd string) ([]byte, error) {
	log.Println(cmd)
	return exec.Command("sh", "-c", cmd).CombinedOutput()
}

func ReplaceVars(str string, replacements map[string]string) string {
	for from, to := range replacements {
		str = strings.Replace(str, from, to, -1)
	}
	return str
}

func Tar(dir string, files []string, tw *tar.Writer, prevdir string) {
	if len(files) == 1 && files[0] == "*" {
		files = []string{}
		d, _err := ioutil.ReadDir(dir)
		handleError(_err)
		for _, e := range d {
			files = append(files, e.Name())
		}
	}

	for _, file := range files {
		f, _err := os.Open(filepath.Join(dir, file))
		if _err != nil {
			continue
		}
		s, _err := f.Stat()
		if s.IsDir() {
			Tar(filepath.Join(dir, file), []string{"*"}, tw, file)
		} else {
			handleError(_err)
			// s, _ := f.Stat()

			header := &tar.Header{
				Name: filepath.Join(prevdir, s.Name()),
				Size: s.Size(),
				Mode: 0777,
			}

			if _err = tw.WriteHeader(header); _err != nil {
				handleError(_err)
			}

			buffer := make([]byte, s.Size())
			buffer, _err = ioutil.ReadFile(filepath.Join(dir, s.Name())) // file.Read(buffer)

			if _, _err = tw.Write(buffer); _err != nil {
				handleError(_err)
			}
		}
	}

}

// func Gzip(pathToFile string, buffer *bytes.Buffer) {
// 	gzFile, _err := os.Create(pathToFile)
// 	handleError(_err)

// 	gzWriter := gzip.NewWriter(gzFile)
// 	_, _err = gzWriter.Write(buffer.Bytes())
// 	handleError(_err)

// 	gzWriter.Close()
// }

func Gzip(buffer *bytes.Buffer) bytes.Buffer {
  var gzFile bytes.Buffer
  gzWriter := gzip.NewWriter(&gzFile)

  _, _err := gzWriter.Write(buffer.Bytes())
  handleError(_err)
  gzWriter.Close()

  return gzFile
}

func Encode(buffer *bytes.Buffer, encodeKey []byte) bytes.Buffer {
  var outbuffer bytes.Buffer

  block, _err := aes.NewCipher(encodeKey)
  handleError(_err)

  var iv [aes.BlockSize]byte
  stream := cipher.NewOFB(block, iv[:])

  writer := &cipher.StreamWriter{S: stream, W: &outbuffer}

  if _, _err = io.Copy(writer, buffer); _err != nil {
    panic(_err)
  }

  writer.Close()

  return outbuffer
}

func WriteToFile(path string, buffer bytes.Buffer) {
  file, _err := os.Create(path)
  handleError(_err)
  file.Write(buffer.Bytes())
  file.Close()
}

func PackAndCompress(dir string, files []string, outputFile string, key []byte, encrypt bool) {
	outdir, _ := filepath.Split(outputFile)
	_err := os.MkdirAll(outdir, 0777)
	handleError(_err)

	var tarFile bytes.Buffer

	tarWriter := tar.NewWriter(&tarFile)
	Tar(dir, files, tarWriter, "")
	tarWriter.Close()
	gzipedBuffer := Gzip(&tarFile)

	if encrypt {
		encryptedBuffer := Encode(&gzipedBuffer, key)
		WriteToFile(outputFile + ".encrypted", encryptedBuffer)
	} else {
		WriteToFile(outputFile, gzipedBuffer)
	}
}

func handleError(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

/*
func ErrString(err error) (s *string) {
	if err != nil {
		tmp := err.Error()
		s = &tmp
	}
	return
}
*/
