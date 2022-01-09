package utils

import (
	"fmt"
	"strconv"
	"time"

	excel "github.com/I0HuKc/go-excel"
	"github.com/gefion-tech/tg-exchanger-server/internal/models"
)

func OneObmenDocumentGenerate(data *models.OneObmen, path string) string {
	tableStartFromLine := 1
	headValues := []string{"№", "From", "To", "In", "Out", "Amount", "MinAmount", "MaxAmount"}

	var tableHeader excel.Header = excel.Header{
		CellParams: excel.CreateHeaderCell(headValues, strconv.Itoa(tableStartFromLine)),
		ColParams: []excel.ColWidth{
			{
				StartCol: "A",
				EndCol:   "A",
				Width:    10, // column width
			},
			{
				StartCol: "B",
				EndCol:   "F",
				Width:    30,
			},
			{
				StartCol: "G",
				EndCol:   "H",
				Width:    25,
			},
		},
		RowParams: []excel.RowHeight{
			{
				Row:    tableStartFromLine,
				Height: 25,
			},
		},
	}

	var tableValue [][]interface{}
	for i, v := range data.Rates {
		tableValue = append(tableValue, []interface{}{
			i + 1,
			v.From,
			v.To,
			v.In,
			v.Out,
			v.Amount,
			v.MinAmount,
			v.MaxAmount,
		})
	}

	f := fmt.Sprintf("%s/file/%s.xlsx", path, time.Now().UTC().Format("2006-01-02T15:04:05.00000000"))
	excel.NewFile(f)

	excel.CreateDefaultTable(excel.DefaultTable{
		PathName:         f,
		TableHeader:      tableHeader,
		Data:             tableValue, // Data array written to the table
		Sheet:            "1obmen",   // Sheet name
		ContentRowHeight: 18,
		ContentLineStart: tableStartFromLine + 1,
	})

	return f
}
