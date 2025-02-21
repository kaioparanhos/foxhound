package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"go/scanner"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	// lib
	"github.com/likexian/whois"
    "github.com/tealeg/xlsx"
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
    Delay      int      `json:"delay"`
	
}

type Resultado struct {

	IP        string    `json:"ip"`
	Porta     int       `json:"porta"`
	Status    string    `json:"status"`
	Servico   string    `json:"servico,omitempty"`
	Mensagem  string    `json:"mensagem,omitempty"`

}

var servicos = map [int] string {

    22:  	"SSH",
    80:  	"HTTP",
    443: 	"HTTPS",
    21:  	"FTP",
    25:  	"SMTP",

}

var portas_comuns = [] int {22, 80, 443, 21, 25, 53, 110, 143, 3306, 3389}

func main() {

    config := UserInput()

    if len(config.IPs) == 0 {

        log.Fatal("Nenhum IP especificado!")

    }

    if err := os.MkdirAll(config.Dir, 0755); err != nil {

        log.Fatalf("Erro ao criar diretório: %v", err)

    }

    var relatorio []Resultado
    for _, ip := range config.IPs {

        resultados := scanIP(ip, config)
        relatorio = append(relatorio, resultados...)

    }

    saveReport(relatorio, config)

}

func UserInput() Config {

    scanner := bufio.NewScanner(os.Stdin)

    fmt.Print("IPs ou sub-redes: ")
    scanner.Scan()
    ips := scanner.Text()

    fmt.Print("Porta inicial: ")
    scanner.Scan()
    inicio, _ := strconv.Atoi(scanner.Text())
    
    if inicio == 0 {
        inicio = 1
    }

    fmt.Print("Porta final: ")
    scanner.Scan()
    fim, _ := strconv.Atoi(scanner.Text())

    if fim == 0 {
        fim = 1024
    }

    fmt.Print("Timeout: ")
    scanner.Scan()
    timeout, _ := strconv.Atoi(scanner.Text())

    if timeout == 0 {
        timeout = 2
    }

    fmt.Print("Workers: ")
    scanner.Scan()
    workers, _ := strconv.Atoi(scanner.Text())

    if workers == 0 {
        workers = 100
    }

    fmt.Print("Prefixo do arquivo de saída: ")
    scanner.Scan()
    output := scanner.Text()

    if output == "" {
        output = "resultado"
    }

    fmt.Print("Formato do relatório: ")
    scanner.Scan()
    format := scanner.Text()

    if format == "" {
        format = "json"
    }

    fmt.Print("Modo rápido? (s/n): ")
    scanner.Scan()
    rapido := strings.ToLower(scanner.Text()) == "s"

    fmt.Print("Diretório para salvar o relatório: ")
    scanner.Scan()
    dir := scanner.Text()

    if dir == "" {
        dir = "."
    }

    fmt.Print("Atraso entre cada conexão: ")
    scanner.Scan()
    delay, _ := strconv.Atoi(scanner.Text())


    return Config{

        IPs:        parseIPs(ips),
        Inicio:     inicio,
        Fim:        fim,
        Timeout:    timeout,
        Workers:    workers,
        Output:     output,
        Format:     format,
        Rapido:     rapido,
        Dir:        dir,
        Delay:      delay,

    }
    
}

func scanIP(ip string, config Config) []Resultado {

    portas := make(chan int, config.Workers)
    resultados := make(chan Resultado, config.Fim-config.Inicio+1)
    var wg sync.WaitGroup

    for i := 0; i < config.Workers; i++ {
        wg.Add(1)
        go worker(ip, portas, resultados, &wg, time.Duration(config.Timeout)*time.Second, time.Duration(config.Delay))
    }

    go func() {
        if config.Rapido {
            for _, porta := range portas_comuns {
                portas <- porta
            }
        } else {
            for porta := config.Inicio; porta <= config.Fim; porta++ {
                portas <- porta
            }
        }

        close(portas)
    }()

    var relatorio []Resultado
    go func() {
        wg.Wait()
        close(resultados)
    }()

    for resultado := range resultados {
        relatorio = append(relatorio, resultado)
    }

    for i := range relatorio {

        if relatorio[i].Status == "ABERTA" {
            relatorio[i].Mensagem = consultarWHOIS(ip)
        }

    }

    return relatorio

}

