package stdgolibs

import (
	pkg "debug/pe"

	"reflect"
)

func init() {
	registerValues("debug/pe", map[string]reflect.Value{
		// Functions
		"Open":    reflect.ValueOf(pkg.Open),
		"NewFile": reflect.ValueOf(pkg.NewFile),

		// Consts

		"COFFSymbolSize":                                 reflect.ValueOf(pkg.COFFSymbolSize),
		"IMAGE_FILE_MACHINE_UNKNOWN":                     reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_UNKNOWN),
		"IMAGE_FILE_MACHINE_AM33":                        reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_AM33),
		"IMAGE_FILE_MACHINE_AMD64":                       reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_AMD64),
		"IMAGE_FILE_MACHINE_ARM":                         reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_ARM),
		"IMAGE_FILE_MACHINE_ARMNT":                       reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_ARMNT),
		"IMAGE_FILE_MACHINE_ARM64":                       reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_ARM64),
		"IMAGE_FILE_MACHINE_EBC":                         reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_EBC),
		"IMAGE_FILE_MACHINE_I386":                        reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_I386),
		"IMAGE_FILE_MACHINE_IA64":                        reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_IA64),
		"IMAGE_FILE_MACHINE_M32R":                        reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_M32R),
		"IMAGE_FILE_MACHINE_MIPS16":                      reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_MIPS16),
		"IMAGE_FILE_MACHINE_MIPSFPU":                     reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_MIPSFPU),
		"IMAGE_FILE_MACHINE_MIPSFPU16":                   reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_MIPSFPU16),
		"IMAGE_FILE_MACHINE_POWERPC":                     reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_POWERPC),
		"IMAGE_FILE_MACHINE_POWERPCFP":                   reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_POWERPCFP),
		"IMAGE_FILE_MACHINE_R4000":                       reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_R4000),
		"IMAGE_FILE_MACHINE_SH3":                         reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_SH3),
		"IMAGE_FILE_MACHINE_SH3DSP":                      reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_SH3DSP),
		"IMAGE_FILE_MACHINE_SH4":                         reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_SH4),
		"IMAGE_FILE_MACHINE_SH5":                         reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_SH5),
		"IMAGE_FILE_MACHINE_THUMB":                       reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_THUMB),
		"IMAGE_FILE_MACHINE_WCEMIPSV2":                   reflect.ValueOf(pkg.IMAGE_FILE_MACHINE_WCEMIPSV2),
		"IMAGE_DIRECTORY_ENTRY_EXPORT":                   reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_EXPORT),
		"IMAGE_DIRECTORY_ENTRY_IMPORT":                   reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_IMPORT),
		"IMAGE_DIRECTORY_ENTRY_RESOURCE":                 reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_RESOURCE),
		"IMAGE_DIRECTORY_ENTRY_EXCEPTION":                reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_EXCEPTION),
		"IMAGE_DIRECTORY_ENTRY_SECURITY":                 reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_SECURITY),
		"IMAGE_DIRECTORY_ENTRY_BASERELOC":                reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_BASERELOC),
		"IMAGE_DIRECTORY_ENTRY_DEBUG":                    reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_DEBUG),
		"IMAGE_DIRECTORY_ENTRY_ARCHITECTURE":             reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_ARCHITECTURE),
		"IMAGE_DIRECTORY_ENTRY_GLOBALPTR":                reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_GLOBALPTR),
		"IMAGE_DIRECTORY_ENTRY_TLS":                      reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_TLS),
		"IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG":              reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG),
		"IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT":             reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT),
		"IMAGE_DIRECTORY_ENTRY_IAT":                      reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_IAT),
		"IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT":             reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT),
		"IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR":           reflect.ValueOf(pkg.IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR),
		"IMAGE_FILE_RELOCS_STRIPPED":                     reflect.ValueOf(pkg.IMAGE_FILE_RELOCS_STRIPPED),
		"IMAGE_FILE_EXECUTABLE_IMAGE":                    reflect.ValueOf(pkg.IMAGE_FILE_EXECUTABLE_IMAGE),
		"IMAGE_FILE_LINE_NUMS_STRIPPED":                  reflect.ValueOf(pkg.IMAGE_FILE_LINE_NUMS_STRIPPED),
		"IMAGE_FILE_LOCAL_SYMS_STRIPPED":                 reflect.ValueOf(pkg.IMAGE_FILE_LOCAL_SYMS_STRIPPED),
		"IMAGE_FILE_AGGRESIVE_WS_TRIM":                   reflect.ValueOf(pkg.IMAGE_FILE_AGGRESIVE_WS_TRIM),
		"IMAGE_FILE_LARGE_ADDRESS_AWARE":                 reflect.ValueOf(pkg.IMAGE_FILE_LARGE_ADDRESS_AWARE),
		"IMAGE_FILE_BYTES_REVERSED_LO":                   reflect.ValueOf(pkg.IMAGE_FILE_BYTES_REVERSED_LO),
		"IMAGE_FILE_32BIT_MACHINE":                       reflect.ValueOf(pkg.IMAGE_FILE_32BIT_MACHINE),
		"IMAGE_FILE_DEBUG_STRIPPED":                      reflect.ValueOf(pkg.IMAGE_FILE_DEBUG_STRIPPED),
		"IMAGE_FILE_REMOVABLE_RUN_FROM_SWAP":             reflect.ValueOf(pkg.IMAGE_FILE_REMOVABLE_RUN_FROM_SWAP),
		"IMAGE_FILE_NET_RUN_FROM_SWAP":                   reflect.ValueOf(pkg.IMAGE_FILE_NET_RUN_FROM_SWAP),
		"IMAGE_FILE_SYSTEM":                              reflect.ValueOf(pkg.IMAGE_FILE_SYSTEM),
		"IMAGE_FILE_DLL":                                 reflect.ValueOf(pkg.IMAGE_FILE_DLL),
		"IMAGE_FILE_UP_SYSTEM_ONLY":                      reflect.ValueOf(pkg.IMAGE_FILE_UP_SYSTEM_ONLY),
		"IMAGE_FILE_BYTES_REVERSED_HI":                   reflect.ValueOf(pkg.IMAGE_FILE_BYTES_REVERSED_HI),
		"IMAGE_SUBSYSTEM_UNKNOWN":                        reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_UNKNOWN),
		"IMAGE_SUBSYSTEM_NATIVE":                         reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_NATIVE),
		"IMAGE_SUBSYSTEM_WINDOWS_GUI":                    reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_WINDOWS_GUI),
		"IMAGE_SUBSYSTEM_WINDOWS_CUI":                    reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_WINDOWS_CUI),
		"IMAGE_SUBSYSTEM_OS2_CUI":                        reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_OS2_CUI),
		"IMAGE_SUBSYSTEM_POSIX_CUI":                      reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_POSIX_CUI),
		"IMAGE_SUBSYSTEM_NATIVE_WINDOWS":                 reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_NATIVE_WINDOWS),
		"IMAGE_SUBSYSTEM_WINDOWS_CE_GUI":                 reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_WINDOWS_CE_GUI),
		"IMAGE_SUBSYSTEM_EFI_APPLICATION":                reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_EFI_APPLICATION),
		"IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER":        reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_EFI_BOOT_SERVICE_DRIVER),
		"IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER":             reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_EFI_RUNTIME_DRIVER),
		"IMAGE_SUBSYSTEM_EFI_ROM":                        reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_EFI_ROM),
		"IMAGE_SUBSYSTEM_XBOX":                           reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_XBOX),
		"IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION":       reflect.ValueOf(pkg.IMAGE_SUBSYSTEM_WINDOWS_BOOT_APPLICATION),
		"IMAGE_DLLCHARACTERISTICS_HIGH_ENTROPY_VA":       reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_HIGH_ENTROPY_VA),
		"IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE":          reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_DYNAMIC_BASE),
		"IMAGE_DLLCHARACTERISTICS_FORCE_INTEGRITY":       reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_FORCE_INTEGRITY),
		"IMAGE_DLLCHARACTERISTICS_NX_COMPAT":             reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_NX_COMPAT),
		"IMAGE_DLLCHARACTERISTICS_NO_ISOLATION":          reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_NO_ISOLATION),
		"IMAGE_DLLCHARACTERISTICS_NO_SEH":                reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_NO_SEH),
		"IMAGE_DLLCHARACTERISTICS_NO_BIND":               reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_NO_BIND),
		"IMAGE_DLLCHARACTERISTICS_APPCONTAINER":          reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_APPCONTAINER),
		"IMAGE_DLLCHARACTERISTICS_WDM_DRIVER":            reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_WDM_DRIVER),
		"IMAGE_DLLCHARACTERISTICS_GUARD_CF":              reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_GUARD_CF),
		"IMAGE_DLLCHARACTERISTICS_TERMINAL_SERVER_AWARE": reflect.ValueOf(pkg.IMAGE_DLLCHARACTERISTICS_TERMINAL_SERVER_AWARE),

		// Variables

	})
	registerTypes("debug/pe", map[string]reflect.Type{
		// Non interfaces

		"StringTable":      reflect.TypeOf((*pkg.StringTable)(nil)).Elem(),
		"COFFSymbol":       reflect.TypeOf((*pkg.COFFSymbol)(nil)).Elem(),
		"Symbol":           reflect.TypeOf((*pkg.Symbol)(nil)).Elem(),
		"File":             reflect.TypeOf((*pkg.File)(nil)).Elem(),
		"ImportDirectory":  reflect.TypeOf((*pkg.ImportDirectory)(nil)).Elem(),
		"FormatError":      reflect.TypeOf((*pkg.FormatError)(nil)).Elem(),
		"FileHeader":       reflect.TypeOf((*pkg.FileHeader)(nil)).Elem(),
		"DataDirectory":    reflect.TypeOf((*pkg.DataDirectory)(nil)).Elem(),
		"OptionalHeader32": reflect.TypeOf((*pkg.OptionalHeader32)(nil)).Elem(),
		"OptionalHeader64": reflect.TypeOf((*pkg.OptionalHeader64)(nil)).Elem(),
		"SectionHeader32":  reflect.TypeOf((*pkg.SectionHeader32)(nil)).Elem(),
		"Reloc":            reflect.TypeOf((*pkg.Reloc)(nil)).Elem(),
		"SectionHeader":    reflect.TypeOf((*pkg.SectionHeader)(nil)).Elem(),
		"Section":          reflect.TypeOf((*pkg.Section)(nil)).Elem(),
	})
}
