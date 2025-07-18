package ra

import (
	"fmt"
)

type FlagSet struct {
	used       *bool
	flags      map[string]any
	positional []string
	subCmds    map[string]*FlagSet
}

func NewCmd() *FlagSet {
	return &FlagSet{
		flags:      make(map[string]any),
		positional: []string{},
		subCmds:    make(map[string]*FlagSet),
	}
}

func (fs *FlagSet) Parse(args []string) *ParseError {
	// todo
	//  remember to check if particular command invoked
	return nil
}

// BASE TYPES

type BaseFlag struct {
	Name     string
	Short    string
	Usage    string
	Optional bool
	Hidden   bool
}

type Flag[T any] struct {
	BaseFlag
	Default *T
	Value   *T
}

// NON-SLICE FLAGS

type BoolFlag struct {
	Flag[bool]
}
type StringFlag struct {
	Flag[string]
}

type IntFlag struct {
	Flag[int]
	min *int
	max *int
}

// SLICE FLAGS

type SliceFlag[T any] struct {
	BaseFlag
	Separator *string
	Variadic  bool
	Default   *[]T
	Value     *[]T
}

type StringSliceFlag = SliceFlag[string]
type IntSliceFlag = SliceFlag[int]

// NEW ARGS

func NewBool(name string) *BoolFlag {
	return &BoolFlag{Flag: Flag[bool]{BaseFlag: BaseFlag{Name: name}}}
}

func NewString(name string) *StringFlag {
	return &StringFlag{Flag: Flag[string]{BaseFlag: BaseFlag{Name: name}}}
}

func NewInt(name string) *IntFlag {
	return &IntFlag{Flag: Flag[int]{BaseFlag: BaseFlag{Name: name}}}
}

func NewStringSlice(name string) *StringSliceFlag {
	return &SliceFlag[string]{BaseFlag: BaseFlag{Name: name}}
}

func NewIntSlice(name string) *IntSliceFlag {
	return &SliceFlag[int]{BaseFlag: BaseFlag{Name: name}}
}

// NON-SLICE SETTERS

func (f BoolFlag) SetShort(s string) BoolFlag {
	f.Short = s
	return f
}

func (f BoolFlag) SetUsage(u string) BoolFlag {
	f.Usage = u
	return f
}

func (f BoolFlag) SetDefault(v bool) BoolFlag {
	f.Default = &v
	return f
}

func (f BoolFlag) SetOptional(b bool) BoolFlag {
	f.Optional = b
	return f
}

func (f BoolFlag) SetHidden(b bool) BoolFlag {
	f.Hidden = b
	return f
}

func (f StringFlag) SetShort(s string) StringFlag {
	f.Short = s
	return f
}

func (f StringFlag) SetUsage(u string) StringFlag {
	f.Usage = u
	return f
}

func (f StringFlag) SetDefault(v string) StringFlag {
	f.Default = &v
	return f
}

func (f StringFlag) SetOptional(b bool) StringFlag {
	f.Optional = b
	return f
}

func (f StringFlag) SetHidden(b bool) StringFlag {
	f.Hidden = b
	return f
}

func (f IntFlag) SetShort(s string) IntFlag {
	f.Short = s
	return f
}

func (f IntFlag) SetUsage(u string) IntFlag {
	f.Usage = u
	return f
}

func (f IntFlag) SetDefault(v int) IntFlag {
	f.Default = &v
	return f
}

func (f IntFlag) SetOptional(b bool) IntFlag {
	f.Optional = b
	return f
}

func (f IntFlag) SetHidden(b bool) IntFlag {
	f.Hidden = b
	return f
}

func (f IntFlag) SetMin(min int) IntFlag {
	f.min = &min
	return f
}

func (f IntFlag) SetMax(max int) IntFlag {
	f.max = &max
	return f
}

func (f SliceFlag[T]) SetShort(s string) SliceFlag[T] {
	f.Short = s
	return f
}

func (f SliceFlag[T]) SetUsage(u string) SliceFlag[T] {
	f.Usage = u
	return f
}

func (f SliceFlag[T]) SetDefault(v []T) SliceFlag[T] {
	f.Default = &v
	return f
}

func (f SliceFlag[T]) SetOptional(b bool) SliceFlag[T] {
	f.Optional = b
	return f
}

func (f SliceFlag[T]) SetHidden(b bool) SliceFlag[T] {
	f.Hidden = b
	return f
}

func (f SliceFlag[T]) SetSeparator(sep string) SliceFlag[T] {
	f.Separator = &sep
	return f
}

func (f SliceFlag[T]) SetVariadic(b bool) SliceFlag[T] {
	f.Variadic = b
	return f
}

// REGISTER METHODS

func (f Flag[T]) Register(fs *FlagSet) (*T, error) {
	if _, exists := fs.flags[f.Name]; exists {
		return nil, fmt.Errorf("flag %q already defined", f.Name)
	}

	valPtr := new(T)

	n := f // shallow copy
	if f.Default != nil {
		def := *f.Default // shallow copy to make default immutable
		n.Default = &def
	}
	n.Value = valPtr

	// register into the set
	fs.flags[f.Name] = n
	fs.positional = append(fs.positional, f.Name)

	return valPtr, nil
}

func (f SliceFlag[T]) Register(fs *FlagSet) (*[]T, error) {
	if _, exists := fs.flags[f.Name]; exists {
		return nil, fmt.Errorf("flag %q already defined", f.Name)
	}

	valPtr := new([]T)

	n := f // shallow copy
	if f.Default != nil {
		def := *f.Default // shallow copy to make default immutable
		n.Default = &def
	}

	n.Value = valPtr

	fs.flags[f.Name] = n
	fs.positional = append(fs.positional, f.Name)

	return valPtr, nil
}

func (fs *FlagSet) RegisterCmd(name string, cmd *FlagSet) (*bool, error) {
	if _, exists := fs.subCmds[name]; exists {
		return nil, fmt.Errorf("command %q already defined", name)
	}

	fs.subCmds[name] = cmd
	cmd.used = new(bool)
	return cmd.used, nil
}
