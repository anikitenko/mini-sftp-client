package controllers

import "io"

type (
	FileStructureStruct struct {
		Path      string
		Directory bool
		Symlink   bool
	}

	PassThru struct {
		io.Reader
		total int64
	}
)
