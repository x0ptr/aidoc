package storage

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Cache map[string]map[string]Answer

type Answer struct {
	Short   string
	Verbose string
}

var (
	cachePath string = ""
	cache     Cache  = Cache{}
)

func CacheFilePath() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, "aidoc")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	cachePath = filepath.Join(dir, "cache.yaml")
	return cachePath, nil
}

func LoadCache() (Cache, error) {
	b, err := os.ReadFile(cachePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cache, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(b, &cache); err != nil {
		return nil, err
	}
	if cache == nil {
		cache = Cache{}
	}
	return cache, nil
}

func SaveCache() error {
	b, err := yaml.Marshal(cache)
	if err != nil {
		return err
	}
	tmp := cachePath + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, cachePath)
}

func ClearCache() error {
	path, err := CacheFilePath()
	if err != nil {
		return err
	}

	return os.Remove(path)
}

func CacheGet(lang, topic string, isVerbose bool) (string, bool) {
	if m, ok := cache[lang]; ok {
		if v, ok := m[topic]; ok {
			if isVerbose {
				return isEmpty(v.Verbose)
			} else {
				return isEmpty(v.Short)
			}
		}
	}
	return "", false
}

func CacheSet(lang, topic, val string, isVerbose bool) {
	_, ok := cache[lang]
	if !ok {
		cache[lang] = make(map[string]Answer)
		cache[lang][topic] = Answer{Verbose: "", Short: ""}
	}
	ans := Answer{
		Short:   cache[lang][topic].Short,
		Verbose: cache[lang][topic].Verbose,
	}
	if isVerbose {
		ans.Verbose = val
	} else {
		ans.Short = val
	}
	cache[lang][topic] = ans
}

func isEmpty(ans string) (string, bool) {
	return ans, len(ans) > 5
}
