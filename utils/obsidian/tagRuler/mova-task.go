package tagRuler

var TagRuleMovaTask = TagRule{
	TagName: "mova-task",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)
		TagRuleTask.CheckRules(manipulator)

		//properties
		cardUrl := "card_url"
		docUrl := "doc_url"

		//Inline
		manipulator.AddPropertyIfNotExist(cardUrl, []string{})
		manipulator.AddPropertyIfNotExist(docUrl, []string{})

		manipulator.IsURI(cardUrl)
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
