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
		"Last name":                 "LASTNAME_SRC",
		"Preferred else First name": "FIRSTNAME_SRC",
	})

	if err != nil {
		fmt.Println("Error: tried to Rename column headers that don't exist in the provided file..")
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

// TODO: refactor
func fixWisDF(df *dataframe.DataFrame) {
	_ = df.Rename(map[string]string{
		"Email":          "EMAIL_SRC",
		"Wireless phone": "MOBILEPHONE_SRC",
	})

	_, err := df.SelectCol("MOBILEPHONE_SRC")
	if err != nil {
		fmt.Println("Error: could not find MOBILEPHONE_SRC column..")
	} else {
		formatPhoneNumber(df)
	}
}

// TODO: Drop Rows with NoCommand for their COMMAND
// MergeMobile: Rather than using the Merge() method, we will iterate over all of the values, checking to see if the PERSONID's match
// if the PERSONID's match, then we will update information using that index
func MergeMobile(s2DF *dataframe.DataFrame, wisDF *dataframe.DataFrame) {
	const BlueDiamondStatus = "PENDING"
	const NFCBundle = "{Winchendon School NFC Bundle~20}"
	Command := [3]string{"NoCommand", "AddPerson", "ModifyPerson"}

	fixWisDF(wisDF)

	result, _ := s2DF.Merge(wisDF, "PERSONID", dataframe.InnerMerge)

	_ = result.Rename(map[string]string{

		"MOBILEPHONE_SRC": "MOBILEPHONE",
		"EMAIL_SRC":       "EMAIL",
	})

	command, _ := result.SelectCol("COMMAND")
	bdEnabled, _ := result.SelectCol("BLUEDIAMONDENABLED")
	bdStatus, _ := result.SelectCol("BLUEDIAMONDSTATUS")
	mcRequest, _ := result.SelectCol("MOBILECREDENTIALREQUEST")

	for i := 0; i < bdEnabled.Len(); i++ {
		phoneNumber, _ := result.ILoc().At(i, 15)
		if phoneNumber == "" {
			_ = command.Set(i, Command[0])
			_ = bdEnabled.Set(i, "FALSE")
		} else {
			_ = command.Set(i, Command[2])
			_ = bdEnabled.Set(i, "TRUE")
			_ = bdStatus.Set(i, BlueDiamondStatus)
		}
	}

	_ = os.Mkdir("output", 0755)
	_, err := result.ToCSV("output/import_first.csv")
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

// MergeGender: Rather than using the Merge() method, we will iterate over all of the values, checking to see if the PERSONID's match
// if the PERSONID's match, then we will update information using that index
func MergeGender(s2DF *dataframe.DataFrame, wisDF *dataframe.DataFrame) {
	AccessArr := [3]string{
		"General - All Students~~~FALSE~FALSE",
		"General - All Male Students~~~FALSE~FALSE",
		"General - All Female Students~~~FALSE~FALSE",
	}
	CommandArr := [3]string{"NoCommand", "AddPerson", "ModifyPerson"}

	result, _ := s2DF.Merge(wisDF, "PERSONID", dataframe.InnerMerge)

	_ = result.Rename(map[string]string{
		"LASTNAME_SRC":  "LASTNAME",
		"FIRSTNAME_SRC": "FIRSTNAME",
		"Gender":        "UDF4",
	})

	command, _ := result.SelectCol("COMMAND")
	accesslevel, _ := result.SelectCol("ACCESSLEVELS")
	gender, _ := result.SelectCol("UDF4")

	for i := 0; i < command.Len(); i++ {
		// TODO: edit current access levels based on gender
		_ = command.Set(i, CommandArr[2])
		g, err := gender.At(i)
		if err != nil {
			continue
		}
		if g == "Male" {
			originalString, _ := accesslevel.At(i)
			newString := strings.ReplaceAll(originalString.(string), AccessArr[0], AccessArr[1])
			_ = accesslevel.Set(i, newString)
		}
		if g == "Female" {
			originalString, _ := accesslevel.At(i)
			newString := strings.ReplaceAll(originalString.(string), AccessArr[0], AccessArr[2])
			_ = accesslevel.Set(i, newString)
		}
	}

	_ = os.Mkdir("output", 0755)
	_, err := result.ToCSV("output/updated_access.csv")
	if err != nil {
		panic(err)
	}
}
