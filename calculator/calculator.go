package calculator

import (
	"bufio"
	"net"
	"os"
	"regexp"
	"strconv"
)

const (
	specialChar = '!'
)

var (
	reExpr = regexp.MustCompile("^(\\d+?)([\\+\\-\\/\\*])(\\d+?)\\!?$")
	// Error é a mensagem padrão de erro
	Error = []byte("0!")
)

// MakeMsg monta a resposta
func MakeMsg(expr []byte) []byte {
	return append(expr, specialChar)
}

// CleanMsg remove !
func CleanMsg(expr []byte) []byte {
	if l := len(expr); l > 0 && expr[l-1] == '!' {
		return expr[:l-1]
	}
	return expr
}

// ReadConn ler o body da mensagem
func ReadConn(conn net.Conn) ([]byte, error) {
	b, e := bufio.NewReader(conn).ReadBytes(specialChar)
	b = CleanMsg(b)
	return b, e
}

// Validate Valida a mensagem
func Validate(expr []byte) bool {
	return reExpr.Match(expr)
}

//CalculateExpression faz a operação matematica
func CalculateExpression(expr []byte) ([]byte, bool) {
	matches := reExpr.FindAllSubmatch(expr, -1)
	if matches == nil {
		return nil, false
	}
	match := matches[0]
	n1, err := strconv.ParseInt(string(match[1]), 10, 0)
	if err != nil {
		return nil, false
	}
	n2, err := strconv.ParseInt(string(match[3]), 10, 0)
	if err != nil {
		return nil, false
	}
	var n int64
	switch string(match[2]) {
	case "+":
		n = n1 + n2
	case "-":
		n = n1 - n2
	case "/":
		n = n1 / n2
	case "*":
		n = n1 * n2
	default:
		return nil, false
	}
	return []byte(strconv.FormatInt(n, 10)), true
}

// EndIfErr lança o erro
func EndIfErr(err error) {
	if err != nil {
		os.Stderr.WriteString("Error: ")
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
