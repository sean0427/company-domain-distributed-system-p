package config

import (
	"errors"
	"fmt"
	"os"
)

func GetPGCompany() (string, error) {
	if v, found := os.LookupEnv("POSTGRES_USER"); found {
		return v, nil
	} else {
		return "", errors.New("POSTGES_USER is not set")
	}
}

func GetPGPassword() (string, error) {
	if v, found := os.LookupEnv("POSTGRES_PASSWORD"); found {
		return v, nil
	} else {
		return "", errors.New("POSTGRES_PASSWORD is not set")
	}
}

func GetPGDB() (string, error) {
	if v, found := os.LookupEnv("POSTGRES_DB"); found {
		return v, nil
	} else {
		return "", errors.New("POSTGRES_DB is not set")
	}
}

func GetPGAddress() (string, error) {
	if v, found := os.LookupEnv("POSTGRES_ADDRESS"); found {
		return v, nil
	} else {
		return "", errors.New("POSTGRES_ADDRESS is not set")
	}
}
func GetPGPort() (string, error) {
	if v, found := os.LookupEnv("POSTGRES_PORT"); found {
		return v, nil
	} else {
		return "", errors.New("POSTGRES_ADDRESS is not set")
	}
}

func GetPostgresDNS() (string, error) {
	company, err := GetPGCompany()
	if err != nil {
		return "", err
	}

	password, err := GetPGPassword()
	if err != nil {
		return "", err
	}

	db, err := GetPGDB()
	if err != nil {
		return "", err
	}

	address, err := GetPGAddress()
	if err != nil {
		return "", err
	}

	port, err := GetPGPort()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", company, password, address, port, db), nil
}
