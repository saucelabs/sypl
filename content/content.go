// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package content

// IContent specifies what a content does.
type IContent interface {
	// GetOriginal returns the original, non-modified content.
	GetOriginal() string

	// GetProcessed returns the content to be processed.
	GetProcessed() string

	// SetProcessed sets the processed content.
	SetProcessed(content string)
}

// Content represents the original, and the processed content. original is set
// when the message is created, and can't be changed - just read, while
// Processed is consumed and modified by Processors.
type content struct {
	// original, non-modified content.
	Original string

	// processed, consumed and modified by Processors.
	Processed string
}

//////
// IContent interface implementation.
//////

// GetOriginal returns the original, non-modified content.
func (c *content) GetOriginal() string {
	return c.Original
}

// GetProcessed returns the content to be processed.
func (c *content) GetProcessed() string {
	return c.Processed
}

// SetProcessed sets the processed content.
func (c *content) SetProcessed(content string) {
	c.Processed = content
}

//////
// Factory.
//////

// New is the Content factory.
func New(c string) IContent {
	return &content{
		Original:  c,
		Processed: c,
	}
}
