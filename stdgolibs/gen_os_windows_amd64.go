package stdgolibs

import (
	pkg "os"

	"reflect"
)

func init() {
	registerValues("os", map[string]reflect.Value{
		// Functions
		"Getpagesize":     reflect.ValueOf(pkg.Getpagesize),
		"SameFile":        reflect.ValueOf(pkg.SameFile),
		"MkdirAll":        reflect.ValueOf(pkg.MkdirAll),
		"RemoveAll":       reflect.ValueOf(pkg.RemoveAll),
		"IsPathSeparator": reflect.ValueOf(pkg.IsPathSeparator),
		"Stat":            reflect.ValueOf(pkg.Stat),
		"Lstat":           reflect.ValueOf(pkg.Lstat),
		"ReadDir":         reflect.ValueOf(pkg.ReadDir),
		"Mkdir":           reflect.ValueOf(pkg.Mkdir),
		"Chdir":           reflect.ValueOf(pkg.Chdir),
		"Open":            reflect.ValueOf(pkg.Open),
		"Create":          reflect.ValueOf(pkg.Create),
		"OpenFile":        reflect.ValueOf(pkg.OpenFile),
		"Rename":          reflect.ValueOf(pkg.Rename),
		"TempDir":         reflect.ValueOf(pkg.TempDir),
		"UserCacheDir":    reflect.ValueOf(pkg.UserCacheDir),
		"UserConfigDir":   reflect.ValueOf(pkg.UserConfigDir),
		"UserHomeDir":     reflect.ValueOf(pkg.UserHomeDir),
		"Chmod":           reflect.ValueOf(pkg.Chmod),
		"DirFS":           reflect.ValueOf(pkg.DirFS),
		"ReadFile":        reflect.ValueOf(pkg.ReadFile),
		"WriteFile":       reflect.ValueOf(pkg.WriteFile),
		"Executable":      reflect.ValueOf(pkg.Executable),
		"Getwd":           reflect.ValueOf(pkg.Getwd),
		"Expand":          reflect.ValueOf(pkg.Expand),
		"ExpandEnv":       reflect.ValueOf(pkg.ExpandEnv),
		"Getenv":          reflect.ValueOf(pkg.Getenv),
		"LookupEnv":       reflect.ValueOf(pkg.LookupEnv),
		"Setenv":          reflect.ValueOf(pkg.Setenv),
		"Unsetenv":        reflect.ValueOf(pkg.Unsetenv),
		"Clearenv":        reflect.ValueOf(pkg.Clearenv),
		"Environ":         reflect.ValueOf(pkg.Environ),
		"Hostname":        reflect.ValueOf(pkg.Hostname),
		"Getuid":          reflect.ValueOf(pkg.Getuid),
		"Geteuid":         reflect.ValueOf(pkg.Geteuid),
		"Getgid":          reflect.ValueOf(pkg.Getgid),
		"Getegid":         reflect.ValueOf(pkg.Getegid),
		"Getgroups":       reflect.ValueOf(pkg.Getgroups),
		"Exit":            reflect.ValueOf(pkg.Exit),
		"NewSyscallError": reflect.ValueOf(pkg.NewSyscallError),
		"IsExist":         reflect.ValueOf(pkg.IsExist),
		"IsNotExist":      reflect.ValueOf(pkg.IsNotExist),
		"IsPermission":    reflect.ValueOf(pkg.IsPermission),
		"IsTimeout":       reflect.ValueOf(pkg.IsTimeout),
		"CreateTemp":      reflect.ValueOf(pkg.CreateTemp),
		"MkdirTemp":       reflect.ValueOf(pkg.MkdirTemp),
		"Getpid":          reflect.ValueOf(pkg.Getpid),
		"Getppid":         reflect.ValueOf(pkg.Getppid),
		"FindProcess":     reflect.ValueOf(pkg.FindProcess),
		"StartProcess":    reflect.ValueOf(pkg.StartProcess),
		"NewFile":         reflect.ValueOf(pkg.NewFile),
		"Truncate":        reflect.ValueOf(pkg.Truncate),
		"Remove":          reflect.ValueOf(pkg.Remove),
		"Pipe":            reflect.ValueOf(pkg.Pipe),
		"Link":            reflect.ValueOf(pkg.Link),
		"Symlink":         reflect.ValueOf(pkg.Symlink),
		"Readlink":        reflect.ValueOf(pkg.Readlink),
		"Chown":           reflect.ValueOf(pkg.Chown),
		"Lchown":          reflect.ValueOf(pkg.Lchown),
		"Chtimes":         reflect.ValueOf(pkg.Chtimes),

		// Consts

		"ModeDir":           reflect.ValueOf(pkg.ModeDir),
		"ModeAppend":        reflect.ValueOf(pkg.ModeAppend),
		"ModeExclusive":     reflect.ValueOf(pkg.ModeExclusive),
		"ModeTemporary":     reflect.ValueOf(pkg.ModeTemporary),
		"ModeSymlink":       reflect.ValueOf(pkg.ModeSymlink),
		"ModeDevice":        reflect.ValueOf(pkg.ModeDevice),
		"ModeNamedPipe":     reflect.ValueOf(pkg.ModeNamedPipe),
		"ModeSocket":        reflect.ValueOf(pkg.ModeSocket),
		"ModeSetuid":        reflect.ValueOf(pkg.ModeSetuid),
		"ModeSetgid":        reflect.ValueOf(pkg.ModeSetgid),
		"ModeCharDevice":    reflect.ValueOf(pkg.ModeCharDevice),
		"ModeSticky":        reflect.ValueOf(pkg.ModeSticky),
		"ModeIrregular":     reflect.ValueOf(pkg.ModeIrregular),
		"ModeType":          reflect.ValueOf(pkg.ModeType),
		"ModePerm":          reflect.ValueOf(pkg.ModePerm),
		"PathSeparator":     reflect.ValueOf(pkg.PathSeparator),
		"PathListSeparator": reflect.ValueOf(pkg.PathListSeparator),
		"O_RDONLY":          reflect.ValueOf(pkg.O_RDONLY),
		"O_WRONLY":          reflect.ValueOf(pkg.O_WRONLY),
		"O_RDWR":            reflect.ValueOf(pkg.O_RDWR),
		"O_APPEND":          reflect.ValueOf(pkg.O_APPEND),
		"O_CREATE":          reflect.ValueOf(pkg.O_CREATE),
		"O_EXCL":            reflect.ValueOf(pkg.O_EXCL),
		"O_SYNC":            reflect.ValueOf(pkg.O_SYNC),
		"O_TRUNC":           reflect.ValueOf(pkg.O_TRUNC),
		"SEEK_SET":          reflect.ValueOf(pkg.SEEK_SET),
		"SEEK_CUR":          reflect.ValueOf(pkg.SEEK_CUR),
		"SEEK_END":          reflect.ValueOf(pkg.SEEK_END),
		"DevNull":           reflect.ValueOf(pkg.DevNull),

		// Variables

		"Stdin":               reflect.ValueOf(&pkg.Stdin),
		"Stdout":              reflect.ValueOf(&pkg.Stdout),
		"Stderr":              reflect.ValueOf(&pkg.Stderr),
		"Interrupt":           reflect.ValueOf(&pkg.Interrupt),
		"Kill":                reflect.ValueOf(&pkg.Kill),
		"Args":                reflect.ValueOf(&pkg.Args),
		"ErrInvalid":          reflect.ValueOf(&pkg.ErrInvalid),
		"ErrPermission":       reflect.ValueOf(&pkg.ErrPermission),
		"ErrExist":            reflect.ValueOf(&pkg.ErrExist),
		"ErrNotExist":         reflect.ValueOf(&pkg.ErrNotExist),
		"ErrClosed":           reflect.ValueOf(&pkg.ErrClosed),
		"ErrNoDeadline":       reflect.ValueOf(&pkg.ErrNoDeadline),
		"ErrDeadlineExceeded": reflect.ValueOf(&pkg.ErrDeadlineExceeded),
		"ErrProcessDone":      reflect.ValueOf(&pkg.ErrProcessDone),
	})
	registerTypes("os", map[string]reflect.Type{
		// Non interfaces

		"File":         reflect.TypeOf((*pkg.File)(nil)).Elem(),
		"FileInfo":     reflect.TypeOf((*pkg.FileInfo)(nil)).Elem(),
		"FileMode":     reflect.TypeOf((*pkg.FileMode)(nil)).Elem(),
		"DirEntry":     reflect.TypeOf((*pkg.DirEntry)(nil)).Elem(),
		"LinkError":    reflect.TypeOf((*pkg.LinkError)(nil)).Elem(),
		"ProcessState": reflect.TypeOf((*pkg.ProcessState)(nil)).Elem(),
		"PathError":    reflect.TypeOf((*pkg.PathError)(nil)).Elem(),
		"SyscallError": reflect.TypeOf((*pkg.SyscallError)(nil)).Elem(),
		"Process":      reflect.TypeOf((*pkg.Process)(nil)).Elem(),
		"ProcAttr":     reflect.TypeOf((*pkg.ProcAttr)(nil)).Elem(),
	})
}