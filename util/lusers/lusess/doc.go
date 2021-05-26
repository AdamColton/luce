// Package lusess attaches lusers.User accounts to gorilla/sessions-backed
// HTTP sessions: logging a user in, reading the logged-in user back out of
// a request, and a util/linject field that injects the current Session or
// User into a midware data struct.
package lusess
