package bot_test

import (
	"fmt"
	"strings"

	"ronoaldo.gopkg.net/bot"
)

func Example() {
	b := bot.New()
	page, err := b.GET("http://www.valor.com.br/valor-data/moedas")
	if err != nil {
		fmt.Printf("Error carregando página: %v", err)
		return
	}
	tables, err := page.Tables()
	if err != nil {
		fmt.Printf("Erro extraindo tabelas: %v", err)
		return
	}

	// Lookup the value in the extracted tables
	if len(tables) > 0 {
		for _, table := range tables {
			// Find specific table to parse data from, by class attribute.
			if strings.Contains(table.Class, "valor_tabela") {
				// Lookup in the table data pre-processed row text values,
				// and print the interesting stuff.
				for _, row := range table.Data {
					if strings.TrimSpace(strings.ToLower(row[0])) == "dólar comercial" {
						fmt.Printf("Dólar Comercial, valor de compra: %s\n", strings.TrimSpace(row[1]))
					}
					if strings.Contains(strings.ToLower(row[0]), "fonte") {
						fmt.Println(row[0])
					}
				}
				break
			}
		}
	}
}
