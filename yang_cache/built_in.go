package yangcache

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"

	jtafCfg "github.com/chrismarget-j/jtaf/config"
	"github.com/chrismarget-j/jtaf/data/yang"
	"github.com/chrismarget-j/jtaf/helpers"
)

// populateBakedIn drops the baked-in yang files
// (data/yang/publisher/*.yang) into the yang cache dir. The
// returned strings represent yang directories the caller might
// use as a module source.
func populateBakedIn(cfg jtafCfg.Cfg) ([]string, error) {
	resultMap := make(map[string]struct{})
	for k, v := range yang.Files {
		fn := path.Join(cfg.YangCacheDir(), k)
		dn := path.Dir(fn)

		err := os.MkdirAll(dn, 0o755)
		if err != nil {
			return nil, fmt.Errorf("while mkdir-ing %q - %w", dn, err)
		}

		_, err = os.Stat(fn)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("while stat-ing %q - %w", fn, err)
		}

		if err == nil {
			// file exists. do checksum and end the loop iteration.
			h := sha256.New()
			_, err = h.Write([]byte(v))
			if err != nil {
				return nil, fmt.Errorf("while writing %q data to sha256 hasher - %w", fn, err)
			}

			err = helpers.CheckSha256(fn, fmt.Sprintf("%x", h.Sum(nil)))
			if err != nil {
				return nil, fmt.Errorf("cached file %q failed checksum; expected %q", fn, fmt.Sprintf("%x", h.Sum(nil)))
			}

			resultMap[dn] = struct{}{}

			continue // checksum is okay - nothing more to do with this one
		}

		// the file does not exist. must be created

		f, err := helpers.NewTmpFileWithRenameOnClose(dn, "."+path.Base(fn), fn)
		if err != nil {
			return nil, fmt.Errorf("while creating temp file in %q - %w", dn, err)
		}

		_, err = f.Write([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("while writing to temporary file - %w", err)
		}

		err = f.Close()
		if err != nil {
			return nil, fmt.Errorf("while closing temporary file - %w", err)
		}

		resultMap[dn] = struct{}{}
	}

	return helpers.Keys(resultMap), nil
}
