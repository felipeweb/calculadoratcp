package main

import (
	"net"
	"os"

	"github.com/felipeweb/tcp/calculator"
)

func main() {

	// Listens.
	server, err := net.ListenTCP("tcp", nil)
	calculator.EndIfErr(err)
	defer server.Close()

	// Logs server address.
	os.Stdout.WriteString("Endereço do servidor TCP: ")
	os.Stdout.WriteString(server.Addr().String())
	os.Stdout.WriteString("\n")

	// Accepts connection.
	conn, err := server.AcceptTCP()
	calculator.EndIfErr(err)

	for {
		// Receives data from the client.
		expr, err := calculator.ReadConn(conn)
		calculator.EndIfErr(err)

		// Logs received package.
		os.Stdout.WriteString("\nPacote recebido de ")
		os.Stdout.WriteString(conn.RemoteAddr().String())
		os.Stdout.WriteString(".\nExpressão recebida: ")
		os.Stdout.Write(expr)

		// Calculates expression.
		result, ok := calculator.CalculateExpression(calculator.CleanMsg(expr))
		if !ok {
			os.Stdout.WriteString("\nExpressão inválida.\n")
			_, err = conn.Write(calculator.Error)
			calculator.EndIfErr(err)
			continue
		}

		// Logs result.
		os.Stdout.WriteString(".\nResultado enviado: ")
		os.Stdout.Write(result)
		os.Stdout.WriteString(".\n")

		// Sends result to the client.
		_, err = conn.Write(calculator.MakeMsg(result))
		calculator.EndIfErr(err)
	}
}
