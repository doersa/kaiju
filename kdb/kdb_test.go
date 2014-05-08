package kdb

import (
    "bytes"
    //"testing"
    "encoding/binary"
    //"encoding/hex"
    "crypto/sha256"
    "os"
    "testing"
    "fmt"
    "path/filepath"
    "github.com/oxfeeefeee/kaiju/config"
)

func cookUint32(key uint32, value uint32)([]byte, []byte) {
    kbuf, vbuf := make([]byte, 6, 6), make([]byte, 4, 4)
    binary.LittleEndian.PutUint32(kbuf[:], key)
    hash := sha256.Sum256(kbuf[0:4])
    copy(kbuf[:], hash[0:6])
    clearMaskBits(kbuf)
    binary.LittleEndian.PutUint32(vbuf, value)  
    return kbuf, vbuf
}

func writeUint32(t *testing.T, db *KDB, key uint32, value uint32) {
    kbuf, vbuf := cookUint32(key, value)
    db.addRecord(kbuf, vbuf)
}

func testUint32(t *testing.T, db *KDB, key uint32, value uint32) {
    kbuf, vbuf := cookUint32(key, value)
    v, getErr := db.getRecord(kbuf, 1)
    if getErr != nil {
        t.Errorf(fmt.Sprintf("Failed to getRecord KDB: %s", getErr))
    }
    if bytes.Compare(vbuf, v) != 0 {
        t.Errorf("Did not get what I set %v %v", vbuf, v)
    }
}

func TestKDB(t *testing.T) {
    err := config.ReadJsonConfigFile()
    if err != nil {
        t.Errorf(fmt.Sprintf("Failed to read config file: %s", err))
    }

    cfg := config.GetConfig()
    path := filepath.Join(config.GetConfigFileDir(), cfg.DBTempDir)
    os.MkdirAll(path, os.ModePerm)

    path = filepath.Join(path, "testdb.dat")
    exists, _ := fileExists(path)
    if exists {
        fmt.Printf("File already there: %s", path)
    }

    fmt.Printf("File Path: %s\n", path)
    f, openErr := os.Create(path)
    if openErr != nil {
        t.Errorf(fmt.Sprintf("Failed to create file: %s", openErr))
    }

    capacity := uint32(1 * 1024 * 1024)
    db, dberr := New(capacity, f)
    if dberr != nil {
        t.Errorf(fmt.Sprintf("Failed to create KDB: %s", dberr))
    }

    for i:=uint32(1); i < capacity; i++ {
        writeUint32(t, db, uint32(i), uint32(i))  
    }
    for i:=uint32(1); i < capacity; i++ {
        testUint32(t, db, uint32(i), uint32(i))  
    }
    fmt.Printf("db: %s\n", db)

    if closeErr := f.Close(); closeErr != nil {
        t.Errorf(fmt.Sprintf("Error closing file: %s", closeErr))
    }
}

func fileExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}