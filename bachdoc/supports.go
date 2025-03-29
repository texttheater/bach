package main

import (
	"errors"
)

type supportsCmd struct {
	RendererName string `arg:"" help:"name of the renderer"`
}

func (s *supportsCmd) Run() error {
	if s.RendererName != "html" {
		return errors.New("renderer not supported")
	}
	return nil
}
