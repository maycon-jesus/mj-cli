package obsidian

import (
	"fmt"
	"slices"
	"strings"
)

type FrontmatterManipulator struct {
	ChMsgs chan []string
	Note   *ObsidianFile
}

func (m *FrontmatterManipulator) addPropertyIfNotExist(propertyName string, propertyValues []string) {
	messages := make([]string, 0)
	defer func() {
		m.ChMsgs <- messages
	}()

	_, ok := m.Note.Frontmatter[propertyName]
	if ok {
		return
	}
	m.Note.AddProperty(propertyName, propertyValues)
	messages = append(messages, fmt.Sprintf("Adicionada propriedade %s", propertyName))
}

func (m *FrontmatterManipulator) enumChecker(propertyName string, expectedEnum []string) {
	messages := make([]string, 0)
	defer func() {
		m.ChMsgs <- messages
	}()

	property, ok := m.Note.Frontmatter[propertyName]
	if !ok {
		return
	}
	for i, value := range property.GetValues() {
		if !slices.Contains(expectedEnum, value) {
			message := fmt.Sprintf("O valor %d da propriedade %s está fora do enum: %s", i, propertyName, strings.Join(expectedEnum, " | "))
			messages = append(messages, message)
		}
	}
}

func (m *FrontmatterManipulator) isFilled(propertyName string) {
	messages := make([]string, 0)
	defer func() {
		m.ChMsgs <- messages
	}()

	property, ok := m.Note.Frontmatter[propertyName]
	if !ok {
		return
	}
	if len(property.GetValues()) == 0 {
		message := fmt.Sprintf("A propriedade %s precisa estar preenchida", propertyName)
		messages = append(messages, message)
	}
}

var TagsRules = map[string]TagRule{
	"mova-task": TagRuleMovaTask,
}

type TagRule struct {
	TagName    string
	ApplyRules func(note *ObsidianFile) []string
}

var TagRuleMovaTask = TagRule{
	TagName: "mova-task",
	ApplyRules: func(note *ObsidianFile) []string {
		manipulator := FrontmatterManipulator{
			ChMsgs: make(chan []string, 128),
			Note:   note,
		}
		messages := make([]string, 0)

		manipulator.addPropertyIfNotExist("title", []string{})
		manipulator.addPropertyIfNotExist("created_at", []string{})
		manipulator.addPropertyIfNotExist("started_at", []string{})
		manipulator.addPropertyIfNotExist("finished_at", []string{})
		manipulator.addPropertyIfNotExist("card_url", []string{})
		manipulator.addPropertyIfNotExist("doc_url", []string{})
		manipulator.addPropertyIfNotExist("status", []string{})

		manipulator.isFilled("title")
		manipulator.isFilled("created_at")
		manipulator.isFilled("status")

		manipulator.enumChecker("status", []string{"created", "doing", "done"})

		if v, _ := note.GetPropertyValues("status"); len(v) > 0 {
			status := v[0]
			switch status {
			case "doing":
				manipulator.isFilled("started_at")
			case "done":
				manipulator.isFilled("started_at")
				manipulator.isFilled("finished_at")
			}

		}

		close(manipulator.ChMsgs)
		for message := range manipulator.ChMsgs {
			messages = append(messages, message...)
		}

		return messages
	},
}
