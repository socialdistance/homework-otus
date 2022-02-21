package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
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

		user := &User{}
		err := easyjson.Unmarshal([]byte(res), user)
		if err != nil {
			return nil, err
		}

		domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		if !strings.Contains(domain, domainStatDot) {
			continue
		}

		domainStat[domain]++
	}

	return domainStat, nil
}

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return GetDomainStatNew(r, domain)
	// u, err := getUsers(r)
	// if err != nil {
	//	return nil, fmt.Errorf("get users error: %w", err)
	// }
	// return countDomains(u, domain)
}

// type users [100_000]User
//
// func getUsers(r io.Reader) (result users, err error) {
//	content, err := ioutil.ReadAll(r)
//	if err != nil {
//		return
//	}
//
//	lines := strings.Split(string(content), "\n")
//	for i, line := range lines {
//		var user User
//		if err = json.Unmarshal([]byte(line), &user); err != nil {
//			return
//		}
//		result[i] = user
//	}
//	return
//}
//
// func countDomains(u users, domain string) (DomainStat, error) {
//	result := make(DomainStat)
//
//	for _, user := range u {
//		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
//		if err != nil {
//			return nil, err
//		}
//
//		if matched {
//			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
//			num++
//			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
//		}
//	}
//	return result, nil
//}
