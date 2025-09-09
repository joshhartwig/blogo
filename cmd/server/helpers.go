package main

import (
	"errors"
	"net/http"
)

func (c *config) render(w http.ResponseWriter, name string, data any) error {
	v, ok := c.tplCache[name]
	if !ok {
		return errors.New("template not found")
	}

	err := v.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}
	return nil
}
