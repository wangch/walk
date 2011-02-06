// Copyright 2010 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package walk

import (
	"os"
)

import (
	. "walk/winapi/user32"
)

type widgetListObserver interface {
	onInsertingWidget(index int, widget IWidget) os.Error
	onInsertedWidget(index int, widget IWidget) os.Error
	onRemovingWidget(index int, widget IWidget) os.Error
	onRemovedWidget(index int, widget IWidget) os.Error
	onClearingWidgets() os.Error
	onClearedWidgets() os.Error
}

type WidgetList struct {
	items    []IWidget
	observer widgetListObserver
}

func newWidgetList(observer widgetListObserver) *WidgetList {
	return &WidgetList{observer: observer}
}

func (l *WidgetList) Add(item IWidget) os.Error {
	return l.Insert(len(l.items), item)
}

func (l *WidgetList) At(index int) IWidget {
	return l.items[index]
}

func (l *WidgetList) Clear() os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onClearingWidgets(); err != nil {
			return err
		}
	}

	oldItems := l.items
	l.items = l.items[:0]

	if observer != nil {
		if err := observer.onClearedWidgets(); err != nil {
			l.items = oldItems
			return err
		}
	}

	return nil
}

func (l *WidgetList) Index(item IWidget) int {
	for i, widget := range l.items {
		if widget == item {
			return i
		}
	}

	return -1
}

func (l *WidgetList) Contains(item IWidget) bool {
	return l.Index(item) > -1
}

func (l *WidgetList) indexHandle(handle HWND) int {
	for i, widget := range l.items {
		if widget.BaseWidget().hWnd == handle {
			return i
		}
	}

	return -1
}

func (l *WidgetList) containsHandle(handle HWND) bool {
	return l.indexHandle(handle) > -1
}

func (l *WidgetList) insertIntoSlice(index int, item IWidget) {
	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = item
}

func (l *WidgetList) Insert(index int, item IWidget) os.Error {
	observer := l.observer
	if observer != nil {
		if err := observer.onInsertingWidget(index, item); err != nil {
			return err
		}
	}

	l.insertIntoSlice(index, item)

	if observer != nil {
		if err := observer.onInsertedWidget(index, item); err != nil {
			l.items = append(l.items[:index], l.items[index+1:]...)
			return err
		}
	}

	return nil
}

func (l *WidgetList) Len() int {
	return len(l.items)
}

func (l *WidgetList) Remove(item IWidget) os.Error {
	index := l.Index(item)
	if index == -1 {
		return nil
	}

	return l.RemoveAt(index)
}

func (l *WidgetList) RemoveAt(index int) os.Error {
	observer := l.observer
	item := l.items[index]
	if observer != nil {
		if err := observer.onRemovingWidget(index, item); err != nil {
			return err
		}
	}

	l.items = append(l.items[:index], l.items[index+1:]...)

	if observer != nil {
		if err := observer.onRemovedWidget(index, item); err != nil {
			l.insertIntoSlice(index, item)
			return err
		}
	}

	return nil
}