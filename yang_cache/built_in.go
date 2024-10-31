package yangcache

import (
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

		tf, err := os.CreateTemp(dn, "."+path.Base(fn))
		if err != nil {
			return nil, fmt.Errorf("while creating temporary file - %w", err)
		}
		tfn := tf.Name()

		_, err = tf.Write([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("while writing to temporary file %q - %w", tfn, err)
		}

		err = tf.Close()
		if err != nil {
			return nil, fmt.Errorf("while closing temporary file %q - %w", tfn, err)
		}

		err = os.Rename(tfn, fn)
		if err != nil {
			return nil, fmt.Errorf("while renaming temporary file %q to %q - %w", tfn, fn, err)
		}

		resultMap[dn] = struct{}{}
	}

	return helpers.Keys(resultMap), nil
}
