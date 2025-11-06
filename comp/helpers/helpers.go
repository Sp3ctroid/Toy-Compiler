package helpers

import "strconv"

func String_plus(s string) string {
	s_int, _ := strconv.Atoi(s)
	s_int++
	s = strconv.Itoa(s_int)
	return s
}

func Is_digit(char byte) bool {
	if char >= '0' && char <= '9' {
		return true
	}

	return false
}

func Is_comparative(char byte) bool {
	if char == '=' || char == '>' || char == '<' {
		return true
	}

	return false
}

func Is_char(char byte) bool {
	if char >= 'A' && char <= 'Z' {
		return true
	}

	if char >= 'a' && char <= 'z' {
		return true
	}

	return false
}

func Is_boolean(char byte) bool {
	if char == '&' || char == '|' {
		return true
	}
	return false
}

func Is_special(char byte) bool {
	switch char {
	case ';':
		return true
	case ',':
		return true
	case '(':
		return true
	case ')':
		return true
	case '{':
		return true
	case '}':
		return true
	case '=':
		return true
	case '+':
		return true
	case '-':
		return true
	case '*':
		return true
	case '/':
		return true
	case '\n':
		return true
	case '\t':
		return true
	case '&':
		return true
	case '|':
		return true
	}
	return false
}
