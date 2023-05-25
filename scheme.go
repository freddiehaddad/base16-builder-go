package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

var bases = []string{
	"00", "01", "02", "03", "04", "05", "06", "07",
	"08", "09", "0A", "0B", "0C", "0D", "0E", "0F",
}

type scheme struct {
	Slug string `yaml:"-"`

	Scheme      string `yaml:"scheme"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`

	// Colors will hold all the "base*" variables.
	Colors map[string]color `yaml:",inline"`
}

func schemeFromFile(schemesFS fs.FS, fileName string) (*scheme, bool) {
	ret := &scheme{}
	ext := ""

	logger := log.WithField("file", fileName)

	if strings.HasSuffix(fileName, ".yaml") {
		ext = ".yaml"
	} else if strings.HasSuffix(fileName, ".yml") {
		ext = ".yml"
	} else {
		logger.Error("Only YAML files are supported")
		return nil, false
	}

	data, err := fs.ReadFile(schemesFS, fileName)
	if err != nil {
		logger.Error(err)
		return nil, false
	}

	err = yaml.Unmarshal(data, ret)
	if err != nil {
		logger.Error(err)
		return nil, false
	}

	// Now that we have the data, we can sanitize it
	ok := true
	if ret.Scheme == "" {
		logger.Error("Scheme name cannot be empty")
		ok = false
	}

	// Author is a warning because there appear to be some themes
	// without them.
	if ret.Author == "" {
		logger.Warn("Scheme author should not be empty")
	}

	if len(bases) != len(ret.Colors) {
		logger.Error("Wrong number of colors in scheme")
		ok = false
	}

	// Sanitize any fields which were added later
	if ret.Description == "" {
		ret.Description = ret.Scheme
	}

	// Now that we've got all that out of the way, we can start
	// processing stuff.

	// Take the last path component and chop off .yaml
	ret.Slug = filepath.Base(strings.TrimSuffix(fileName, ext))

	for _, base := range bases {
		baseKey := "base" + base
		if _, innerOk := ret.Colors[baseKey]; !innerOk {
			logger.Errorf("Scheme missing %q", baseKey)
			ok = false
			continue
		}
	}

	return ret, ok
}

func (s *scheme) mustacheContext() map[string]interface{} {
	ret := map[string]interface{}{
		"scheme-name":        s.Scheme,
		"scheme-author":      s.Author,
		"scheme-slug":        s.Slug,
		"scheme-description": s.Description,

		// Any extensions on the spec should go here
		"scheme-slug-underscored": strings.ReplaceAll(s.Slug, "-", "_"),
	}

	for _, base := range bases {
		baseKey := "base" + base
		baseVal := s.Colors[baseKey]

		// Note that we only lowercase the output of this to match the reference
		// repo.
		ret[baseKey+"-hex"] = fmt.Sprintf("%02x%02x%02x", baseVal.R, baseVal.G, baseVal.B)
		ret[baseKey+"-hex-bgr"] = fmt.Sprintf("%02x%02x%02x", baseVal.B, baseVal.G, baseVal.R)

		ret[baseKey+"-rgb-r"] = baseVal.R
		ret[baseKey+"-rgb-g"] = baseVal.G
		ret[baseKey+"-rgb-b"] = baseVal.B
		ret[baseKey+"-dec-r"] = float32(baseVal.R) / 255
		ret[baseKey+"-dec-g"] = float32(baseVal.G) / 255
		ret[baseKey+"-dec-b"] = float32(baseVal.B) / 255
		ret[baseKey+"-hex-r"] = fmt.Sprintf("%02x", baseVal.R)
		ret[baseKey+"-hex-g"] = fmt.Sprintf("%02x", baseVal.G)
		ret[baseKey+"-hex-b"] = fmt.Sprintf("%02x", baseVal.B)
	}

	return ret
}

func loadSchemes(schemesFS fs.FS) ([]*scheme, bool) {
	schemes := make(map[string]*scheme)

	schemePaths, err := fs.Glob(schemesFS, "*.y*")
	if err != nil {
		log.Error(err)
		return nil, false
	}

	for _, schemePath := range schemePaths {
		scheme, ok := schemeFromFile(schemesFS, schemePath)
		if !ok {
			log.Errorf("Failed to load scheme")
			continue
		}

		// XXX: this should never happen because it's now a single schemes dir,
		// but we include this check just in case someone messed something up.
		if _, ok := schemes[scheme.Slug]; ok {
			log.WithField("scheme", scheme.Slug).Warnf("Conflicting scheme")
		}

		log.Debugf("Found scheme %q", scheme.Slug)

		schemes[scheme.Slug] = scheme
	}

	ret := []*scheme{}
	for _, scheme := range schemes {
		ret = append(ret, scheme)
	}

	return ret, true
}
