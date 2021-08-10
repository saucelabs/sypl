// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package message

import (
	"log"

	"github.com/google/uuid"
	"github.com/saucelabs/sypl/shared"
)

// generateUUID generates UUIDv4 for message ID.
func generateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println(shared.ErrorPrefix, "generateUUID: Failed to generate UUID for message", err)
	}

	return id.String()
}
