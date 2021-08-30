// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import "errors"

var ErrSyplNotInitialized = errors.New("sypl isn't initialized. Have you instantiated it?")
