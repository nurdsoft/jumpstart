package pkg

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

func renderTemplates(ctx pongo2.Context, inDir, outDir string) error {
	return filepath.WalkDir(inDir, func(path string, d fs.DirEntry, e error) error {
		if d.IsDir() {
			return nil
		}
		relPath, err := filepath.Rel(inDir, path)
		if err != nil {
			return err
		}
		outPath := filepath.Join(outDir, relPath)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		template, err := pongo2.FromFile(path)
		if err != nil {
			return err
		}
		output, err := template.ExecuteBytes(ctx)
		if err != nil {
			return err
		}
		return os.WriteFile(outPath, output, 0644)
	})
}

func renderConfig(ctx pongo2.Context, inFile string) ([]byte, error) {
	template, err := pongo2.FromFile(inFile)
	if err != nil {
		return nil, err
	}
	return template.ExecuteBytes(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// outFile := filepath.Join(outDir, filepath.Base(inFile))
	// return os.WriteFile(outFile, output, 0644)
}
