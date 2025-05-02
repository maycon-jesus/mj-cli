package obsidian

import (
	"fmt"
	"github.com/maycon-jesus/mj-cli/utils/mySlices"
	"slices"
)

type InlineProperties struct {
	Properties []*InlineProperty `json:"Properties"`
	Tags       []string          `json:"Tags"`
}

func (t *InlineProperties) GetAllPropertiesByTag(tag string) []*InlineProperty {
	return mySlices.Filter(t.Properties, func(p *InlineProperty) bool {
		return p.Tag == tag
	})
}

func (t *InlineProperties) AddProperty(tag string, name string, values []string) []*InlineProperty {
	property := &InlineProperty{
		Tag:    tag,
		Name:   name,
		Values: values,
	}
	t.Properties = append(t.Properties, property)
	if !slices.Contains(t.Tags, tag) {
		t.Tags = append(t.Tags, tag)
	}
	return t.Properties
}

func (t *InlineProperties) GetProperty(tag string, name string) *InlineProperty {
	return t.GetPropertyById(fmt.Sprintf("%s.%s", tag, name))
}

func (t *InlineProperties) GetPropertyById(id string) *InlineProperty {
	for _, p := range t.Properties {
		if p.GetPropertyId() == id {
			return p
		}
	}
	return nil
}

func (inlineProperties *InlineProperties) ToLinesArr() []string {
	lines := make([]string, 0)
	lines = append(lines, "%% start inline-frontmatter %%\n")

	for _, tag := range inlineProperties.Tags {
		lines = append(lines, fmt.Sprintf("- %s", tag))

		properties := inlineProperties.GetAllPropertiesByTag(tag)
		for _, property := range properties {
			if len(property.Values) == 0 {
				lines = append(lines, fmt.Sprintf("\t- %s:", property.Name))
			} else if len(property.Values) == 1 {
				lines = append(lines, fmt.Sprintf("\t- %s: %s", property.Name, property.Values[0]))
			} else {
				lines = append(lines, fmt.Sprintf("\t- %s:", property.Name))
				for _, value := range property.Values {
					lines = append(lines, fmt.Sprintf("\t\t- %s", value))
				}
			}
		}
	}

	lines = append(lines, "\n%% end inline-frontmatter %%")

	return lines
}

func CreateInlineProperties() *InlineProperties {
	return &InlineProperties{
		Properties: make([]*InlineProperty, 0),
		Tags:       make([]string, 0),
	}
}

type InlineProperty struct {
	Tag    string   `json:"Tag"`
	Name   string   `json:"Name"`
	Values []string `json:"Values"`
}

func (p *InlineProperty) GetPropertyId() string {
	propertyId := fmt.Sprintf("%s.%s", p.Tag, p.Name)
	return propertyId
}
