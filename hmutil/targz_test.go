package hmutil

import (
  "testing"
  "os"
  "strconv"
  "sync"
)

const (
  testDirPath = "/tmp/testpath/"
)

var wg *sync.WaitGroup = new(sync.WaitGroup)

func before() {
  os.MkdirAll(testDirPath, 0777)
  os.Mkdir(testDirPath + "1", 0777)
  os.Mkdir(testDirPath + "2", 0777)

  f, _ := os.Create(testDirPath + "aa.txt")
  f.Write([]byte("anisoudhioyuhgw uydgaistu dfaystd fatsd gcvbxzkhxjna lssahuosda iwsdy agsdaitsfda usyfd auyfsvdzkjhsbkahsbdalhk"))
  f.Close()

  i := 0
  for i < 5 {
    go func(n int) {
      f1, _ := os.Create(testDirPath + "1/" + strconv.Itoa(n) + ".txt")
      f1.Write([]byte("aiudo isudipus disdhpius dsuijdhap isjbdhaiw uhdpa uoshdap usodhai sjdbaishdb"))
      f1.Close()
      wg.Done()
    }(i)
    go func(n int) {
      f2, _ := os.Create(testDirPath + "2/" + strconv.Itoa(n) + ".txt")
      f2.Write([]byte("aiudo isudipus disdhpius dsuijdhap isjbdhaiw uhdpa uoshdap usodhai sjdbaishdb"))
      f2.Close()
      wg.Done()
    }(i)
    i++
  }
}

func after() {
  os.RemoveAll(testDirPath)
}

func TestTarGZ(t *testing.T) {
  wg.Add(10)
  before()
  _err := PackAndCompress(testDirPath,[]string{"*"}, testDirPath + "result.tar.gz", []byte("qqqqwwwweeeensdu"), true)
  if _err != nil {
    t.Error(_err)
  }
  after()
  wg.Wait()
}
