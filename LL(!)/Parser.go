package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Token struct {
	valor string
}

func tokenizar(entrada string) []Token {
	tokens := strings.Fields(entrada)
	var resultado []Token
	for _, token := range tokens {
		resultado = append(resultado, Token{valor: token})
	}
	return resultado
}

type LL1Entrada struct {
	naoTerminal string
	terminal    string
	producao    string
}

type LL1Tabela struct {
	entradas []LL1Entrada
}

func (tabela *LL1Tabela) adicionarEntrada(naoTerminal, terminal, producao string) {
	tabela.entradas = append(tabela.entradas, LL1Entrada{naoTerminal: naoTerminal, terminal: terminal, producao: producao})
}

func carregarTabela(arquivoCSV string) (*LL1Tabela, error) {
	tabela := &LL1Tabela{}

	arquivo, err := os.Open(arquivoCSV)
	if err != nil {
		return nil, fmt.Errorf("Erro: Falha ao abrir o arquivo CSV '%s'", arquivoCSV)
	}
	defer arquivo.Close()

	leitorCSV := csv.NewReader(arquivo)
	for {
		linha, err := leitorCSV.Read()
		if err != nil {
			break
		}
		if len(linha) == 3 {
			tabela.adicionarEntrada(strings.TrimSpace(linha[0]), strings.TrimSpace(linha[1]), strings.TrimSpace(linha[2]))
		}
	}

	return tabela, nil
}

func analisarEntrada(entrada string, tabela *LL1Tabela) bool {
	tokens := tokenizar(entrada)
	tokens = append(tokens, Token{valor: "$"})

	pilha := []string{"$"}
	pilha = append(pilha, tabela.entradas[0].naoTerminal)

	indiceEntrada := 0

	for len(pilha) > 0 {
		topoPilha := pilha[len(pilha)-1]
		pilha = pilha[:len(pilha)-1]

		if topoPilha == "$" {
			if indiceEntrada == len(tokens)-1 {
				return true
			} else {
				fmt.Println("Erro: Não foi possível consumir toda a entrada.")
				for indiceEntrada < len(tokens)-1 {
					fmt.Printf("Tokens remanescentes: %s ", tokens[indiceEntrada].valor)
					indiceEntrada++
				}
				return false
			}
		}

		if indiceEntrada < len(tokens) && topoPilha == tokens[indiceEntrada].valor {
			indiceEntrada++
		} else {
			var entradaProducao *LL1Entrada
			for _, entrada := range tabela.entradas {
				if entrada.naoTerminal == topoPilha && entrada.terminal == tokens[indiceEntrada].valor {
					entradaProducao = &entrada
					break
				}
			}

			if entradaProducao == nil {
				fmt.Printf("Erro: Nenhuma produção encontrada para '%s' com '%s'\n", topoPilha, tokens[indiceEntrada].valor)
				if indiceEntrada > 0 {
					fmt.Printf("O token '%s' foi inesperado após '%s'\n", tokens[indiceEntrada].valor, tokens[indiceEntrada-1].valor)
				}
				return false
			}

			simbolos := strings.Fields(entradaProducao.producao)
			for i := len(simbolos) - 1; i >= 0; i-- {
				if simbolos[i] != "ε" {
					pilha = append(pilha, simbolos[i])
				}
			}
		}
	}

	fmt.Println("Erro: A entrada foi rejeitada pela gramática.")
	return false
}

func executarLL1Interpretador(arquivoCSV, arquivoEntrada string) bool {
	tabela, err := carregarTabela(arquivoCSV)
	if err != nil {
		fmt.Println(err)
		return false
	}

	entrada, err := os.ReadFile(arquivoEntrada)
	if err != nil {
		fmt.Printf("Erro: Não foi possível abrir o arquivo de entrada '%s'.\n", arquivoEntrada)
		return false
	}

	return analisarEntrada(string(entrada), tabela)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Uso: go run main.go <arquivo_tabela_csv> <arquivo_entrada>")
		return
	}

	arquivoCSV := os.Args[1]
	arquivoEntrada := os.Args[2]

	sucesso := executarLL1Interpretador(arquivoCSV, arquivoEntrada)
	if sucesso {
		fmt.Println("Análise concluída com sucesso: Entrada aceita.")
	} else {
		fmt.Println("Análise falhou: Entrada não aceita.")
	}
}
