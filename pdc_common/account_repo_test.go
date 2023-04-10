package pdc_common_test

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdcgo/common_conf/pdc_common"
	"github.com/stretchr/testify/assert"
)

func MockAccount() (string, func()) {
	base := "test_repo"

	wd, _ := os.Getwd()

	basepath := filepath.Join(wd, base)
	log.Println(basepath)
	os.MkdirAll(basepath, 0755)

	fname := filepath.Join(basepath, "test.txt")
	csvfile, err := os.Create(fname)
	if err != nil {
		panic(err)
	}

	defer csvfile.Close()

	csvwriter := csv.NewWriter(csvfile)
	csvwriter.Comma = '|'

	defer csvwriter.Flush()

	csvwriter.Write([]string{"user", "pass"})
	csvwriter.Write([]string{"user", "pass", "email"})
	csvwriter.Write([]string{"user", "pass", "email", "emailpas"})
	csvwriter.Write([]string{"gudanggaram", "pass", "email"})

	return basepath, func() {
		os.RemoveAll(basepath)
	}
}

func TestAccountRepo(t *testing.T) {
	base, cancel := MockAccount()
	defer cancel()

	repo := pdc_common.NewAccountRepo(base)

	assert.NotEmpty(t, repo.Accounts)

	t.Run("test finding akun", func(t *testing.T) {
		acdata, err := repo.Get("gudanggaram")

		assert.NotEmpty(t, acdata)
		assert.Nil(t, err)

	})

}
