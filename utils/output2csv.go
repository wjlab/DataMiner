package utils

import (
	"dataMiner/models"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"strconv"
	"time"
)

/*
  Overview func: output to csv file
  @Param  csv (the result of overview function)
  @Param  outputID (the output file name)
*/
func SavetocsvO(csv []models.OverviewData,outputID InfoStruct){
	categories := map[string]string{
		"A1": "Database", "B1": "Table/Collection","C1":"RowsCount/DocumentsCount"}
	var values = make(map[string]string)
	sumTmp:=0
	for i:=0;i<len(csv);i++{
		sumTmp++
		//set excel values
		excelValuetmpA:="A"+strconv.Itoa(sumTmp+1)
		excelValuetmpB:="B"+strconv.Itoa(sumTmp+1)
		excelValuetmpC:="C"+strconv.Itoa(sumTmp+1)
		values[excelValuetmpA]= csv[i].DatabaseName
		values[excelValuetmpB]= csv[i].TableName
		values[excelValuetmpC]= csv[i].RowCount
	}
	//output to a excel
	f := excelize.NewFile()
	err := f.SetColWidth("Sheet1", "A", "C", 50)
	if err!=nil{
		log.Fatalf("SetColWidth goes wrong!!!")
	}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
		},
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = f.SetCellStyle("Sheet1", "A1", "A1", style)
	err = f.SetCellStyle("Sheet1", "B1", "B1", style)
	err = f.SetCellStyle("Sheet1", "C1", "C1", style)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// save the result to xlsx
	currentTime := time.Now().Format("20060102150405")
	filename:=outputID.IP+"_"+outputID.Port+"_"+outputID.User+"_Overview_"+currentTime+".xlsx"
	if err := f.SaveAs(filename); err != nil {
		log.Fatalf(err.Error())
	}
	f.Close()
	log.Println("The results(csv) saved in "+filename)
}

/*
  SampleData func: output to csv file
  @Param  csv (the result of SampleData function)
  @Param  outputID (the output file name)
  @Param  num the number of rows returned from database
*/
func Savetocsv(csv []models.SampleStruct,outputID InfoStruct,num int) {
	// Map to group items by database name
	groupedItems := make(map[string][]models.SampleStruct)

	// Iterate over the table info items and group them by database name
	for _, item := range csv {
		if _, ok := groupedItems[item.DatabaseName]; !ok {
			groupedItems[item.DatabaseName] = make([]models.SampleStruct, 0)
		}
		groupedItems[item.DatabaseName] = append(groupedItems[item.DatabaseName], item)
	}

	// Create a new Excel file
	f := excelize.NewFile()

	// Map to store sheets and their corresponding rows
	sheets := make(map[string]int)

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   13,
		},
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	categories := map[string]string{
		 "A1": "Table/Collection","B1":"Data"}

	// align center
	centerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal:      "center",
			Vertical:        "center",
		},
		Font: &excelize.Font{
			Bold:   true,
			Size:   12,
		},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// bold the columns names
	fontStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
		},
	})

	for dbName, items := range groupedItems {
		// Get the sheet for the database
		sheetName := dbName
		sheetIndex, ok := sheets[sheetName]
		if !ok {
			// Create a new sheet for the database
			sheetIndex = f.NewSheet(sheetName)
			sheets[sheetName] = sheetIndex
		}

		rowCount:=2
		for i:=0;i<len(items);i++{
			f.SetCellValue(dbName, "A"+strconv.Itoa(rowCount), items[i].TableName)
			f.SetSheetRow(dbName, "B"+strconv.Itoa(rowCount), &items[i].ColumnName)
			f.SetCellStyle(dbName, "B"+strconv.Itoa(rowCount), "AJ"+strconv.Itoa(rowCount), fontStyle)

			//Merge table name
			f.MergeCell(dbName, "A"+strconv.Itoa(rowCount), "A"+strconv.Itoa(rowCount+num))
			f.SetCellStyle(dbName, "A"+strconv.Itoa(rowCount), "A"+strconv.Itoa(rowCount+num), centerStyle)
			//store the data from every table
			for j:=1;j<=len(items[i].Rows);j++{
				f.SetSheetRow(dbName, "B"+strconv.Itoa(rowCount+j), &items[i].Rows[j-1])
			}
			rowCount=rowCount+num+3
		}

		for k, v := range categories {
			f.SetCellValue(dbName, k, v)
		}

		err = f.SetCellStyle(dbName, "A1", "A1", style)
		err = f.SetCellStyle(dbName, "B1", "B1", style)
		err = f.SetColWidth(dbName, "A", "A", 20)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// delete the default Sheet1 sheet
	f.DeleteSheet("Sheet1")

	// save the result to xlsx
	currentTime := time.Now().Format("20060102150405")
	filename:=outputID.IP+"_"+outputID.Port+"_"+outputID.User+"_Sample_"+currentTime+".xlsx"

	if err := f.SaveAs(filename); err != nil {
		log.Fatalf(err.Error())
	}
	f.Close()
	log.Println("The results(csv) saved in "+filename)
}

/*
  SensitiveData func: output to csv file
  @Param  csv (the result of SensitiveData function)
  @Param  outputID (the output file name)
*/
func SavetocsvD(csv []models.SensitiveData,outputID InfoStruct){
	categories := map[string]string{
		"A1": "Database", "B1": "Table/Collection","C1":"Data","D1":"Type"}
	var values = make(map[string]string)
	sumTmp:=0
	for i:=0;i<len(csv);i++{
		sumTmp++
		//set excel values
		excelValuetmpA:="A"+strconv.Itoa(sumTmp+1)
		excelValuetmpB:="B"+strconv.Itoa(sumTmp+1)
		excelValuetmpC:="C"+strconv.Itoa(sumTmp+1)
		excelValuetmpD:="D"+strconv.Itoa(sumTmp+1)
		values[excelValuetmpA]= csv[i].DatabaseName
		values[excelValuetmpB]= csv[i].TableName
		values[excelValuetmpC]= csv[i].Data
		values[excelValuetmpD]= csv[i].Type
	}
	//output to a excel
	f := excelize.NewFile()
	err := f.SetColWidth("Sheet1", "A", "C", 50)
	if err!=nil{
		log.Fatalf("SetColWidth goes wrong!!!")
	}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   11,
		},
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = f.SetCellStyle("Sheet1", "A1", "A1", style)
	err = f.SetCellStyle("Sheet1", "B1", "B1", style)
	err = f.SetCellStyle("Sheet1", "C1", "C1", style)
	err = f.SetCellStyle("Sheet1", "D1", "D1", style)
	if err != nil {
		log.Fatalf(err.Error())
	}
	// save the result to xlsx
	currentTime := time.Now().Format("20060102150405")
	filename:=outputID.IP+"_"+outputID.Port+"_"+outputID.User+"_Sensitive_"+currentTime+".xlsx"
	if err := f.SaveAs(filename); err != nil {
		log.Fatalf(err.Error())
	}
	f.Close()
	log.Println("The results(csv) saved in "+filename)
}