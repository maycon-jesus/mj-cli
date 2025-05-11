package utils

import "time"

func GetMonthName(month time.Month) string {
	monthName := ""

	switch month {
	case 1:
		monthName = "Janeiro"
	case 2:
		monthName = "Fevereiro"
	case 3:
		monthName = "Março"
	case 4:
		monthName = "Abril"
	case 5:
		monthName = "Maio"
	case 6:
		monthName = "Junho"
	case 7:
		monthName = "Julho"
	case 8:
		monthName = "Agosto"
	case 9:
		monthName = "Setembro"
	case 10:
		monthName = "Outubro"
	case 11:
		monthName = "Novembro"
	case 12:
		monthName = "Dezembro"
	}
	return monthName
}
