package tagRuler

import "github.com/maycon-jesus/mj-cli/utils/obsidian"

var TagRuleBook = TagRule{
	TagName: "book",
	ApplyRules: func(note *obsidian.ObsidianFile) []string {
		manipulator := FrontmatterManipulator{
			ChMsgs: make(chan []string, 128),
			Note:   note,
		}
		messages := make([]string, 0)

		manipulator.AddPropertyIfNotExist("title", []string{})
		manipulator.AddPropertyIfNotExist("author", []string{})
		manipulator.AddPropertyIfNotExist("started_at", []string{})
		manipulator.AddPropertyIfNotExist("finished_at", []string{})
		manipulator.AddPropertyIfNotExist("status", []string{})

		manipulator.IsFilled("title")
		manipulator.IsFilled("author")
		manipulator.IsFilled("status")

		manipulator.EnumChecker("status", []string{"to-read", "reading", "read"})

		if v, _ := note.GetPropertyValues("status"); len(v) > 0 {
			status := v[0]
			switch status {
			case "reading":
				manipulator.IsFilled("started_at")
			case "read":
				manipulator.IsFilled("started_at")
				manipulator.IsFilled("finished_at")
			}

		}

		close(manipulator.ChMsgs)
		for message := range manipulator.ChMsgs {
			messages = append(messages, message...)
		}

		return messages
	},
}
