package tagRuler

var TagRuleAulaTask = TagRule{
	TagName: "aula-task",
	CheckRules: func(manipulator *FrontmatterManipulator) {
		TagRuleDefault.CheckRules(manipulator)
		TagRuleAula.CheckRules(manipulator)
		TagRuleTask.CheckRules(manipulator)
	},
}
