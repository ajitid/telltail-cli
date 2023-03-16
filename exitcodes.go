package main

// do (line num - 4) to get the exit code number
const (
	exitFileNotWriteable = 1 // In rwx, w includes creation, deletion and modification. That's why we've used the word 'writeable'.
	exitDirNotCreatable
	exitUrlNotDownloadable
	exitMissingDependency
	exitCannotDetermineUserHomeDir
	exitInvokingStartupScriptFailed // TODO add this for linux
)
