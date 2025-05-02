package tagRuler

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils/obsidian"
	"slices"
	"strings"
)

type FrontmatterManipulator struct {
	ChMsgs chan []string
	Note   *obsidian.ObsidianFile
}

func (m *FrontmatterManipulator) ReadAllMessagesInChannel() []string {
	storeMessages := make([]string, 0)
	for message := range m.ChMsgs {
		for _, msg := range message {
			storeMessages = append(storeMessages, msg)
		}
	}
	return storeMessages
}

func (m *FrontmatterManipulator) AddPropertyIfNotExist(propertyName string, propertyValues []string) {
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

func (m *FrontmatterManipulator) EnumChecker(propertyName string, expectedEnum []string) {
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

func (m *FrontmatterManipulator) IsFilled(propertyName string) {
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

func NewFrontmatterManipulator(note *obsidian.ObsidianFile) *FrontmatterManipulator {
	return &FrontmatterManipulator{
		ChMsgs: make(chan []string, 128),
		Note:   note,
	}
}

type TagRule struct {
	TagName    string
	CheckRules func(manipulator *FrontmatterManipulator)
}

func (t TagRule) ApplyRules(note *obsidian.ObsidianFile) []string {
	manipulator := NewFrontmatterManipulator(note)
	t.CheckRules(manipulator)
	close(manipulator.ChMsgs)

	messages := manipulator.ReadAllMessagesInChannel()

	return messages
}

var TagsRules = map[string]TagRule{
	"mova-task":   TagRuleMovaTask,
	"book":        TagRuleBook,
	"aula/nota":   TagRuleAulaNota,
	"aula/tarefa": TagRuleAulaTarefa,
}
