package main

// do (line num - 5) to get the exit status number
const (
	// --- Intentionally Left Empty (to make calculating exit status using line number easy)
	// ---
	// ---
	// starts from 5 because generic errors from cli pkg gives 1 as exit status (errors like missing cli arg, or the program returned an error that is not from cli.Exit())
	// and unknown command gives 3
	exitUnsupportedOsVariant = iota + 5
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
