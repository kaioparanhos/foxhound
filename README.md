
# **FOXHOUND - PORT SCANNER**

Este é um projeto de escaneamento de portas desenvolvido em Go, que permite verificar portas abertas em um ou mais IPs, consultar informações WHOIS e salvar os resultados em diferentes formatos (JSON, CSV, XLSX, TXT).

---

## **Funcionalidades**

- **Escaneamento de Portas:**
  - Verifica portas abertas em um intervalo especificado.
  - Modo rápido para escanear apenas portas comuns.
- **Consulta WHOIS:**
  - Obtém informações WHOIS para IPs com portas abertas.
- **Relatórios Personalizáveis:**
  - Salva os resultados em JSON, CSV,  XLSX ou TXT.
  - Permite escolher o diretório de saída e o nome do arquivo.
- **Interação com o Usuário:**
  - Solicita configurações diretamente no terminal.

---

## **Requisitos**

- **Go 1.20 ou superior:** [Instalação do Go](https://golang.org/doc/install)
- **Bibliotecas Externas:**
  - `github.com/likexian/whois`
  - Instale as dependências com:
    ```bash
    go get github.com/likexian/whois
    ```
    
  - `github.com/tealeg/xlsx`
  - Instale as dependências com:
    ```bash
    go get github.com/tealeg/xlsx
    ```

---

## **Como Usar**

### **1. Clonar o Repositório**

```bash
git clone https://github.com/kaioparanhos/foxhound
cd foxhound/scanner
```

### **2. Executar o Programa**

Execute o programa com o seguinte comando:

```bash
go run foxhound-scan.go
```

O programa solicitará as configurações diretamente no terminal. Siga as instruções abaixo:

---

### **Configurações do Usuário**

1. **IPs ou Sub-redes:**
   - Digite os IPs ou sub-redes que deseja escanear, separados por vírgula.
   - Exemplo: `192.168.1.1,192.168.1.0/24`

2. **Porta Inicial e Final:**
   - Defina o intervalo de portas a serem escaneadas.
   - Padrão: `1` (início) e `1024` (fim).

3. **Timeout:**
   - Tempo máximo de espera para cada conexão (em segundos).
   - Padrão: `2`.

4. **Número de Workers:**
   - Quantidade de goroutines para escanear portas simultaneamente.
   - Padrão: `100`.

5. **Prefixo do Arquivo de Saída:**
   - Nome base do arquivo de relatório.
   - Padrão: `resultado`.

6. **Formato do Relatório:**
   - Escolha entre `json`, `csv`, `xlsx` ou `txt`.
   - Padrão: `json`.

7. **Modo Rápido:**
   - Escaneia apenas portas comuns (22, 80, 443, etc.).
   - Responda `s` para ativar ou `n` para desativar.
   - Padrão: `n`.

8. **Diretório de Saída:**
   - Pasta onde o relatório será salvo.
   - Padrão: `.` (diretório atual).

---

### **Exemplo de Uso**

```plaintext
Digite os IPs ou sub-redes (separados por vírgula): 192.168.1.1,192.168.1.2
Porta inicial (padrão: 1): 1
Porta final (padrão: 1024): 1024
Timeout em segundos (padrão: 2): 2
Número de workers (padrão: 100): 50
Prefixo do arquivo de saída (padrão: resultado): scan
Formato do relatório (json, csv, xlsx, txt) (padrão: json): csv
Modo rápido (escaneia apenas portas comuns)? (s/n) (padrão: n): n
Diretório para salvar o relatório (padrão: .): relatorios
Atraso entre cada conexão (em milissegundos) (padrão: 0): 0
```

---

### **Estrutura do Relatório**

O relatório contém as seguintes informações para cada porta escaneada:

- **IP:** Endereço IP escaneado.
- **Porta:** Número da porta.
- **Status:** `ABERTA` ou `FECHADA`.
- **Serviço:** Serviço associado à porta (ex: HTTP, SSH).
- **Mensagem:** Informações adicionais (ex: resultado do WHOIS).

---

### **Diretório de Saída**

Os relatórios são salvos no diretório especificado, com o nome no formato:

```
<nome-do-arquivo>-<timestamp>.<formato>
```

Exemplo: `scan-02-01-2023-15-04-05.csv`
