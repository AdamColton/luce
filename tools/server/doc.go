// Package server ties together the pieces of the luce server process:
// HTTP routes over tools/server/core, users and sessions through
// util/lusers and lusess, an admin CLI, and the service-registration socket
// that external services connect to.
package server

// == projects.Code.luce.server ==
// [ ] stand alone services #chunk
//  Services can fulfill an interface which can then either be invoked
//  as a service downstream of a luce.server, or built as a standalone service.
// [ ] browser plugin server << stand alone services
