package obsidian

type FileProperty struct {
	values   []string
	metadata FilePropertyMetadata
}

type FilePropertyMetadata = map[string]string
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

// GetMetadataValue retrieves the metadata value associated with the provided key and a boolean indicating existence.
func (f FileProperty) GetMetadataValue(key string) (metadataValue string, ok bool) {
	metadataValue, ok = f.metadata[key]
	return metadataValue, ok
}

// SetMetadataValue sets a metadata value for the given key in the FileProperty's metadata map.
func (f *FileProperty) AddMetadata(key string, value string) {
	f.metadata[key] = value
	return
}

func (f *FileProperty) AddMetadataMap(metadata FilePropertyMetadata) {
	for k, v := range metadata {
		f.metadata[k] = v
	}
	return
}

func (f *FileProperty) SetMetadata(metadata FilePropertyMetadata) {
	f.metadata = metadata
	return
}
