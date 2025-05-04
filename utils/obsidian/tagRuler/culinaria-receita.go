package tagRuler

var TagRuleCulinariaReceita = TagRule{
	TagName: "culinaria-receita",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)

		categories := "categories"
		capaUrl := "capa_url"

		manipulator.AddPropertyIfNotExist(categories, []string{})
		manipulator.AddPropertyIfNotExist(capaUrl, []string{})

		manipulator.IsFilled(categories)
		manipulator.IsFilled(capaUrl)
	},
}
