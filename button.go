// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"github.com/wangch/win"
)

type clickable interface {
	raiseClicked()
}

type setCheckeder interface {
	setChecked(checked bool)
}

type Button struct {
	WidgetBase
	checkedChangedPublisher EventPublisher
	clickedPublisher        EventPublisher
	textChangedPublisher    EventPublisher
	image                   Image
	dbclickPublisher        EventPublisher
}

func NewButton(parent Container) (*Button, error) {
	btn := &Button{}
	if err := InitWidget(
		btn,
		parent,
		"BUTTON",
		win.WS_TABSTOP|win.WS_VISIBLE,
		0); err != nil {
		return nil, err
	}
	btn.init()
	return btn, nil
}

func (b *Button) DbClicked() *Event {
	return b.dbclickPublisher.Event()
}

func (b *Button) init() {
	b.MustRegisterProperty("Checked", NewBoolProperty(
		func() bool {
			return b.Checked()
		},
		func(v bool) error {
			b.SetChecked(v)
			return nil
		},
		b.CheckedChanged()))

	b.MustRegisterProperty("Text", NewProperty(
		func() interface{} {
			return b.Text()
		},
		func(v interface{}) error {
			return b.SetText(v.(string))
		},
		b.textChangedPublisher.Event()))
}

func (b *Button) Image() Image {
	return b.image
}

func (b *Button) SetImage(image Image) error {
	var typ uintptr
	var handle uintptr
	switch img := image.(type) {
	case nil:
		// zeroes are good

	case *Bitmap:
		typ = win.IMAGE_BITMAP
		handle = uintptr(img.hBmp)

	case *Icon:
		typ = win.IMAGE_ICON
		handle = uintptr(img.hIcon)

	default:
		return newError("image must be either *walk.Bitmap or *walk.Icon")
	}

	b.SendMessage(win.BM_SETIMAGE, typ, handle)

	b.image = image

	return b.updateParentLayout()
}

func (b *Button) Text() string {
	return windowText(b.hWnd)
}

func (b *Button) SetText(value string) error {
	if value == b.Text() {
		return nil
	}

	if err := setWindowText(b.hWnd, value); err != nil {
		return err
	}

	return b.updateParentLayout()
}

func (b *Button) Checked() bool {
	return b.SendMessage(win.BM_GETCHECK, 0, 0) == win.BST_CHECKED
}

func (b *Button) SetChecked(checked bool) {
	if checked == b.Checked() {
		return
	}

	b.window.(setCheckeder).setChecked(checked)
}

func (b *Button) setChecked(checked bool) {
	var chk uintptr

	if checked {
		chk = win.BST_CHECKED
	} else {
		chk = win.BST_UNCHECKED
	}

	b.SendMessage(win.BM_SETCHECK, chk, 0)

	b.checkedChangedPublisher.Publish()
}

func (b *Button) CheckedChanged() *Event {
	return b.checkedChangedPublisher.Event()
}

func (b *Button) Clicked() *Event {
	return b.clickedPublisher.Event()
}

func (b *Button) raiseClicked() {
	b.clickedPublisher.Publish()
}

func (b *Button) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_COMMAND:
		switch win.HIWORD(uint32(wParam)) {
		case win.BN_CLICKED:
			b.raiseClicked()
		}

	case win.WM_LBUTTONDBLCLK:
		b.dbclickPublisher.Publish()

	case win.WM_SETTEXT:
		b.textChangedPublisher.Publish()
	}

	return b.WidgetBase.WndProc(hwnd, msg, wParam, lParam)
}
