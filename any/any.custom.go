/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package any

func NewAny(data []byte) (*RawData, error) {
	return &RawData{Data: data}, nil
}

func (a *RawData) MarshalJSON() ([]byte, error) {
	if len(a.Data) == 0 {
		return []byte("null"), nil
	}
	return a.Data, nil
}

func (a *RawData) Size() int {
	return len(a.Data)
}

func (a *RawData) MarshalTo(b []byte) (int, error) {
	return copy(b, a.Data), nil
}

func (a *RawData) Unmarshal(b []byte) error {
	a.Data = b
	return nil
}

func (a *RawData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(a.Data)
	copy(dAtA[i:], a.Data)
	return len(a.Data), nil
}
