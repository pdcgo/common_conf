package pdc_common

import (
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrAccountNotFound error = errors.New("account not found")

type AccountData struct {
	Username      string
	Password      string
	Email         string
	EmailPassword string
	GroupName     string
}

type AccountRepo struct {
	BaseDir  string
	Accounts []*AccountData
}

func (ac *AccountRepo) Load() {

	globpath := filepath.Join(ac.BaseDir, "*.txt")

	files, _ := filepath.Glob(globpath)
	for _, file := range files {
		ac.Accounts = append(ac.Accounts, ReadCsvFile(file)...)
	}

}

func (ac *AccountRepo) Get(name string) (*AccountData, error) {
	for _, akun := range ac.Accounts {
		if akun.Username == name {
			return akun, nil
		}
	}

	return nil, ErrAccountNotFound
}

func ReadCsvFile(file string) []*AccountData {
	var hasil []*AccountData

	csvf, _ := os.Open(file)
	defer csvf.Close()

	// getting groupname
	groupName := strings.ReplaceAll(file, "\\", "/")
	gnames := strings.Split(groupName, "/")
	groupName = gnames[len(gnames)-1]
	groupName = strings.ReplaceAll(groupName, ".txt", "")

	// reading csv
	csvreader := csv.NewReader(csvf)
	csvreader.Comma = '|'
	csvreader.FieldsPerRecord = -1

	csvdata, err := csvreader.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, each := range csvdata {
		lines := make([]string, 5)
		copy(lines, each)
		acdata := AccountData{
			Username:      lines[0],
			Password:      lines[1],
			Email:         lines[2],
			EmailPassword: lines[3],
			GroupName:     groupName,
		}
		hasil = append(hasil, &acdata)
	}

	return hasil
}

func NewAccountRepo(base string) *AccountRepo {
	repo := &AccountRepo{
		BaseDir:  base,
		Accounts: []*AccountData{},
	}

	repo.Load()

	return repo
}
