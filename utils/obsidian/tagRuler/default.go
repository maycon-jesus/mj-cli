package tagRuler

import "strconv"

var TagRuleDefault = TagRule{
	TagName: "",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		manipulator.AddPropertyIfNotExist("title", []string{
			manipulator.Note.Name,
		})
		manipulator.IsFilled("title")

		manipulator.AddPropertyIfNotExist("id", []string{
			strconv.FormatInt(manipulator.Note.ModTime, 10),
		})
		manipulator.IsFilled("id")
	},
}
