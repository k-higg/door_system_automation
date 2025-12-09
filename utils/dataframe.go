package utils

import (
	"Door_System_User_Automation/cmd"
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

	if path == cmd.WisFile {
		Normalize(df)
	}

	return df
}

// Normalize: used to normalize the headers so they match what is required to import to S2
func Normalize(df *dataframe.DataFrame) {
	err := df.Rename(map[string]string{
		"User ID":        "PERSONID",
		"Email":          "EMAIL_SRC",
		"Wireless phone": "MOBILEPHONE_SRC",
	})
	if err != nil {
		panic(err)
	}
	formatPhoneNumber(df)
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

// MergeAndExport: Rather than using the Merge() method, we will iterate over all of the values, checking to see if the PERSONID's match
// if the PERSONID's match, then we will update information using that index
func MergeAndExport(s2DF *dataframe.DataFrame, wisDF *dataframe.DataFrame) {
	const BlueDiamondStatus = "PENDING"
	const NFCBundle = "{Winchendon School NFC Bundle~20}"
	const NFCRequestStatus = "{Winchendon School NFC Bundle~Awaiting_REG}"
	const AddPerson = "AddPerson"
	const ModPerson = "ModifyPerson"
	const NoCommand = "NoCommand"

	result, _ := s2DF.Merge(wisDF, "PERSONID", dataframe.InnerMerge)

	_ = result.Rename(map[string]string{
		"MOBILEPHONE_SRC": "MOBILEPHONE",
		"EMAIL_SRC":       "EMAIL",
	})

	command, _ := result.SelectCol("COMMAND")
	bdEnabled, _ := result.SelectCol("BLUEDIAMONDENABLED")
	bdStatus, _ := result.SelectCol("BLUEDIAMONDSTATUS")
	mcRequest, _ := result.SelectCol("MOBILECREDENTIALREQUEST")
	mcRequestStatus, _ := result.SelectCol("MOBILECREDENTIALREQUESTSTATUS")

	for i := 0; i < bdEnabled.Len(); i++ {
		phoneNumber, _ := result.ILoc().At(i, 15)
		if phoneNumber == "" {
			_ = command.Set(i, NoCommand)
			_ = bdEnabled.Set(i, "FALSE")
		} else {
			_ = command.Set(i, ModPerson)
			_ = bdEnabled.Set(i, "TRUE")
			_ = bdStatus.Set(i, BlueDiamondStatus)
			_ = mcRequest.Set(i, NFCBundle)
			_ = mcRequestStatus.Set(i, NFCRequestStatus)
		}
	}

	_, err := result.ToCSV("output.csv")
	if err != nil {
		return
	}
}
