package main

import (
	"log"
	"os"
	"regexp"
	"sort"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/jszwec/csvutil"
)

type Record struct {
	DeptCd string `csv:"deptCd"`
	EmpNm  string `csv:"empNm"`
	OrigNm string `csv:"origNm"`
	PolyNm string `csv:"polyNm"`
	Issues string
}

type issue struct {
	ID    string `csv:"id"`
	Name  string `csv:"name"`
	Issue string `csv:"issue"`
}

func main() {
	var records []Record
	{
		_data, err := os.ReadFile("assembly.csv")
		_ = _data
		if err != nil {
			panic(err)
		}
		err = csvutil.Unmarshal(_data, &records)
		if err != nil {
			panic(err)
		}
		log.Println(records)
	}

	var issueMap = make(map[string][]string)
	var issues []issue

	{
		_assembly_etc, err := os.ReadFile("assembly_issue.csv")
		if err != nil {
			panic(err)
		}
		err = csvutil.Unmarshal(_assembly_etc, &issues)
		if err != nil {
			panic(err)
		}

		for _, issue := range issues {
			issueMap[issue.ID] = append(issueMap[issue.ID], issue.Issue)
		}
	}

	for i, record := range records {
		if issues, ok := issueMap[record.DeptCd]; ok {
			for _, issue := range issues {
				records[i].Issues = records[i].Issues + markdownToHTML(issue) + " "
			}
		}
	}

	for _, record := range records {
		spew.Dump(record)
	}

	// sort the records by length of issues
	sort.Slice(records, func(i, j int) bool {
		return len(records[i].Issues) > len(records[j].Issues)
	})

	tmpl, err := template.ParseFiles("templates/index.html.tmpl")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("out", os.ModePerm)
	if err != nil {
		panic(err)
	}

	output, err := os.Create("out/index.html")
	if err != nil {
		panic(err)
	}
	defer output.Close()

	err = tmpl.Execute(output, records)
	if err != nil {
		panic(err)
	}
}
func markdownToHTML(markdown string) string {
	re := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	html := re.ReplaceAllString(markdown, `<a href="$2">$1</a>`)
	return html
}
