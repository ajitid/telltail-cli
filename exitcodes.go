package main

// do (line num - 4) to get the exit code number
const (
	exitUnsupportedOsVariant = iota + 1
	exitFileNotReadable
	exitFileNotModifiable // In rwx, w includes creation, deletion and makings edits. That's why we've used a word similar to writeable.
	exitDirNotModifiable
	exitUrlNotDownloadable
	exitMissingDependency
	exitCannotDetermineUserHomeDir
	exitInvokingStartupScriptFailed // TODO add this for linux
	exitServiceNotPresent
	exitInvalidConfig
)
