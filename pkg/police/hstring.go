package police

import "strings"

const DELIMITER_HSTRING = ":"

func (p *ACL) GetHString() ([]string, []string, []string) {
	sub, act, res := ConverHString(p.Subject), ConverHString(p.Action), ConverHString(p.Resource)
	return sub, act, res
}

func ConverHString(str string) (result []string) {
	if str == "" {
		result = make([]string, 0)
	} else {
		result = strings.Split(str, DELIMITER_HSTRING)
	}
	return
}
