package writer

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	internalmodels "bitbucket/api_format/internal/models"
)

// BuildCSV creates CSV output from parsed statement data.
func BuildCSV(info *internalmodels.StatementInfo, includeMetadata bool) (string, error) {
	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)

	if includeMetadata {
		metadataRows := [][2]string{
			{"Bank", string(info.Bank)},
			{"Account Holder", info.AccountHolder},
			{"Account Number", info.AccountNumber},
			{"Sort Code", info.SortCode},
			{"Statement Period", info.StatementPeriod},
		}

		for _, row := range metadataRows {
			if row[1] == "" {
				continue
			}
			if err := csvWriter.Write([]string{"# " + row[0], row[1]}); err != nil {
				return "", fmt.Errorf("write metadata csv row: %w", err)
			}
		}
	}

	if err := csvWriter.Write([]string{"Date", "Description", "Type", "Amount", "Balance"}); err != nil {
		return "", fmt.Errorf("write header csv row: %w", err)
	}

	for _, tx := range info.Transactions {
		if err := csvWriter.Write([]string{
			tx.Date,
			tx.Description,
			tx.Type,
			formatAmount(tx.Amount),
			formatAmount(tx.Balance),
		}); err != nil {
			return "", fmt.Errorf("write transaction csv row: %w", err)
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return "", fmt.Errorf("flush csv writer: %w", err)
	}

	return buffer.String(), nil
}

func formatAmount(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
