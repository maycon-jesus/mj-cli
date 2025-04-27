package tagRuler

import "github.com/maycon-jesus/mj-cli/utils/obsidian"

var TagRuleMovaTask = TagRule{
	TagName: "mova-task",
	ApplyRules: func(note *obsidian.ObsidianFile) []string {
		manipulator := FrontmatterManipulator{
			ChMsgs: make(chan []string, 128),
			Note:   note,
		}
		messages := make([]string, 0)

		manipulator.AddPropertyIfNotExist("title", []string{})
		manipulator.AddPropertyIfNotExist("created_at", []string{})
		manipulator.AddPropertyIfNotExist("started_at", []string{})
		manipulator.AddPropertyIfNotExist("finished_at", []string{})
		manipulator.AddPropertyIfNotExist("card_url", []string{})
		manipulator.AddPropertyIfNotExist("doc_url", []string{})
		manipulator.AddPropertyIfNotExist("status", []string{})

		manipulator.IsFilled("title")
		manipulator.IsFilled("created_at")
		manipulator.IsFilled("status")

		manipulator.EnumChecker("status", []string{"created", "doing", "done"})

		if v, _ := note.GetPropertyValues("status"); len(v) > 0 {
			status := v[0]
			switch status {
			case "doing":
				manipulator.IsFilled("started_at")
			case "done":
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
