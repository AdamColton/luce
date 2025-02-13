// Package server ties together the pieces of the luce server process:
// HTTP routes over tools/server/core, users and sessions through
// util/lusers and lusess, an admin CLI, and the service-registration socket
// that external services connect to.
package server
