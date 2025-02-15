package main

// lib
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"
    "strings"

	// libs p/ consulta
	"github.com/likexian/whois"
	"github.com/ns3777k/go-shodan/shodan"
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
	Shodan     string   `json:"shodan"`
	VirusTotal string   `json:"virustotal"`
}

type Result struct {
	IP       string `json:"ip"`
	Porta    int    `json:"porta"`
	Status   string `json:"status"`
	Servico  string `json:"servico,omitempty"`
	Mensagem string `json:"mensagem,omitempty"`
}
