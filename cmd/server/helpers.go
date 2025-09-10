package main

import (
	"errors"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, name string, data any) error {
	v, ok := app.templateCache[name]
	if !ok {
		return errors.New("template not found")
	}

	err := v.ExecuteTemplate(w, "base", data)
	if err != nil {

		return err
	}
	return nil
}
