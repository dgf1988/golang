package main

import (
	"net/http"
	"html/template"
)

type HeaderInfo struct{
	Title string
	Description string
	Keywords string
}

func (this HeaderInfo) WriteTo(w http.ResponseWriter, tpl_filename string) error {
	t, err := template.ParseFiles(tpl_filename)
	if err != nil {
		return err
	}
	return t.Execute(w, this)
}

func HtmlHeader(w http.ResponseWriter, tpl_filename string, headerinfo HeaderInfo) error {
	t, err := template.ParseFiles(tpl_filename)
	if err != nil {
		return err
	}
	return t.Execute(w, headerinfo)
}

type FooterInfo struct {
	Addr string
	Email string
	ICP string
	Author string
	Supper string
	CopyRight string
}

func (this FooterInfo) WriteTo(w http.ResponseWriter, tpl_filename string) error {
	t, err := template.ParseFiles(tpl_filename)
	if err != nil {
		return err
	}
	return t.Execute(w, this)
}

func HtmlFooter(w http.ResponseWriter, tpl_filename string, footerInfo FooterInfo) error {
	t, err := template.ParseFiles(tpl_filename)
	if err != nil {
		return err
	}
	return t.Execute(w, footerInfo)
}



type indexPageStep struct {
	IsNow  bool
	Number int64
}

type indexPage struct {
	Steps []indexPageStep
}

func HtmlIndexPage(w http.ResponseWriter, total, now int64) error {
	t, err := template.ParseFiles("hoetom/server/tpl/pageindex.html")
	if err != nil {
		return err
	}
	t.Execute(w, NewIndexPage(total, now))
	return nil
}
