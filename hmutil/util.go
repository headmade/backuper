package hmutil

import (
	// "log"
	// "os/exec"
	"strings"
	"fmt"
	"os"
	"path/filepath"
	"bytes"
	"io/ioutil"
	"io"
	"crypto/aes"
	"crypto/cipher"
	"compress/gzip"
	"archive/tar"

	"github.com/rlmcpherson/s3gof3r"
	"github.com/Twizty/scp-go/scp"
	"github.com/dutchcoders/goftp"
)

// func System(cmd string) ([]byte, error) {
// 	log.Println(cmd)
// 	return exec.Command("sh", "-c", cmd).CombinedOutput()
// }

func ReplaceVars(str string, replacements map[string]string) string {
	for from, to := range replacements {
		str = strings.Replace(str, from, to, -1)
	}
	return str
}

func find(prefix string, arr []string, elem string) int {
	for i, e := range arr {
		if filepath.Join(prefix, e) == elem {
			return i
		}
	}
	return -1
}

func removeFromSlice(prefix string, arr []string, elem string) ([]string, bool) {
	index := find(prefix, arr, elem)

	if index == -1 {
		return arr, false
	} else {
		return append(arr[:index], arr[index+1:]...), true
	}
}

func removeExcludedFiles(dir string, files, excludeFiles []string) {
	var changed bool
	var deleted []string = []string{}

	for _, e := range excludeFiles {
		files, changed = removeFromSlice(dir, files, e)
		if changed {
			deleted = append(deleted, e)
		}
	}
	for _, e := range deleted {
		removeFromSlice("", excludeFiles, e)
	}
}

func Tar(dir string, files, excludeFiles []string, tw *tar.Writer) error {
	if len(files) == 1 && files[0] == "*" {
		files = []string{}
		d, _err := ioutil.ReadDir(dir)
		
		if _err != nil {
			return _err
		}

		for _, e := range d {
			files = append(files, e.Name())
		}
	}

	// for _, e := range excludeFiles {
	// 	fmt.Println(filpath.Join(dir, e))
	// 	files = removeFromSlice(files, e)
	// 	excludeFiles = removeFromSlice(excludeFiles, e)
	// }

	removeExcludedFiles(dir, files, excludeFiles)

	for _, file := range files {
		f, _err := os.Open(filepath.Join(dir, file))
		if _err != nil {
			continue
		}

		s, _err := f.Stat()

		if s.IsDir() {
			Tar(filepath.Join(dir, file), []string{"*"}, excludeFiles, tw)
		} else {
			if _err != nil {
				return _err
			}

			header := &tar.Header{
				Name: filepath.Join(dir, s.Name()),
				Size: s.Size(),
				Mode: 0777,
			}

			if _err = tw.WriteHeader(header); _err != nil {
				return _err
			}

			buffer := make([]byte, s.Size())
			buffer, _err = ioutil.ReadFile(filepath.Join(dir, s.Name())) // file.Read(buffer)
			if _err != nil {
				return _err
			}

			if _, _err = tw.Write(buffer); _err != nil {
				return _err
			}
		}
	}

	return nil
}

// func Gzip(pathToFile string, buffer *bytes.Buffer) {
// 	gzFile, _err := os.Create(pathToFile)
// 	handleError(_err)

// 	gzWriter := gzip.NewWriter(gzFile)
// 	_, _err = gzWriter.Write(buffer.Bytes())
// 	handleError(_err)

// 	gzWriter.Close()
// }

func Gzip(buffer *bytes.Buffer) (bytes.Buffer, error) {
	var gzFile bytes.Buffer
	gzWriter := gzip.NewWriter(&gzFile)
	_, _err := gzWriter.Write(buffer.Bytes())
	gzWriter.Flush()
	gzWriter.Close()
	if _err != nil {
		return bytes.Buffer{}, _err
	}
	
	return gzFile, nil
}

func Encode(buffer *bytes.Buffer, encodeKey []byte) (bytes.Buffer, error) {
	var outbuffer bytes.Buffer

	block, _err := aes.NewCipher(encodeKey)
	if _err != nil {
		return bytes.Buffer{}, _err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewOFB(block, iv[:])

	writer := &cipher.StreamWriter{S: stream, W: &outbuffer}
	defer writer.Close()
	if _, _err = io.Copy(writer, buffer); _err != nil {
		return bytes.Buffer{}, _err
	}

	return outbuffer, nil
}

func WriteToFile(path string, buffer bytes.Buffer) error {
	file, _err := os.Create(path)
	defer file.Close()
	if _err != nil {
		return _err
	}
	file.Write(buffer.Bytes())
	return nil
}

func PackAndCompress(dir string, files, excludeFiles []string, outputFile string, key []byte, encrypt bool) error {
	outdir, _ := filepath.Split(outputFile)
	_err := os.MkdirAll(outdir, 0777)
	if _err != nil {
		return _err
	}

	var tarFile bytes.Buffer

	tarWriter := tar.NewWriter(&tarFile)
	_err = Tar(dir, files, excludeFiles, tarWriter)

	if _err != nil {
		return _err
	}

	tarWriter.Close()
	gzipedBuffer, _err := Gzip(&tarFile)

	if _err != nil {
		return _err
	}

	if encrypt {
		encryptedBuffer, _err := Encode(&gzipedBuffer, key)
		if _err != nil {
			return _err
		}
		return WriteToFile(outputFile + ".encrypted", encryptedBuffer)
	} else {
		return WriteToFile(outputFile, gzipedBuffer)
	}
}

func handleError(e error) {
	if e != nil {
		panic(e)
	}
}

func UploadToS3(k s3gof3r.Keys, bucketName string, pathToFile string) error {
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(bucketName)

	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}

	stats, _ := file.Stat()

	w, err := b.PutWriter(stats.Name(), nil, nil)
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, file); err != nil { // Copy into S3
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return nil
}

func DownloadFromS3(k s3gof3r.Keys, bucketName, filename string) (bytes.Buffer, error) {
	s3 := s3gof3r.New("", k)
	b := s3.Bucket(bucketName)

	r, _, err := b.GetReader(filename, nil)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// outputFile, err := os.Create(filepath.Join(outputDirPath, filename))
	// if err != nil {
	// 	return nil, err
	// }

	var output bytes.Buffer

	if _, err = io.Copy(&output, r); err != nil {
		return bytes.Buffer{}, err
	}

	err = r.Close()
	if err != nil {
		return bytes.Buffer{}, err
	}

	return output, nil
}

func SSHUploader(port int, id_rsa, host, user, srcFile, dstFile string) *scp.SecureCopier {
	return scp.NewSecureCopier(port, false, false, false, true, false, false, id_rsa, "", "", srcFile, host, user, dstFile)
}

func SSHDownloader(port int, id_rsa, host, user, srcFile, dstFile string) *scp.SecureCopier {
	return scp.NewSecureCopier(port, false, false, false, true, false, false, id_rsa, host, user, srcFile, "", "", dstFile)
}

func SSHExec(ssh *scp.SecureCopier) (bytes.Buffer, error) {
	var o, e, result bytes.Buffer
	err, _ := ssh.Exec(nil, &o, &e, &result)
	return result, err
}

func prepareFTPConn(port int, host, login, password string) (*goftp.FTP, error) {
	var ftp *goftp.FTP
	var err error

	if ftp, err = goftp.Connect(fmt.Sprintf("%s:%d", host, port)); err != nil {
		return nil, err
	}

	if err = ftp.Login(login, password); err != nil {
		return nil, err
	}

	return ftp, nil
}

func FTPUpload(port int, host, login, password, pathToFile, pathToRemoteFile string) error {
	ftp, err := prepareFTPConn(port, host, login, password)
	if err != nil {
		return err
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}

	if err = ftp.Stor(pathToRemoteFile, file); err != nil {
		return err
	}

	return ftp.Quit()
}

func FTPDownload(port int, host, login, password, pathToRemoteFile string) (bytes.Buffer, error) {
	ftp, err := prepareFTPConn(port, host, login, password)
	if err != nil {
		return bytes.Buffer{}, err
	}

	var buffer bytes.Buffer

	_, err = ftp.Retr(pathToRemoteFile, func(r io.Reader) error {
		// file, _err := os.Create(pathToFile)
		// if _err != nil {
		// 	return _err
		// }

		if _, _err := io.Copy(&buffer, r); _err != nil {
			return _err
		}

		return nil
	})

	if err != nil {
		return bytes.Buffer{}, err
	}

	if err = ftp.Quit(); err != nil {
		return buffer, err
	}

	return buffer, nil
}

func Decode(buffer *bytes.Buffer, key []byte) (bytes.Buffer, error) {
	block, err := aes.NewCipher(key)
  if err != nil {
    return bytes.Buffer{}, err
  }

  var iv [aes.BlockSize]byte
  stream := cipher.NewOFB(block, iv[:])

  var outbuffer bytes.Buffer
  reader := &cipher.StreamReader{S: stream, R: buffer}

  if _, err = io.Copy(&outbuffer, reader); err != nil {
    return bytes.Buffer{}, err
  }

  return outbuffer, nil
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
