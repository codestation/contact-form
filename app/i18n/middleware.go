// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package i18n

import (
	"context"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type LanguageTagKey struct{}

func (l LanguageTagKey) String() string {
	return "LanguageTag"
}

func LoadMessagePrinter(preferLangKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, _ := c.Get(preferLangKey).(string)
			pref := c.Request().Header.Get("Accept-Language")
			lang := c.QueryParam("lang")
			tags := message.MatchLanguage(lang, user, pref)
			c.Set(LanguageTagKey{}.String(), tags)

			schemaCtx := context.WithValue(c.Request().Context(), LanguageTagKey{}, tags)
			request := c.Request().WithContext(schemaCtx)
			c.SetRequest(request)

			return next(c)
		}
	}
}

func GetLanguageTags(c echo.Context) language.Tag {
	if tags, ok := c.Get(LanguageTagKey{}.String()).(language.Tag); ok {
		return tags
	}
	return language.English
}

func GetLanguageTagsContext(ctx context.Context) language.Tag {
	if tags, ok := ctx.Value(LanguageTagKey{}).(language.Tag); ok {
		return tags
	}
	return language.English
}
