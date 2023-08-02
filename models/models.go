package models


//initial information
type InitData struct {
	DatabaseType      string
    DatabaseAddress   string
    DatabaseUser      string
    DatabasePassword  string
	DatabaseInstance  string
    ProxyAddress      string
    ProxyUser         string
    ProxyPassword     string
	WindowsAuth       bool       //Use windows authentication to connect database
	AuthSource        string     //Mongodb authentication need to provie database name
}
//Overview data format
type OverviewData struct{
	DatabaseName string
	TableName    string
	RowCount     string
}
type Overviews struct{
	OverviewList []OverviewData
}

//Sample data format
type SampleStruct struct {
	DatabaseName  string
	TableName     string
	ColumnName   []string
	Rows       [][]string
}
type Samples struct{
	SampleList []SampleStruct
}

//Sensitive data format
type SensitiveData struct{
	DatabaseName string
	TableName    string
	Data         string
	Type         string
}
type Sensitive struct{
	SensitiveList []SensitiveData
}

//document struct
type Document struct{
	Key string
	Value string
}
//dccuments array
type Documents []Document
//add the sort function for later sorting the documents array
func (doc Documents) Len() int { return len(doc) }
func (doc Documents) Swap(i, j int) { doc[i], doc[j] = doc[j], doc[i] }
func (doc Documents) Less(i, j int) bool { return doc[i].Key < doc[j].Key }

