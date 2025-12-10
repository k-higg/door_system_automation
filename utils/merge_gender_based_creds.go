package utils

import (
	"os"
	"strings"

	"github.com/apoplexi24/gpandas/dataframe"
)

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
