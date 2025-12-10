package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/apoplexi24/gpandas"
	"github.com/apoplexi24/gpandas/dataframe"
	"github.com/nyaruka/phonenumbers"
)

// CreateDataFrame: creates a dataframe from the provided csv
func CreateDataFrame(path string) *dataframe.DataFrame {
	cleanPath, err := StripBOM(path)
	if err != nil {
		panic(err)
	}

	pd := gpandas.GoPandas{}

	df, err := pd.Read_csv(cleanPath)
	if err != nil {
		panic(err)
	}

	return df
}

// Normalize: used to normalize the headers so they match what is required to import to S2
func Normalize(df *dataframe.DataFrame) {
	err := df.Rename(map[string]string{
		"User ID":                   "PERSONID",
		"Email":                     "EMAIL_SRC",
		"Wireless phone":            "MOBILEPHONE_SRC",
		"Last name":                 "LASTNAME_SRC",
		"Preferred else First name": "FIRSTNAME_SRC",
	})

	if err != nil {
		fmt.Println("Error: tried to Rename column headers that don't exist in the provided file..")
	}

	_, err = df.SelectCol("MOBILEPHONE_SRC")
	if err != nil {
		fmt.Println("Error: could not find MOBILEPHONE_SRC column..")
	} else {
		formatPhoneNumber(df)
	}
}

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

func formatPhoneNumber(df *dataframe.DataFrame) {
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

// TODO: Drop Rows with NoCommand for their COMMAND
// MergeAndExport: Rather than using the Merge() method, we will iterate over all of the values, checking to see if the PERSONID's match
// if the PERSONID's match, then we will update information using that index
func MergeAndExport(s2DF *dataframe.DataFrame, wisDF *dataframe.DataFrame) {
	const BlueDiamondStatus = "PENDING"
	const NFCBundle = "{Winchendon School NFC Bundle~20}"
	const AddPerson = "AddPerson"
	const ModPerson = "ModifyPerson"
	const NoCommand = "NoCommand"

	result, _ := s2DF.Merge(wisDF, "PERSONID", dataframe.InnerMerge)

	_ = result.Rename(map[string]string{
		"MOBILEPHONE_SRC": "MOBILEPHONE",
		"EMAIL_SRC":       "EMAIL",
	})

	command, _ := result.SelectCol("COMMAND")
	bdEnabled, err := result.SelectCol("BLUEDIAMONDENABLED")
	bdStatus, _ := result.SelectCol("BLUEDIAMONDSTATUS")
	mcRequest, _ := result.SelectCol("MOBILECREDENTIALREQUEST")

	if err != nil {
		return
	}

	for i := 0; i < bdEnabled.Len(); i++ {
		phoneNumber, _ := result.ILoc().At(i, 15)
		if phoneNumber == "" {
			_ = command.Set(i, NoCommand)
			_ = bdEnabled.Set(i, "FALSE")
		} else {
			_ = command.Set(i, ModPerson)
			_ = bdEnabled.Set(i, "TRUE")
			_ = bdStatus.Set(i, BlueDiamondStatus)

		}
	}

	os.Mkdir("output", 0755)
	_, err = result.ToCSV("output/import_first.csv")
	if err != nil {
		return
	}

	for i := 0; i < command.Len(); i++ {
		_ = mcRequest.Set(i, NFCBundle)
	}

	_, err = result.ToCSV("output/import_second.csv")
	if err != nil {
		return
	}
}
