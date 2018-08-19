package main

import (
	"testing"
	"io/ioutil"
	"strings"
	"path/filepath"
	"bytes"
	"html/template"
	"os"
)

func TestRenderSGV(t *testing.T) {
	runTestsWithRenderer(t, "svg", SvgRenderer)
}

func TestRenderPNG(t *testing.T) {
	runTestsWithRenderer(t, "png", PngRenderer)
}

func runTestsWithRenderer(t *testing.T, outFileExt string, renderer Renderer) {
	tmpDir, err := ioutil.TempDir("", "goseq-test-"+outFileExt)
	if err != nil {
		t.Errorf("cannot create temporary directory '%s': %v", tmpDir, err)
	}

	t.Logf("%s tmpdir = %s", outFileExt, tmpDir)

	inFiles := testFiles(t)
	outFiles := make([]string, 0, len(inFiles))

	for _, f := range inFiles {
		t.Run(f, func(t *testing.T) {
			content, err := ioutil.ReadFile(f)
			if err != nil {
				t.Errorf("cannot read file '%s': %v", f, err)
			}

			outFile := filepath.Join(tmpDir, filepath.Base(f)+"."+outFileExt)
			outFiles = append(outFiles, outFile)

			if err := ProcessSeqDiagram(bytes.NewBuffer(content), f, outFile, renderer); err != nil {
				t.Errorf("cannot generate '%s' -> '%s': %v", f, outFile, err)
			}
		})
	}

	generateOutputBrowserPage(t, tmpDir, inFiles, outFiles)
}

func testFiles(t *testing.T) []string {
	dir, err := ioutil.ReadDir("_tests")
	if err != nil {
		t.Fatalf("error opening test dir: %v", err)
	}

	filenames := make([]string, 0, len(dir))
	for _, f := range dir {
		if strings.HasSuffix(f.Name(), ".seq") {
			filenames = append(filenames, filepath.Join("_tests", f.Name()))
		}
	}

	return filenames
}

func generateOutputBrowserPage(t *testing.T, tmpDir string, inFiles []string, outFiles []string) {
	type fileData struct {
		SeqFileName    string
		SeqFileContent string
		ImgFileName    string
	}

	tmpl := template.Must(template.New("outputBrowser").Parse(outputBrowserPage))
	files := make([]fileData, len(inFiles))
	for i := range inFiles {
		content, err := ioutil.ReadFile(inFiles[i])
		if err != nil {
			t.Errorf("cannot read file '%s': %v", inFiles[i], err)
		}

		files[i] = fileData{
			SeqFileName:    inFiles[i],
			SeqFileContent: string(content),
			ImgFileName:    outFiles[i],
		}
	}

	browserOut, err := os.Create(filepath.Join(tmpDir, "out.html"))
	if err != nil {
		t.Fatalf("cannot generate browser out.html: %v", err)
	}

	if err := tmpl.Execute(browserOut, struct {
		TestFiles []fileData
	}{TestFiles: files}); err != nil {
		browserOut.Close()
		t.Fatalf("cannot generate browser out.html: %v", err)
	}

	if err := browserOut.Close(); err != nil {
		t.Fatalf("cannot generate browser out.html: %v", err)
	}
}

const outputBrowserPage = `<!DOCTYPE html>
<html>
<head>
  <style>
    table { border: solid thin black; border-collapse: collapse; }
    td { border: solid; }
  </style>
</head>
<body>
  {{ range .TestFiles }}
    <p>{{ .SeqFileName }}</p>

	<table>
	  <tr>
		<td>
          <pre>{{ .SeqFileContent }}</pre>
        </td>
        <td>
          <img src="{{ .ImgFileName }}">
        </td>
      </tr>
	</table>
  {{ end }}
</body>
</html>`
