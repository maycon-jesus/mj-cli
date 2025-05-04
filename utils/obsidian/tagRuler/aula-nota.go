package tagRuler

var TagRuleAulaNota = TagRule{
	TagName: "aula-nota",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)
		TagRuleAula.CheckRules(manipulator)

		date := manipulator.GenPropertyId("date")

		manipulator.AddPropertyIfNotExist(date, []string{})

		manipulator.IsFilled(date)
	},
}
