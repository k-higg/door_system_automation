package utils

import (
	"regexp"
	"strings"

	"github.com/apoplexi24/gpandas/dataframe"
	"github.com/nyaruka/phonenumbers"
)

func sanitizePhoneNumber(number string) string {
	if number == "" {
		return ""
	}

	number = strings.TrimPrefix(number, "+")
	number = strings.TrimPrefix(number, "1")

	re := regexp.MustCompile("[ ()-]")
	number = re.ReplaceAllString(number, "")

	parsedNumber, err := phonenumbers.Parse(number, "US")
	if err == nil {
		formattedNumber := phonenumbers.Format(
			parsedNumber,
			phonenumbers.INTERNATIONAL,
		)
		return strings.ReplaceAll(formattedNumber, "-", " ")
	}

	return number
}

func FormatPhoneNumber(df *dataframe.DataFrame) {
	series, err := df.SelectCol("MOBILEPHONE_SRC")
	if err != nil {
		panic(err)
	}

	for i := 0; i < series.Len(); i++ {
		val, err := series.At(i)
		if err != nil {
			continue
		}

		src, ok := val.(string)
		if !ok || src == "" {
			continue
		}

		formatted := sanitizePhoneNumber(src)

		err = series.Set(i, formatted)
		if err != nil {
			panic(err)
		}
	}
}
