package hw10programoptimization

import (
	"bufio"
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

func GetDomainStatNew(r io.Reader, domain string) (DomainStat, error) {
	domainStat := make(DomainStat)
	domainStatDot := "." + strings.ToLower(domain)

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		res := scanner.Text()
		if !strings.Contains(res, domain) {
			continue
		}

		var user User

		if err := user.UnmarshalJSON([]byte(res)); err != nil {
			return nil, err
		}

		domainRes := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		if !strings.Contains(domainRes, domainStatDot) {
			continue
		}

		domainStat[domainRes]++
	}

	return domainStat, nil
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return GetDomainStatNew(r, domain)
}
