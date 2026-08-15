package aws

import (
	"bytes"
	"os"
	"path/filepath"
)

type AwsConfig struct {
	Profiles []Profile
}

func (c *AwsConfig) processBlock(b []byte) {
	profile := profileFromBytes(b)
	if profile == nil {
		return
	}
	c.Profiles = append(c.Profiles, *profile)
}

func GetAwsConfig(dir string) (*AwsConfig, error) {
	file := filepath.Join(dir, "config")

	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		// no .aws/config present
		return &AwsConfig{
			Profiles: nil,
		}, nil
	} else if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &AwsConfig{}
	blocks := bytes.SplitSeq(data, []byte("["))
	for block := range blocks {
		if len(block) == 0 {
			continue
		}
		cfg.processBlock(block)
	}

	return cfg, nil
}
