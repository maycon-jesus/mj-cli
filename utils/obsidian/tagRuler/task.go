package tagRuler

var TagRuleTask = TagRule{
	TagName: "task",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)

		//properties
		createdAt := "created_at"
		startedAt := "started_at"
		finishedAt := "finished_at"
		status := "status"

		//Inline
		manipulator.AddPropertyIfNotExist(createdAt, []string{})
		manipulator.AddPropertyIfNotExist(startedAt, []string{})
		manipulator.AddPropertyIfNotExist(finishedAt, []string{})
		manipulator.AddPropertyIfNotExist(status, []string{})

		manipulator.IsFilled(createdAt)
		manipulator.IsFilled(status)

		manipulator.EnumChecker(status, []string{"created", "doing", "done"})

		//manipulator.InlineIsDate("created_at")
		//manipulator.InlineIsDate("started_at")
		//manipulator.InlineIsDate("finished_at")
		//
		//manipulator.InlineIsURI("card_url")
		//manipulator.InlineIsURI("doc_url")

		//if status := manipulator.Note.InlineProperties.GetProperty(manipulator.Tag, "status"); status != nil && len(status.Values) > 0 {
		//	switch status.Values[0] {
		//	case "doing":
		//		manipulator.InlineIsFilled("started_at")
		//	case "done":
		//		manipulator.InlineIsFilled("started_at")
		//		manipulator.InlineIsFilled("finished_at")
		//	}
		//}
	},
}
