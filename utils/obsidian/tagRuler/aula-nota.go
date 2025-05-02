package tagRuler

var TagRuleAulaNota = TagRule{
	TagName: "book",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		manipulator.AddPropertyIfNotExist("title", []string{})
		manipulator.AddPropertyIfNotExist("disciplina", []string{})
		manipulator.AddPropertyIfNotExist("professor", []string{})
		manipulator.AddPropertyIfNotExist("date", []string{})

		manipulator.IsFilled("title")
		manipulator.IsFilled("disciplina")
		manipulator.IsFilled("professor")
		manipulator.IsFilled("date")
	},
}
