package main

import (

    "encoding/csv"
    "encoding/json"
	"bufio"
    "fmt"
    "log"
    "net"
    "os"
	"strconv"
    "strings"
    "sync"
    "time"

	// lib
	"github.com/likexian/whois"

)

type Config struct {

	IPs        []string `json:"ips"`
	Inicio     int      `json:"inicio"`
	Fim        int      `json:"fim"`
	Timeout    int      `json:"timeout"`
	Workers    int      `json:"workers"`
	Output     string   `json:"output"`
	Format     string   `json:"format"`
	Rapido     bool     `json:"rapido"`
	Dir        string   `json:"dir"`
	
}

type Result struct {

	IP        string    `json:"ip"`
	Porta     int    	`json:"porta"`
	Status    string 	`json:"status"`
	Servico   string 	`json:"servico,omitempty"`
	Mensagem  string 	`json:"mensagem,omitempty"`
	Porta     int       `json:"porta"`
	Status    string    `json:"status"`
	Servico   string    `json:"servico,omitempty"`
	Mensagem  string    `json:"mensagem,omitempty"`

}

var services = map [int] string {

    22:  	"SSH",
    80:  	"HTTP",
    443: 	"HTTPS",
    21:  	"FTP",
    25:  	"SMTP",

}

var cports = [] int {22, 80, 443, 21, 25, 53, 110, 143, 3306, 3389}

func main() {

    config := UserInput()

    if len(config.IPs) == 0 {

        log.Fatal("Nenhum IP especificado!")

    }

    if err := os.MkdirAll(config.Dir, 0755); err != nil {

        log.Fatalf("Erro ao criar diretório: %v", err)

    }

    var report []Result
    for _, ip := range config.IPs {

        finds := IPscanner(ip, config)
        report = append(report, finds...)

    }

    saveReport(report, config)

}
 