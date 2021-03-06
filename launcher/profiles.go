package launcher

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/brawaru/marct/sdtypes"
	"github.com/brawaru/marct/utils"
	"github.com/relvacode/iso8601"
)

func (w *Instance) ReadProfiles() (profiles *Profiles, err error) {
	err = unmarshalJSONFile(filepath.Join(w.Path, launcherProfilesPath), &profiles)

	return
}

func (w *Instance) WriteProfiles(profiles *Profiles) (err error) {
	if profiles == nil {
		return errors.New("cannot write null profiles")
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	err = enc.Encode(profiles)

	if err == nil {
		var file *os.File
		file, err = os.Create(filepath.Join(w.Path, launcherProfilesPath))

		if err == nil {
			defer utils.DClose(file)
			_, err = io.Copy(file, buf)
		}
	}

	return
}

func initDefaultProfiles() *Profiles {
	version := 3
	releaseIcon := "Grass"
	snapshotIcon := "Crafting_Table"
	defaultDate, _ := iso8601.ParseString("1970-01-01T00:00:00.000Z")
	return &Profiles{
		Profiles: map[string]Profile{
			utils.NewUUID(): {
				Type:          "latest-release",
				LastVersionID: "latest-release",
				Icon:          &releaseIcon,
				Created:       (*sdtypes.ISOTime)(&defaultDate),
				LastUsed:      sdtypes.ISOTime(defaultDate),
			},
			utils.NewUUID(): {
				Type:          "latest-snapshot",
				LastVersionID: "latest-snapshot",
				Icon:          &snapshotIcon,
				Created:       (*sdtypes.ISOTime)(&defaultDate),
				LastUsed:      sdtypes.ISOTime(defaultDate),
			},
		},
		Version: &version,
	}
}

func (w *Instance) ReadOrCreateProfiles() (profiles *Profiles, err error) {
	profiles, err = w.ReadProfiles()

	if err != nil {
		if utils.DoesNotExist(err) {
			profiles = initDefaultProfiles()
			err = nil
		}
	}

	return
}
