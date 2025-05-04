package tagRuler

var TagRuleAula = TagRule{
	TagName: "aula",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)

		disciplina := "disciplina"
		professor := "professor"

		manipulator.AddPropertyIfNotExist(disciplina, []string{})
		manipulator.AddPropertyIfNotExist(professor, []string{})

		manipulator.IsFilled(disciplina)
		manipulator.IsFilled(professor)
	},
}
