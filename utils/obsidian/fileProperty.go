package obsidian

type FileProperty struct {
	Key    string   `json:"Key"`
	Values []string `json:"Values"`
}

type FilePropertyName = string
type FilePropertiesMap = map[FilePropertyName]*FileProperty

// GetValues returns a slice of strings representing the Values stored in the FileProperty instance.
func (f FileProperty) GetValues() []string {
	return f.Values
}

// SetValues assigns the provided slice of strings to the Values field of the FileProperty instance.
func (f *FileProperty) SetValues(values []string) {
	f.Values = values
	return
}

// AddValues appends one or more string Values to the Values field of the FileProperty instance.
func (f *FileProperty) AddValues(values ...string) {
	for _, value := range values {
		f.Values = append(f.Values, value)
	}
	return
}
