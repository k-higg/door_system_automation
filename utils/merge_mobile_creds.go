package utils

import (
	"fmt"
	"os"

	"github.com/apoplexi24/gpandas/dataframe"
)

// MergeMobile: Rather than using the Merge() method, we will iterate over all of the values, checking to see if the PERSONID's match
// if the PERSONID's match, then we will update information using that index
func MergeMobile(s2DF *dataframe.DataFrame, wisDF *dataframe.DataFrame) {
	const BlueDiamondStatus = "PENDING"
	const NFCBundle = "{Winchendon School NFC Bundle~20}"
	Command := [3]string{"NoCommand", "AddPerson", "ModifyPerson"}

	FormatWisDF(wisDF)

	result, _ := s2DF.Merge(wisDF, "PERSONID", dataframe.InnerMerge)

	_ = result.Rename(map[string]string{
		"MOBILEPHONE_SRC": "MOBILEPHONE",
		"EMAIL_SRC":       "EMAIL",
	})

	fmt.Println(result)

	command, _ := result.SelectCol("COMMAND")
	bdEnabled, _ := result.SelectCol("BLUEDIAMONDENABLED")
	bdStatus, _ := result.SelectCol("BLUEDIAMONDSTATUS")
	mcRequest, _ := result.SelectCol("MOBILECREDENTIALREQUEST")

	for i := 0; i < bdEnabled.Len(); i++ {
		phoneNumber, _ := result.ILoc().At(i, 15)
		if phoneNumber == "" {
			for j := 0; j < len(result.Columns); j++ {
				col, _ := result.ILoc().Col(j)
				err := col.Set(i, nil)
				if err != nil {
					panic(err)
				}
			}
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
