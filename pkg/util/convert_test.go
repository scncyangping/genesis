package util

import (
	"fmt"
	"testing"
)

func TestCamelize(t *testing.T) {
	fmt.Println(Camelize("sys_user", false))
	fmt.Println(Camelize("sys_user", false))
}

func TestAESEncrypt(t *testing.T) {
	data := "123456"
	if encrypt, err := AESEncrypt([]byte(data)); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%s Encrypt => %s", data, encrypt)
	}
}

func TestAESDecrypt(t *testing.T) {
	data := ""
	if decrypt, err := AESDecrypt(data); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%s decrypt => %s", data, decrypt)
	}
}
