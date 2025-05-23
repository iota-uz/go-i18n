package i18n

import (
	"fmt"
	texttemplate "text/template"

	"github.com/iota-uz/go-i18n/v2/i18n/template"
	"github.com/iota-uz/go-i18n/v2/internal"
	"github.com/iota-uz/go-i18n/v2/internal/plural"
)

// MessageTemplate is an executable template for a message.
type MessageTemplate struct {
	*Message
	PluralTemplates map[plural.Form]*internal.Template
}

// NewMessageTemplate returns a new message template.
func NewMessageTemplate(m *Message) *MessageTemplate {
	pluralTemplates := map[plural.Form]*internal.Template{}
	setPluralTemplate(pluralTemplates, plural.Zero, m.Zero, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.One, m.One, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Two, m.Two, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Few, m.Few, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Many, m.Many, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Other, m.Other, m.LeftDelim, m.RightDelim)
	if len(pluralTemplates) == 0 {
		return nil
	}
	return &MessageTemplate{
		Message:         m,
		PluralTemplates: pluralTemplates,
	}
}

func setPluralTemplate(pluralTemplates map[plural.Form]*internal.Template, pluralForm plural.Form, src, leftDelim, rightDelim string) {
	if src != "" {
		pluralTemplates[pluralForm] = &internal.Template{
			Src:        src,
			LeftDelim:  leftDelim,
			RightDelim: rightDelim,
		}
	}
}

type pluralFormNotFoundError struct {
	pluralForm plural.Form
	messageID  string
}

func (e pluralFormNotFoundError) Error() string {
	return fmt.Sprintf("message %q has no plural form %q", e.messageID, e.pluralForm)
}

// Execute executes the template for the plural form and template data.
// Deprecated: This method is no longer used internally by go-i18n and it probably should not have been exported to
// begin with. Its replacement is not exported. If you depend on this method for some reason and/or have
// a use case for exporting execute, please file an issue.
func (mt *MessageTemplate) Execute(pluralForm plural.Form, data interface{}, funcs texttemplate.FuncMap) (string, error) {
	t := mt.PluralTemplates[pluralForm]
	if t == nil {
		return "", pluralFormNotFoundError{
			pluralForm: pluralForm,
			messageID:  mt.ID,
		}
	}
	parser := &template.TextParser{
		Funcs: funcs,
	}
	return t.Execute(parser, data)
}

func (mt *MessageTemplate) execute(pluralForm plural.Form, data interface{}, parser template.Parser) (string, error) {
	t := mt.PluralTemplates[pluralForm]
	if t == nil {
		return "", pluralFormNotFoundError{
			pluralForm: pluralForm,
			messageID:  mt.ID,
		}
	}
	return t.Execute(parser, data)
}
