package main

import (
	// "crypto/tls"
	// "crypto/x509"
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/simonvetter/modbus"
)

func main() {
	var client *modbus.ModbusClient
	var err error
	// var clientKeyPair tls.Certificate
	// var serverCertPool *x509.CertPool
	var regs float32

	// load the client certificate and its associated private key, which
	// are used to authenticate the client to the server
	// clientKeyPair, err = tls.LoadX509KeyPair(
	// 	"certs/client.cert.pem", "certs/client.key.pem")
	// if err != nil {
	// 	fmt.Printf("failed to load client key pair: %v\n", err)
	// 	os.Exit(1)
	// }

	// load either the server certificate or the certificate of the CA
	// (Certificate Authority) which signed the server certificate
	// serverCertPool, err = modbus.LoadCertPool("certs/server.cert.pem")
	// if err != nil {
	// 	fmt.Printf("failed to load server certificate/CA: %v\n", err)
	// 	os.Exit(1)
	// }

	// create a client targetting host secure-plc on port 802 using
	// modbus TCP over TLS (MBAPS)
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		// tcp+tls is the moniker for MBAPS (modbus/tcp encapsulated in
		// TLS),
		// 802/tcp is the IANA-registered port for MBAPS.
		URL: "tcp://192.168.1.6:12345",
		// set the client-side cert and key
		// TLSClientCert: &clientKeyPair,
		// set the server/CA certificate
		// TLSRootCAs: serverCertPool,
	})
	if err != nil {
		fmt.Printf("failed to create modbus client: %v\n", err)
		os.Exit(1)
	}

	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		fmt.Printf("failed to connect: %v\n", err)
		os.Exit(2)
	}

	// read two 16-bit holding registers at address 0x4000
	file, _ := os.Open("./addrMap.csv")

	// csv reader 생성
	rdr := csv.NewReader(bufio.NewReader(file))

	// csv 내용 모두 읽기
	rows, _ := rdr.ReadAll()

	// 행,열 읽기
	for i := range rows {
		var a, err = strconv.ParseUint(rows[i][0], 10, 16)
		regs, err = client.ReadFloat32(uint16(a), modbus.INPUT_REGISTER)
		if err != nil {
			//fmt.Printf("failed to read registers 0x4000 and 0x4001: %v\n", err)
		} else {
			fmt.Printf("modbus address %d: %f\n", uint16(a), regs)
			// fmt.Printf("register 0x4001: 0x%04x\n", regs[1])
		}
	}

	// for i := 0; i < 9999; i++ {
	// 	regs, err = client.ReadFloat32(uint16(i), modbus.INPUT_REGISTER)
	// 	if err != nil {
	// 		//fmt.Printf("failed to read registers 0x4000 and 0x4001: %v\n", err)
	// 	} else {
	// 		fmt.Printf("register %d: %f\n", uint16(i), regs)
	// 		// fmt.Printf("register 0x4001: 0x%04x\n", regs[1])
	// 	}
	// }
	// regs, err = client.ReadRegister(30101, modbus.INPUT_REGISTER)
	// if err != nil {
	// 	fmt.Printf("failed to read registers 0x4000 and 0x4001: %v\n", err)
	// } else {
	// 	fmt.Printf("register 0x4000: 0x%04x\n", regs)
	// 	// fmt.Printf("register 0x4001: 0x%04x\n", regs[1])
	// }

	// set register 0x4002 to 500
	// err = client.WriteRegister(0x7595, 500)
	// if err != nil {
	// 	fmt.Printf("failed to write to register 0x4002: %v\n", err)
	// } else {
	// 	fmt.Printf("set register 0x4002 to 500\n")
	// }

	// close the connection
	err = client.Close()
	if err != nil {
		fmt.Printf("failed to close connection: %v\n", err)
	}

	os.Exit(0)
}
