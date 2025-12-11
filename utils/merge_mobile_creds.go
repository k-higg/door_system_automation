package utils

import (
	"fmt"

	"github.com/apoplexi24/gpandas/dataframe"
)

// MergeMobile: Rather than using the Merge() method, we will iterate over all of the values, checking to see if the PERSONID's match
// if the PERSONID's match, then we will update information using that index
func MergeMobile(s2DF *dataframe.DataFrame, wisDF *dataframe.DataFrame) {
	Command := [3]string{"NoCommand", "AddPerson", "ModifyPerson"}

	FormatWisDF(wisDF)

	result, _ := s2DF.Merge(wisDF, "PERSONID", dataframe.InnerMerge)
	command, _ := result.SelectCol("COMMAND")

	for i := 0; i < command.Len(); i++ {
		phoneNumber, _ := result.ILoc().At(i, 15)
		if phoneNumber == "" {
			fmt.Println("Didn't format correctly......")
		} else {
			bdEnabled, _ := result.SelectCol("BLUEDIAMONDENABLED")
			_ = command.Set(i, Command[2])
			_ = bdEnabled.Set(i, "TRUE")
		}
	}

	ExportDataFrame(result, "phonenumber_email_update.csv")
	MergeMobileCreds(result)
}

func MergeMobileCreds(df *dataframe.DataFrame) {
	const NFCBundle = "{Winchendon School NFC Bundle~20}"

	command, _ := df.SelectCol("COMMAND")
	mcRequest, _ := df.SelectCol("MOBILECREDENTIALREQUEST")

	for i := 0; i < command.Len(); i++ {
		_ = mcRequest.Set(i, NFCBundle)
	}

	ExportDataFrame(df, "mobilecredentials.csv")
}
