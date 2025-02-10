// Package lfilemock provides mock file objects and a mock file system, so code
// that uses lfile can be tested without touching the disk. Parse builds a tree
// from a map, and Directory.Repository serves it as an lfile.Repository.
//
// The mock behaves as the os package does where that matters to a caller: a
// missing path is an *fs.PathError that is fs.ErrNotExist, a directory can't be
// read as a file and a file can't be listed, a directory listing is in order and
// is read in pages, and a file is opened by the name it is given. Every node
// has an Err that is returned by its operations, to test error handling, and
// Mode, Mod and SysData that Stat reports.
//
// It is not meant to be exhaustive, it is used to cover the cases needed to test
// luce. There are no links, and the tree is not safe for concurrent changes.
package lfilemock
