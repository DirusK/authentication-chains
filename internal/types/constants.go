/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package types

// InfinityTTL is the value of the infinite TTL.
const InfinityTTL = 0

const (
	// BucketBlocks is the name of the bucket that will store our blocks
	BucketBlocks = "blocks"
	// BucketIndexes is the name of the bucket that will store hashes.
	BucketIndexes = "indexes"
	// BucketAuthenticationTable is the name of the bucket that will store authentication table.
	BucketAuthenticationTable = "authentication-table"
	// BucketCipher is the name of the bucket that will store cipher.
	BucketCipher = "cipher"
	// BucketClusterHead is the name of the bucket that will store cluster head.
	BucketClusterHead = "cluster-head"
	// BucketClusterNodes is the name of the bucket that will store cluster nodes.
	BucketClusterNodes = "cluster-nodes"
	// BucketChildrenNodes is the name of the bucket that will store children nodes.
	BucketChildrenNodes = "children-nodes"
)

var (
	KeyCipher      = []byte("cipher")
	KeyClusterHead = []byte("cluster-head")
	KeyLastBlock   = []byte("last-block")
)
