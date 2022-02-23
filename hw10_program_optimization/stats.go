package hw10programoptimization

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

var (
	ErrReaderEmpty = errors.New("невалидный reader")
	ErrDomainEmpty = errors.New("невалидный domain")
	ErrInvalidJson = errors.New("невалидный json")
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	if r == nil {
		return result, ErrReaderEmpty
	}

	if len(domain) == 0 {
		return result, ErrDomainEmpty
	}
	user := User{}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if !json.Valid(scanner.Bytes()) {
			return nil, ErrInvalidJson
		}
		if !bytes.Contains(scanner.Bytes(), []byte(domain)) {
			continue
		}

		if err := json.Unmarshal(scanner.Bytes(), &user); err == nil {
			userDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			if strings.HasSuffix(userDomain, domain) {
				result[userDomain]++
			}
		}
	}

	return result, nil
}
