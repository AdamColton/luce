// Package fileservice provides tools to build a service for luce.server that
// handles file requests. A Service maps a BaseURL to a Root Directory. Any
// requests relative to that BaseURL that map to file in the Root Directory
// will be handled by the Service. The Service handles files by their extension.
package fileservice
