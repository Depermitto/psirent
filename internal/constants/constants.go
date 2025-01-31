package constants

import (
	"time"
)

const (
	MaxAddrNum       = 8          // Maximum number of clients that we can download from / can download from us
	FileChunk        = 512 * 1024 // Minimum fragment size of a file in bytes
	PeerPrefix       = "PEER>"    // Prefix for messages from the peer
	HostPrefix       = "HOST>"    // Prefix for messages from the host
	CoordinatorPort  = 6000
	PeerPort         = 6001
	MaxRetryAttempts = 5               // Maximum number of retry attempts for the operation
	RetryDelay       = 2 * time.Second // Delay duration between retry attempts
	ConnectTimeout   = 1 * time.Second // Timeout duration for establishing a connection
)
