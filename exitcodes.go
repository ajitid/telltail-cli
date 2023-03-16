package main

// do (line num - 4) to get the exit code number
const (
	exitUnsupportedOs = 1
	exitUnsupportedOsVariant
	exitFileNotWriteable // In rwx, w includes creation, deletion and modification. That's why we've used the word 'writeable'
	exitDirNotCreatable
	exitUrlNotDownloadable
	exitMissingDependency
	exitCannotDetermineUserHomeDir
)
