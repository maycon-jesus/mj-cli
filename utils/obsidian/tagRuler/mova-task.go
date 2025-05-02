package tagRuler

var TagRuleMovaTask = TagRule{
	TagName: "mova-task",
	CheckRules: func(manipulator *FrontmatterManipulator) {
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

		if v, _ := manipulator.Note.GetPropertyValues("status"); len(v) > 0 {
			status := v[0]
			switch status {
			case "doing":
				manipulator.IsFilled("started_at")
			case "done":
				manipulator.IsFilled("started_at")
				manipulator.IsFilled("finished_at")
			}

		}
	},
}
