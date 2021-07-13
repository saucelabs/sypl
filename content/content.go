package content

// Content represents the original, and the processed content. `original` is set
// when the message is created, and can't be changed - just read, while
// `Processed` is consumed and modified by `Processor`s.
type Content struct {
	// original, non-modified content.
	original string

	// processed, consumed and modified by `Processor`s.
	processed string
}

// GetOriginal returns the original, non-modified content.
func (c *Content) GetOriginal() string {
	return c.original
}

// GetProcessed returns the content to be processed.
func (c *Content) GetProcessed() string {
	return c.processed
}

// SetProcessed sets the processed content.
func (c *Content) SetProcessed(content string) {
	c.processed = content
}

// NewContent creates a new content.
func NewContent(content string) *Content {
	return &Content{
		original:  content,
		processed: content,
	}
}
