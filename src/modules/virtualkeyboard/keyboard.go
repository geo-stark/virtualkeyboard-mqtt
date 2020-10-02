package virtualkeyboard

// #include "keyboard.h"
import "C"

import "fmt"

var keys = map[string]C.int{
	"0": C.KEY_0,
	"1": C.KEY_1,
	"2": C.KEY_2,
	"3": C.KEY_3,
	"4": C.KEY_4,
	"5": C.KEY_5,
	"6": C.KEY_6,
	"7": C.KEY_7,
	"8": C.KEY_8,
	"9": C.KEY_9,

	"f1":  C.KEY_F1,
	"f2":  C.KEY_F2,
	"f3":  C.KEY_F3,
	"f4":  C.KEY_F4,
	"f5":  C.KEY_F5,
	"f6":  C.KEY_F6,
	"f7":  C.KEY_F7,
	"f8":  C.KEY_F8,
	"f9":  C.KEY_F9,
	"f10": C.KEY_F10,
	"f11": C.KEY_F11,
	"f12": C.KEY_F12,

	"q": C.KEY_Q,
	"w": C.KEY_W,
	"e": C.KEY_E,
	"r": C.KEY_R,
	"t": C.KEY_T,
	"y": C.KEY_Y,
	"u": C.KEY_U,
	"i": C.KEY_I,
	"o": C.KEY_O,
	"p": C.KEY_P,
	"a": C.KEY_A,
	"s": C.KEY_S,
	"d": C.KEY_D,
	"f": C.KEY_F,
	"g": C.KEY_G,
	"h": C.KEY_H,
	"j": C.KEY_J,
	"k": C.KEY_K,
	"l": C.KEY_L,
	"z": C.KEY_Z,
	"x": C.KEY_X,
	"c": C.KEY_C,
	"v": C.KEY_V,
	"b": C.KEY_B,
	"n": C.KEY_N,
	"m": C.KEY_M,

	" ":         C.KEY_SPACE,
	"tab":       C.KEY_TAB,
	"esc":       C.KEY_ESC,
	"enter":     C.KEY_ENTER,
	"sysrq":     C.KEY_SYSRQ,
	"backspace": C.KEY_BACKSPACE,

	"lwin":   C.KEY_LEFTMETA,
	"lctrl":  C.KEY_LEFTCTRL,
	"lalt":   C.KEY_LEFTALT,
	"lshift": C.KEY_LEFTSHIFT,

	"caps":   C.KEY_CAPSLOCK,
	"scroll": C.KEY_SCROLLLOCK,
	"num":    C.KEY_NUMLOCK,

	"mute":       C.KEY_MUTE,
	"volumedown": C.KEY_VOLUMEDOWN,
	"volumeup":   C.KEY_VOLUMEUP,
	"playpause":  C.KEY_PLAYPAUSE,

	"home":     C.KEY_HOME,
	"end":      C.KEY_END,
	"pageup":   C.KEY_PAGEUP,
	"pagedown": C.KEY_PAGEDOWN,
	"ins":      C.KEY_INSERT,
	"del":      C.KEY_DELETE,
	"up":       C.KEY_UP,
	"down":     C.KEY_DOWN,
	"left":     C.KEY_LEFT,
	"right":    C.KEY_RIGHT,

	"power":   C.KEY_POWER,
	"suspend": C.KEY_SUSPEND,

	"select":   C.KEY_SELECT,
	"back":     C.KEY_BACK,
	"forward":  C.KEY_FORWARD,
	"homepage": C.KEY_HOMEPAGE,
	"search":   C.KEY_SEARCH,

	// aliases
	"ctrl":  C.KEY_LEFTCTRL,
	"alt":   C.KEY_LEFTALT,
	"shift": C.KEY_LEFTSHIFT,
	"win":   C.KEY_LEFTMETA,
}

type Options struct {
	Vendor  int
	Product int
	Name    string
}

func OpenEx(opts *Options) error {
	if res := C.vkb_open(); res != 0 {
		return fmt.Errorf("errno %v", res)
	}

	for _, v := range keys {
		C.vkb_add_key(v)
	}

	if res := C.vkb_register(); res != 0 {
		return fmt.Errorf("errno %v", res)
	}
	return nil
}

func Open() error {
	return OpenEx(&Options{})
}

func Close() {
	C.vkb_close()
}

func Emit(items []string) error {
	list := make([]C.int, len(items))

	for index, item := range items {
		if item == "" {
			items[index] = ","
		}
		key, present := keys[item]
		if !present {
			return fmt.Errorf("emit unknown key '%v'", item)
		}
		list[index] = key
	}

	for _, key := range list {
		if C.vkb_emit_push(key, 1) != 0 || C.vkb_emit_sync() != 0 {
			return fmt.Errorf("emit error")
		}
	}
	for _, key := range list {
		if C.vkb_emit_push(key, 0) != 0 || C.vkb_emit_sync() != 0 {
			return fmt.Errorf("emit error")
		}
	}
	return nil
}
