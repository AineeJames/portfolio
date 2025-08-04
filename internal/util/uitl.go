package util

import (
	"bytes"
	"html/template"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aineejames/portfolio/internal/models"

	"github.com/lmittmann/tint"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

func ParseTemplatesDir(glob string) *template.Template {
	slog.Info("Parsing templates.")
	templates, err := template.ParseGlob(glob)
	if err != nil {
		slog.Error("Failed to parse templates:", tint.Err(err))
		os.Exit(1)
	}
	return templates
}

func CreateDistFolderAndCopy(args ...string) {
	err := os.MkdirAll("dist/projects", 0o0755)
	if err != nil && !os.IsExist(err) {
		slog.Error("Failed to create dist folder:", tint.Err(err))
		os.Exit(1)
	}

	for _, arg := range args {
		cmd := exec.Command("cp", "-r", arg, "dist")
		if err = cmd.Run(); err != nil {
			slog.Error("Failed to copy folder:", tint.Err(err))
			os.Exit(1)
		}
	}
}

func GetPaths(glob string) []string {
	files, err := filepath.Glob(glob)
	if err != nil {
		slog.Error("Failed to find files:", tint.Err(err))
		os.Exit(1)
	}
	return files
}

func NewGoldmark() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&frontmatter.Extender{},
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
			),
		),
	)
}

func GenerateNewProject(templates *template.Template, md goldmark.Markdown, path string) {

	b, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read file:", tint.Err(err))
		os.Exit(1)
	}

	ctx := parser.NewContext()
	var buf bytes.Buffer
	err = md.Convert(b, &buf, parser.WithContext(ctx))
	if err != nil {
		slog.Error("Failed to convert md to html:", tint.Err(err))
		os.Exit(1)
	}

	var meta struct {
		Title string        `yaml:"title"`
		Date  string        `yaml:"date"`
		Tags  []string      `yaml:"tags"`
		Links []models.Link `yaml:"links,omitempty"`
	}
	if err := frontmatter.Get(ctx).Decode(&meta); err != nil {
		slog.Error("Failed to decode frontmatter", tint.Err(err))
		os.Exit(1)
	}

	baseFile := filepath.Base(path)
	fileName := strings.TrimSuffix(baseFile, ".md") + ".html"
	outFile, err := os.Create(filepath.Join("dist", "projects", fileName))
	if err != nil {
		slog.Error("Could not create file:", tint.Err(err))
		os.Exit(1)
	}
	defer outFile.Close()

	err = templates.ExecuteTemplate(outFile, "base", models.SiteData{
		Title: meta.Title,
		Date:  meta.Date,
		Tags:  meta.Tags,
		Links: meta.Links,
		Name:  "Aiden Olsen",
		Year:  time.Now().Year(),
		Body:  template.HTML(buf.String()),
	})
	if err != nil {
		slog.Error("Could not execute template:", tint.Err(err))
		os.Exit(1)
	}
}
