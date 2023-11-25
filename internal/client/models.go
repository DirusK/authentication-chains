/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

type (
	DeviceAuthenticationRequest struct {
		DeviceID      string
		ClusterHeadID string
		Signature     string
	}
)
