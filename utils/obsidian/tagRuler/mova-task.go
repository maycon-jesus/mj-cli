package tagRuler

var TagRuleMovaTask = TagRule{
	TagName: "mova-task",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		//Inline
		manipulator.InlineAddPropertyIfNotExist("title", []string{manipulator.Note.Name})
		manipulator.InlineAddPropertyIfNotExist("created_at", []string{})
		manipulator.InlineAddPropertyIfNotExist("started_at", []string{})
		manipulator.InlineAddPropertyIfNotExist("finished_at", []string{})
		manipulator.InlineAddPropertyIfNotExist("card_url", []string{})
		manipulator.InlineAddPropertyIfNotExist("doc_url", []string{})
		manipulator.InlineAddPropertyIfNotExist("status", []string{})

		manipulator.InlineIsFilled("title")
		manipulator.InlineIsFilled("created_at")
		manipulator.InlineIsFilled("status")

		manipulator.InlineEnumChecker("status", []string{"created", "doing", "done"})

		manipulator.InlineIsDate("created_at")
		manipulator.InlineIsDate("started_at")
		manipulator.InlineIsDate("finished_at")

		manipulator.InlineIsURI("card_url")
		manipulator.InlineIsURI("doc_url")

		if status := manipulator.Note.InlineProperties.GetProperty(manipulator.Tag, "status"); status != nil && len(status.Values) > 0 {
			switch status.Values[0] {
			case "doing":
				manipulator.InlineIsFilled("started_at")
			case "done":
				manipulator.InlineIsFilled("started_at")
				manipulator.InlineIsFilled("finished_at")
			}
		}
	},
}
