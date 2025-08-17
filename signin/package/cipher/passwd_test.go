package cipher

import (
	"fmt"
	"log"
	"testing"
)

func Test_Passwd(t *testing.T) {

	passwd := "hello world!"
	encodeHash, err := GenerateHashPasswd(passwd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(encodeHash)

	ok, err := CompareHashPasswd(passwd, encodeHash)
	if err != nil {
		log.Fatal(err)
	}

	if !ok {
		t.Errorf("CompareHashPasswd fail")
	} else {
		t.Logf("CompareHashPasswd pass")
	}

}
