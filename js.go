package gopherpc

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"path"
)

var (
	//go:embed gopherpc.js
	gopherpcJs        []byte
	gopherJsIntegrity string
	gopherpcJsName    string
)

func init() {
	sha256HashB := sha256.Sum256(gopherpcJs)
	gopherJsIntegrity = "sha256-" + base64.StdEncoding.EncodeToString(sha256HashB[:])

	gopherpcJsName = fmt.Sprintf("%x.gopherpc.js", sha256HashB)
}

func ImportJs() string {
	return fmt.Sprintf(
		`<script src="%s" integrity="%s" crossorigin="anonymous" async type="module"></script>`,
		path.Join("/", "__gopherpc__", gopherpcJsName),
		gopherJsIntegrity,
	)
}
