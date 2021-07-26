package templates

import (
	"sync"

	"github.com/monzo/terrors"
)

var (
	mu       sync.RWMutex
	registry = map[SheetType]GooglesheetsTemplate{}
)

func registerTemplate(t SheetType, template GooglesheetsTemplate) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := registry[t]; ok {
		panic("Cannot register the same template twice")
	}

	registry[t] = template
}

func GetTemplateByType(t SheetType) (GooglesheetsTemplate, error) {
	v, ok := registry[t]
	if !ok {
		return nil, terrors.BadRequest("template-does-not-exist", "Template with this type does not exist", map[string]string{
			"template_type": t.String(),
		})
	}

	return v, nil
}
