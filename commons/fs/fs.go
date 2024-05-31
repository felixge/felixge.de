package fs

import (
	"io/fs"

	"github.com/felixge/felixge.de/commons/config"
)

func New(embed, os fs.FS, config *config.Config) fs.FS {
	if config.EmbedFS {
		return embed
	}
	// TODO: Wrap os fs to only return files in embed.FS
	return os
}
