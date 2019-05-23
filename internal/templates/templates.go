package templates

import (
	"os"
	"text/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func Write(src, dest string, data interface{}) error {
	box := packr.New("templates", "../../templates")

	text, err := box.FindString(src)
	if err != nil {
		return errors.Wrapf(err, "failed to read template %s", src)
	}

	tmpl, err := template.New(src).Parse(text)
	if err != nil {
		return errors.Wrapf(err, "failed to parse template %s", src)
	}

	file, err := os.Create(dest)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", dest)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		return errors.Wrapf(err, "failed to write template %s to file %s", src, dest)
	}

	return nil
}
