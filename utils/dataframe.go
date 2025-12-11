package utils

import (
	"fmt"

	"github.com/apoplexi24/gpandas"
	"github.com/apoplexi24/gpandas/dataframe"
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

// NormalizeUID: used to normalize the headers so they match what is required to import to S2
func NormalizeUID(df *dataframe.DataFrame) {
	err := df.Rename(map[string]string{
		"User ID": "PERSONID",
	})

	if err != nil {
		fmt.Println("Error: tried to Rename column headers that don't exist in the provided file..")
	}
}

// FormatWisDF: used to format WisDF when using the mobile flag
// sets entire row to null if no phonenumber is present, this will cause the row to drop when merging
func FormatWisDF(df *dataframe.DataFrame) {
	_ = df.Rename(map[string]string{
		"Email":          "EMAIL",
		"Wireless phone": "MOBILEPHONE",
	})

	phoneNumberCol, err := df.SelectCol("MOBILEPHONE")
	if err != nil {
		fmt.Println("Error: could not find MOBILEPHONE column..")
	}
	for i := 0; i < phoneNumberCol.Len(); i++ {
		phoneNumber, _ := phoneNumberCol.At(i)
		if phoneNumber == "" {
			for j := 0; j < len(df.Columns); j++ {
				col, _ := df.ILoc().Col(j)
				err := col.Set(i, nil)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	FormatPhoneNumber(df)
}
