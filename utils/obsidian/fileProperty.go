package obsidian

type FileProperty struct {
	values []string
}

type FilePropertyName = string
type FilePropertiesMap = map[FilePropertyName]*FileProperty

// GetValues returns a slice of strings representing the values stored in the FileProperty instance.
func (f FileProperty) GetValues() []string {
	return f.values
}

// SetValues assigns the provided slice of strings to the values field of the FileProperty instance.
func (f *FileProperty) SetValues(values []string) {
	f.values = values
	return
}

// AddValues appends one or more string values to the values field of the FileProperty instance.
func (f *FileProperty) AddValues(values ...string) {
	for _, value := range values {
		f.values = append(f.values, value)
	}
	return
}
