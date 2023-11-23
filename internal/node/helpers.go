/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package node

import (
	"fmt"

	"authentication-chains/internal/types"
)

// bucketAuthTableLevel returns the name of the bucket that will store authentication table by level.
func bucketAuthTableLevel(level uint) string {
	return fmt.Sprintf("%s %d", types.BucketAuthenticationTable, level)
}
