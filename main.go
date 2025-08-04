package main

import (
	"log/slog"
	"os"

	"github.com/aineejames/portfolio/internal/util"
	"github.com/lmittmann/tint"
)

func init() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, nil)))
}

func main() {
	templates := util.ParseTemplatesDir("templates/*html")
	util.CreateDistFolderAndCopy("static", "assets")
	files := util.GetPaths("content/projects/*.md")
	md := util.NewGoldmark()

	for _, path := range files {
		util.GenerateNewProject(templates, md, path)
	}
}
