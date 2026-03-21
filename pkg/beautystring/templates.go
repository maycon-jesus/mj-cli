package beautystring

// === STATUS ===

func (sb *StrBuilder) Success(message string) *StrBuilder {
	sb.Green().Bold().Text("✔ " + message).Reset()
	return sb
}

func (sb *StrBuilder) Failure(message string) *StrBuilder {
	sb.Red().Bold().Text("✖ " + message).Reset()
	return sb
}

func (sb *StrBuilder) Warning(message string) *StrBuilder {
	sb.Yellow().Bold().Text("⚠ " + message).Reset()
	return sb
}

func (sb *StrBuilder) Info(message string) *StrBuilder {
	sb.Blue().Bold().Text("ℹ " + message).Reset()
	return sb
}

func (sb *StrBuilder) Pending(message string) *StrBuilder {
	sb.Cyan().Bold().Text("… " + message).Reset()
	return sb
}

func (sb *StrBuilder) Lambda(message string) *StrBuilder {
	sb.Magenta().Bold().Text("λ " + message).Reset()
	return sb
}

func (sb *StrBuilder) Lambdaf(format string, a ...any) *StrBuilder {
	sb.Magenta().Bold().Textf("λ "+format, a...).Reset()
	return sb
}

// === TITULOS E SECOES ===

func (sb *StrBuilder) TitleLine(message string) *StrBuilder {
	sb.Underline().Bold().Text(message).NewLine().Reset()
	return sb
}

func (sb *StrBuilder) TitleLinef(format string, a ...any) *StrBuilder {
	sb.Underline().Bold().Textf(format, a...).NewLine().Reset()
	return sb
}

func (sb *StrBuilder) SectionLine(message string) *StrBuilder {
	sb.Bold().Text(message).Reset().NewLine().Reset()
	return sb
}

// === LISTAS E TABELAS

func (sb *StrBuilder) List(items []string) *StrBuilder {
	for _, item := range items {
		sb.Text("• " + item).NewLine()
	}
	return sb
}

func (sb *StrBuilder) NumberedList(items []string) *StrBuilder {
	for i, item := range items {
		sb.Textf("%d. %s", i+1, item).NewLine()
	}
	return sb
}
