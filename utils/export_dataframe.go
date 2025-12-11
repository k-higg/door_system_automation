package utils

import (
	"os"

	"github.com/apoplexi24/gpandas/dataframe"
)

func ExportDataFrame(df *dataframe.DataFrame, path string) {
	_ = os.Mkdir("output", 0755)
	_, err := df.ToCSV("output/" + path)
	if err != nil {
		return
	}
}
