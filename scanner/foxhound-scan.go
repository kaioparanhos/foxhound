package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
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

    salvarRelatorio(relatorio, config)

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

func worker(ip string, portas <-chan int, resultados chan<- Resultado, wg *sync.WaitGroup, timeout time.Duration, delay time.Duration) {
    defer wg.Done()

    for porta := range portas {
        resultados <- verificarPorta(ip, porta, timeout, delay)
    }
}

func verificarPorta(ip string, porta int, timeout time.Duration, delay time.Duration) Resultado {

    if delay > 0 {
        time.Sleep(delay * time.Millisecond) // Atraso entre conexões
    }

    target := fmt.Sprintf("%s:%d", ip, porta)
    conn, err := net.DialTimeout("tcp", target, timeout)
    if err != nil {
        return Resultado{IP: ip, Porta: porta, Status: "FECHADA", Mensagem: err.Error()}
    }
    defer conn.Close()

    servico := servicos[porta]

    if servico == "" {
        servico = "Desconhecido"
    }

    return Resultado{IP: ip, Porta: porta, Status: "ABERTA", Servico: servico}
}

func consultarWHOIS(ip string) string {
    result, err := whois.Whois(ip)
    if err != nil {
        return fmt.Sprintf("Erro ao consultar WHOIS: %v", err)
    }
    return fmt.Sprintf("WHOIS Info:\n%s", result)
}

func salvarRelatorio(relatorio []Resultado, config Config) {
    switch config.Format {
    case "json":
        salvarRelatorioJSON(relatorio, config)
    case "csv":
        salvarRelatorioCSV(relatorio, config)
    case "txt":
        salvarRelatorioTXT(relatorio, config)
    case "xlsx":
        salvarRelatorioODS(relatorio, config)
    default:
        log.Fatalf("Formato de relatório inválido: %s", config.Format)
    }
}


func salvarRelatorioJSON(relatorio []Resultado, config Config) {
    timestamp := time.Now().Format("02-01-2006-15-04-05")
    arquivoNome := fmt.Sprintf("%s/%s-%s.json", config.Dir, config.Output, timestamp)

    arquivo, _ := json.MarshalIndent(relatorio, "", "  ")
    os.WriteFile(arquivoNome, arquivo, 0644)

    log.Printf("Relatório salvo em %s!", arquivoNome)
}


func salvarRelatorioCSV(relatorio []Resultado, config Config) {
    timestamp := time.Now().Format("02-01-2006-15-04-05")
    arquivoNome := fmt.Sprintf("%s/%s-%s.csv", config.Dir, config.Output, timestamp)

    arquivo, _ := os.Create(arquivoNome)
    defer arquivo.Close()
    writer := csv.NewWriter(arquivo)
    writer.Write([]string{"IP", "Porta", "Status", "Servico", "Mensagem"})
    for _, r := range relatorio {
        writer.Write([]string{r.IP, fmt.Sprint(r.Porta), r.Status, r.Servico, r.Mensagem})
    }
    writer.Flush()

    log.Printf("Relatório salvo em %s!", arquivoNome)
}


func salvarRelatorioTXT(relatorio []Resultado, config Config) {
    timestamp := time.Now().Format("02-01-2006-15-04-05")
    arquivoNome := fmt.Sprintf("%s/%s-%s.txt", config.Dir, config.Output, timestamp)

    arquivo, _ := os.Create(arquivoNome)
    defer arquivo.Close()
    for _, r := range relatorio {
        arquivo.WriteString(fmt.Sprintf("IP: %s, Porta: %d, Status: %s, Servico: %s, Mensagem: %s\n", r.IP, r.Porta, r.Status, r.Servico, r.Mensagem))
    }

    log.Printf("Relatório salvo em %s!", arquivoNome)
}


func salvarRelatorioODS(relatorio []Resultado, config Config) {
    timestamp := time.Now().Format("02-01-2006-15-04-05")
    arquivoNome := fmt.Sprintf("%s/%s-%s.xlsx", config.Dir, config.Output, timestamp)

    file := xlsx.NewFile()
    sheet, err := file.AddSheet("Relatório")
    if err != nil {
        log.Fatalf("Erro ao criar planilha: %v", err)
    }

    
    row := sheet.AddRow()
    row.AddCell().SetString("IP")
    row.AddCell().SetString("Porta")
    row.AddCell().SetString("Status")
    row.AddCell().SetString("Serviço")
    row.AddCell().SetString("Mensagem")

    
    for _, r := range relatorio {
        row := sheet.AddRow()
        row.AddCell().SetString(r.IP)
        row.AddCell().SetInt(r.Porta)
        row.AddCell().SetString(r.Status)
        row.AddCell().SetString(r.Servico)
        row.AddCell().SetString(r.Mensagem)
    }

    if err := file.Save(arquivoNome); err != nil {
        log.Fatalf("Erro ao salvar arquivo .xlsx: %v", err)
    }

    log.Printf("Relatório salvo em %s!", arquivoNome)
}


func parseIPs(ips string) []string {
    return strings.Split(ips, ",")
}