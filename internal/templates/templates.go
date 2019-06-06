package templates

import (
	"io"
	"os"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func WriteFile(src, dest string, data interface{}) error {
	file, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", dest)
	}
	defer file.Close()

	return Write(src, file, data)
}

func Write(src string, dest io.Writer, data interface{}) error {
	box := packr.New("templates", "../../templates")

	text, err := box.FindString(src)
	if err != nil {
		return errors.Wrapf(err, "failed to read template %s", src)
	}

	tmpl, err := template.New(src).Parse(text)
	if err != nil {
		return errors.Wrapf(err, "failed to parse template %s", src)
	}

	err = tmpl.Execute(dest, data)
	if err != nil {
		return errors.Wrapf(err, "failed to write template %s", src)
	}

	return nil
}
