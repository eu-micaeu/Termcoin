# Termcoin (cot)

Um monitor de cotações em tempo real para o terminal (CLI) escrito em **Go (GoLang)**. Permite consultar moedas tradicionais (Fiat) e principais criptomoedas de maneira extremamente rápida, leve e com interface colorida diretamente no seu terminal.

## Recursos
- **Arquitetura Modular Limpa**: Organizado em pacotes internos seguindo boas práticas de design Go (separação de responsabilidades em API, Lógica de Negócios/Moedas, e Interface do Usuário).
- **Zero Dependências Externas**: Utiliza apenas pacotes da biblioteca padrão do Go, gerando um binário único e autônomo.
- **Detecção Inteligente de TTY**: Identifica automaticamente se a saída padrão está sendo redirecionada para um arquivo ou pipe (ex: `cot usd > cotacao.txt`) e remove as cores ANSI para gerar um texto limpo.
- **Formatação Localizada (PT-BR)**: Formata os valores monetários no padrão brasileiro (ex: `R$ 323.367,00` e variação `-0,70%`).
- **Resolução de Sinônimos**: Suporta busca flexível (ex: `cot btc`, `cot bitcoin`, `cot dolar`, `cot USD`).

---

## APIs Utilizadas (Gratuitas e sem Token)
1. **Moedas Tradicionais**: [AwesomeAPI](https://docs.awesomeapi.com.br/) (Retorna dados em relação ao Real - BRL).
2. **Criptomoedas**: [CoinGecko API](https://www.coingecko.com/en/api) (Retorna cotações em Dólar - USD e Real - BRL).

---

## Moedas Suportadas

### Tradicionais (Fiat)
- `usd` (Dólar Americano)
- `eur` (Euro)
- `gbp` (Libra Esterlina)
- `ars` (Peso Argentino)
- `jpy` (Iene Japonês)
- `cad` (Dólar Canadense)
- `aud` (Dólar Australiano)
- `chf` (Franco Suíço)
- `cny` (Yuan Chinês)

### Criptomoedas
- `btc` (Bitcoin)
- `eth` (Ethereum)
- `sol` (Solana)
- `ada` (Cardano)
- `doge` (Dogecoin)
- `xrp` (Ripple)

---

## Estrutura de Arquivos no VS Code
Para abrir o projeto no VS Code, a estrutura de arquivos recomendada é a seguinte:

```text
Termcoin/
├── cmd/
│   └── cot/
│       └── main.go       # Ponto de entrada: resolve argumentos e invoca serviços
├── internal/
│   ├── api/
│   │   └── client.go     # Requisições HTTP e consumo JSON
│   ├── currency/
│   │   ├── formatter.go  # Formatadores numéricos e de precisão
│   │   ├── model.go      # Modelos de dados para desserialização
│   │   └── registry.go   # Lista de moedas suportadas e sinônimos
│   └── ui/
│       ├── colors.go     # Códigos ANSI e detecção de TTY
│       └── display.go    # Exibição de caixas, tabelas e ajuda
├── go.mod                # Definição do módulo Go (cot)
└── README.md             # Documentação do projeto (este arquivo)
```

---

## Como Compilar e Instalar no Ubuntu (Globais)

### Passo 1: Compilar Localmente
Se você deseja gerar apenas o binário na pasta do projeto para testes locais:
```bash
go build -o cot ./cmd/cot
./cot usd
```

### Passo 2: Instalar Globalmente com `go install`
O comando `go install` compila a ferramenta a partir da sua pasta principal e a instala no diretório padrão de binários do Go (`$HOME/go/bin` ou `$GOPATH/bin`).

De dentro do diretório `/home/micael/GitHub/Termcoin`, execute:
```bash
go install ./cmd/cot
```
Como o módulo no `go.mod` está nomeado como `cot`, o binário gerado terá exatamente o nome `cot`.

### Passo 3: Configurar o PATH no `~/.bashrc`
Para rodar o comando `cot` de qualquer diretório no seu terminal Ubuntu, o caminho `$HOME/go/bin` precisa estar contido na sua variável `$PATH`.

1. Adicione a exportação do caminho ao final do seu arquivo `~/.bashrc`:
   ```bash
   echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
   ```

2. Atualize o terminal atual para carregar a alteração imediatamente:
   ```bash
   source ~/.bashrc
   ```

3. Pronto! Agora você pode monitorar suas cotações digitando apenas `cot <moeda>` em qualquer terminal do seu sistema.

---

## Exemplos de Uso

### Consultando Dólar:
```bash
cot usd
```
```text
┌───────────────────────────────────────────────────────┐
│     Dólar Americano (USD) ➔ Real Brasileiro (BRL)     │
├───────────────────────────────────────────────────────┤
│  Cotação (Compra):     R$ 5,1407                      │
│  Cotação (Venda):      R$ 5,1437                      │
│  Máxima do Dia:        R$ 5,1832                      │
│  Mínima do Dia:        R$ 5,1325                      │
│  Variação do Dia:      -0,0362 (-0,70%)               │
│  Atualizado em:        2026-06-08 09:04:03            │
└───────────────────────────────────────────────────────┘
```

### Consultando Bitcoin:
```bash
cot btc
```
```text
┌───────────────────────────────────────────────────────┐
│         Bitcoin (BTC) ➔ Cotação em Tempo Real         │
├───────────────────────────────────────────────────────┤
│  Dólar (USD):          US$ 62.843,00                  │
│  Real (BRL):           R$ 323.367,00                  │
│  Fonte:                CoinGecko API                  │
└───────────────────────────────────────────────────────┘
```