// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

type Contact struct {
	Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Company   string `json:"company,omitempty"`
	Subject   string `json:"subject,omitempty"`
	Message   string `json:"message"`
	Tag       string `json:"tag"`
}

func NewContact(opts ...Option) *Contact {
	p := &Contact{
		Model: NewModel(opts...),
	}
	return p
}

type ContactRequest struct {
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name"  validate:"required"`
	Email           string `json:"email"  validate:"required,email"`
	Message         string `json:"message"  validate:"required"`
	Company         string `json:"company,omitempty"  validate:"omitempty"`
	Phone           string `json:"phone,omitempty"  validate:"omitempty"`
	Subject         string `json:"subject,omitempty"  validate:"omitempty"`
	CaptchaResponse string `json:"captcha_response,omitempty"`
}

func (p *ContactRequest) Contact(tag string, opts ...Option) *Contact {
	c := NewContact(opts...)
	c.FirstName = p.FirstName
	c.LastName = p.LastName
	c.Email = p.Email
	c.Message = p.Message
	c.Phone = p.Phone
	c.Company = p.Company
	c.Subject = p.Subject
	c.Tag = tag

	return c
}
