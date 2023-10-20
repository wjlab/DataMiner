package utils

import (
	"log"
	"os"
	"strings"
	"time"
	"html/template"
	"dataMiner/models"
)

/*
  Overview func: output to html file
  @Param  csv (the result of Overview function)
  @Param  outputID (the output file name)
*/
func SavetohtmlO(csv []models.OverviewData,outputID InfoStruct) {
	csvs := models.Overviews{OverviewList: csv}

	// Define the HTML template
	const tpl = `
	<html> 
	<head> 
    <meta charset="utf-8">
	<title>Database Overview</title> 
	</head> <body> 
	<h1>Database Overview</h1>
	<table border="1">
		 <tr>
		 	<th>Database Name</th>
            <th>Table/Collection Name</th>
            <th>RowsCount/DocumentsCount</th>
		 </tr>

    {{ range $a, $b := $.OverviewList }} 

         <tr> 
            <td>  {{ $b.DatabaseName }}</td>
            <td>  {{ $b.TableName }}</td>
            <td>  {{ $b.RowCount }}</td>
         </tr>
	      
    {{ end }} 

	</table>
	</body> </html>`

	// parse template
	tmpl, err := template.New("overview").Parse(tpl)
	if err != nil {
		log.Fatalf("Error on creating template: %v", err)
	}

	// output html
	currentTime := time.Now().Format("20060102150405")
	filename := outputID.IP + "_" + outputID.Port + "_" + outputID.User + "_Overview_" + currentTime + ".html"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error on creating html: %v", err)
	}
	err = tmpl.Execute(f, csvs)
	if err != nil {
		log.Fatalf("Error on executing template: %v", err)
	}
	log.Println("The results(html) saved in " + filename)
}

/*
  SampleData func: output to html file
  @Param  csv (the result of SampleData function)
  @Param  outputID (the output file name)
*/
func Savetohtml(csv []models.SampleStruct,outputID InfoStruct) {
	// Create a template for the HTML output
	tmpl, err := template.New("table").Funcs(template.FuncMap{
		"join": strings.Join,
	}).Parse(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="UTF-8">
				<title>Sample Data</title>
				<style>
					table {
						border-collapse: collapse;
						margin: 20px 0;
						width: 100%;
					}
					table th, table td {
						padding: 8px;
						text-align: left;
						border: 1px solid #ddd;
					}
					table th {
						background-color: #f2f2f2;
					}
					table tr:nth-child(even) {
						background-color: #f2f2f2;
					}
					table th:nth-child(1), table td:nth-child(1) {
						width: 20%;
					}
					table th:nth-child(2), table td:nth-child(2) {
						width: 30%;
					}
				</style>
			</head>
			<body>
				{{ range $table := . }}
					<h2>{{ $table.DatabaseName }} - {{ $table.TableName }}</h2>
					<table>
						<thead>
							<tr>
                                {{ $numCols := len $table.ColumnName }}
                                {{ if gt $numCols 0 }}
									{{ range $col := slice $table.ColumnName 0 }}
										<th>{{ $col }}</th>
									{{ end }}
                                {{ else }}
									<th>This table/collection is empty.</th>
								{{ end }}
							</tr>
						</thead>
						<tbody>
                           {{ $numRows := len $table.Rows }}
                           {{ if gt $numRows 0 }}
							   {{ range $row := $table.Rows }}
								  <tr>
									{{ range $value := slice $row 0 }}
										<td>{{ $value }}</td>
									{{ end }}
								  </tr>
							   {{ end }}
                           {{ end }}
						</tbody>
					</table>
				{{ end }}
			</body>
		</html>
	`)
	if err != nil {
		log.Fatal(err)
	}

	// output html
	currentTime := time.Now().Format("20060102150405")
	filename:=outputID.IP+"_"+outputID.Port+"_"+outputID.User+"_Sample_"+currentTime+".html"
	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error on creating html: %v", err)
	}
	defer file.Close()

	// Execute the template and write the output to the file
	err = tmpl.Execute(file, csv)
	if err != nil {
		log.Fatalf("Error on executing template: %v", err)
	}
	log.Println("The results(html) saved in "+filename)
}

/*
  SensitiveData func: output to html file
  @Param  csv (the result of SensitiveData function)
  @Param  outputID (the output file name)
*/
func SavetohtmlD(csv []models.SensitiveData,outputID InfoStruct) {
	csvs := models.Sensitive{SensitiveList: csv}
	// Define the HTML template
	const tpl = `
	<html> 
	<head> 
    <meta charset="utf-8">
	<title>Sensitive Data</title> 
	</head> <body> 
	<h1>Sensitive Data</h1>
	<table border="1">
		 <tr>
		 	<th>Database Name</th>
            <th>Table/Collection Name</th>
            <th>Data</th>
            <th>Type</th>
		 </tr>

    {{ range $a, $b := $.SensitiveList }} 

         <tr> 
            <td>  {{ $b.DatabaseName }}</td>
            <td>  {{ $b.TableName }}</td>
            <td>  {{ $b.Data }}</td>
            <td>  {{ $b.Type }}</td>
         </tr>

    {{ end }} 

	</table>
	</body> </html>`

	// parse tmeplate
	tmpl, err := template.New("sensitive").Parse(tpl)
	if err != nil {
		log.Fatalf("Error on creating template: %v", err)
	}

	// output html
	currentTime := time.Now().Format("20060102150405")
	filename := outputID.IP + "_" + outputID.Port + "_" + outputID.User + "_Sensitive_" + currentTime + ".html"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error on creating html: %v", err)
	}
	err = tmpl.Execute(f, csvs)
	if err != nil {
		log.Fatalf("Error on executing template: %v", err)
	}
	log.Println("The results(html) saved in " + filename)
}