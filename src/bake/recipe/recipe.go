// Copyright 2013 Sean Kelleher. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

// Package recipe provides access to bake language extensions.
package recipe

import (
	"bake/recipe/test"
	"errors"
	"fmt"
	"os"
	"path"
)

const (
	bakeVar    = "BAKE"    // The name of the bake environment variable
	recipesDir = "recipes" // The directory containing bake recipes
)

type Recipe interface {
	Lang() string
}

func For(lang string) (Recipe, error) {
	return newRecipeFor(lang)
}

type recipe struct {
	lang string
	path string
}

func newRecipeFor(lang string) (*recipe, error) {
	bakePath := os.Getenv(bakeVar)
	if len(bakePath) == 0 {
		return nil, errors.New("bake environment variable not set")
	}

	recipesPath := path.Join(bakePath, recipesDir)
	if _, err := os.Stat(recipesPath); os.IsNotExist(err) {
		return nil, errors.New("bake root doesn't contain recipes dir")
	}

	langRecpPath := path.Join(recipesPath, lang)
	if _, err := os.Stat(langRecpPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("'%s' is not a valid language", lang)
	}

	return &recipe{lang, langRecpPath}, nil
}

func (r *recipe) Lang() string {
	return r.lang
}

func (r *recipe) Path() string {
	return r.path
}

func Test(lang string) (bool, error) {
	r, err := newRecipeFor(lang)
	if err != nil {
		return false, err
	}
	return test.TestRecipe(r.Lang(), r.Path())
}
