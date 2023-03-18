package main

// do (line num - 4) to get the exit code number
const (
	// starts from 2 because generic errors from cli pkg gives 1 as exit status (errors like missing cli arg, or the program returned an error that is not from cli.Exit())
	exitUnsupportedOsVariant = iota + 2
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
