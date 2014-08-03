package main_test

import (
	"os"
	"path"
	"testing"

	. "github.com/marciotrindade/clean_csv"
	"github.com/stretchr/testify/assert"
)

func TestReadCsv(t *testing.T) {
	assert := assert.New(t)

	fileName := "./sample.csv"
	header := []string{"email"}
	email1 := []string{"email1@test.com"}
	email2 := []string{"email2@test.com"}
	email3 := []string{"email3@test.com"}
	content := [][]string{header, email1, email2, email3}

	assert.Equal(ReadCsv(fileName), content)
}

func TestClearEmail(t *testing.T) {
	assert := assert.New(t)

	email := "MárcioTrindade@test.com"

	assert.Equal(ClearEmail(email), "marciotrindade@test.com")
}

func TestValidEmail(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(ValidEmail("MárcioTrindade@test.com"), false)
	assert.Equal(ValidEmail("marciotrindade@hotmail"), false)
	assert.Equal(ValidEmail("marciotrindade@test.com"), true)
	assert.Equal(ValidEmail("marciotrindade@test.com.br"), true)
	assert.Equal(ValidEmail("marciotrindade@hotmail.com"), true)
}

func TestCreateFolder(t *testing.T) {
	assert := assert.New(t)
	path := path.Dir(os.Args[0]) + "/output"

	_, err := os.Stat(path)
	assert.Equal(os.IsNotExist(err), true)

	CreateFolder(path)

	_, err = os.Stat(path)
	assert.Equal(os.IsNotExist(err), false)

	// after remove the created folder
	os.RemoveAll(path)
	_, err = os.Stat(path)
	assert.Equal(os.IsNotExist(err), true)
}
