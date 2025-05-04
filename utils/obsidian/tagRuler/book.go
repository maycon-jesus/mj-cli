package tagRuler

var TagRuleBook = TagRule{
	TagName: "book",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)

		manipulator.AddPropertyIfNotExist("author", []string{})
		manipulator.AddPropertyIfNotExist("started_at", []string{})
		manipulator.AddPropertyIfNotExist("finished_at", []string{})
		manipulator.AddPropertyIfNotExist("status", []string{})

		manipulator.IsFilled("author")
		manipulator.IsFilled("status")

		manipulator.EnumChecker("status", []string{"to-read", "reading", "read"})

		if v, _ := manipulator.Note.GetPropertyValues("status"); len(v) > 0 {
			status := v[0]
			switch status {
			case "reading":
				manipulator.IsFilled("started_at")
			case "read":
				manipulator.IsFilled("started_at")
				manipulator.IsFilled("finished_at")
			}

		}
	},
}
