# lualatex-serienbrief

Siehe auch [lualatex-brief](https://github.com/jojomi/lualatex-brief) für eine passende Briefvorlage.

Die Template-Syntax von [go](https://golang.org) kann in `tex`- und `lco`-Dateien innerhalb des Template-Ordners verwendet werden. Ein kleines Tutorial findet sich [hier](https://themes.gohugo.io/theme/minimal/post/goisforlovers/). Verfügbar sind dort alle Spalten aus der `csv`-Datei, die als Datenquelle dient. Diese benötigt eine Kopfzeile mit Spaltentiteln.


## Verwendung

```
lualatex-serienbrief --help
Usage:
  lualatex-serienbrief [flags]

Flags:
  -d, --data-file string              data file (default is data.csv) (default "data.csv")
  -h, --help                          help for lualatex-serienbrief
  -f, --output-file-template string   template for output PDF (without extension, default is {{ .Name }}) (default "{{.Name }}")
  -o, --output-folder string          output directory (default is output) (default "output")
  -t, --template-dir string           template directory (default is template) (default "template")
  -l, --tex-file string               tex file (default is main.tex) (default "main.tex")
  -v, --verbose                       show more output (e.g. while compiling document)
```