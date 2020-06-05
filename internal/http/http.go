package http

import (
	"github.com/antchfx/htmlquery"
)

// FetchImage is a method.
func FetchImage() error{
	doc, err := htmlquery.LoadURL("https://prnt.sc/aaaaab")
	
	if err != nil{
		return err
	}

	nodes, err := htmlquery.QueryAll(doc, "/html/body/div[3]/div/div/img")
	if err != nil {
		panic(`not a valid XPath expression.`)
	}
}