package backend

import (
	"fmt"

	"github.com/nedrocks/delphisbe/graph/model"
)

type mockImportedContentIter struct{}

func (m *mockImportedContentIter) Next(content *model.ImportedContent) bool { return true }
func (m *mockImportedContentIter) Close() error                             { return fmt.Errorf("error") }
