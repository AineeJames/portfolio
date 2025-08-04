package models

import (
	"html/template"
)

type Link struct {
	Name string `yaml:"name"`
	Link string `yaml:"link"`
}

type SiteData struct {
	Year  int
	Title string
	Body  template.HTML
	Date  string
	Tags  []string
	Links []Link
	Name  string
}
