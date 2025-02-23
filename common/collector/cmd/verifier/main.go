package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"ivxv.ee/common/collector/command/exit"
	"ivxv.ee/common/collector/conf"
	"ivxv.ee/common/collector/errors"
	"ivxv.ee/common/collector/log"
	//ivxv:modules common/collector/container
)

func main() {
	code, err := verifierMain()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
	}
	os.Exit(code)
}

func verifierMain() (int, error) {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+os.Args[0]+` [options] <container file>

verifier uses the trust root given in order to verify a container's signatures
and return the signers and signing times.

The trust container must contain a single file. This is the trust configuration
which must be identified by the key "trust.yaml".

The container to be verified must have an extension corresponding to the
container type it is, e.g., foo.bdoc.

options:`)
		flag.PrintDefaults()
	}

	trust := flag.String("trust", "/etc/ivxv/trust.bdoc",
		"`path` to the trust container. Must have an extension corresponding to\n"+
			"the container type it is, e.g., trust.bdoc.\n")
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		return exit.Usage, nil
	}
	path := flag.Arg(0)

	// We do not want the verifier application to log anything, but it is
	// still assumed that the context has a logger. Use TestContext which
	// provides a test logger that does nothing.
	ctx := log.TestContext(context.Background())

	cfg, code, err := conf.New(ctx, *trust, "", "")
	if err != nil {
		return code, fmt.Errorf("failed to load trust root: %v", err)
	}

	c, err := cfg.Container.OpenFile(path)
	if err != nil {
		code = exit.DataErr
		if perr := errors.CausedBy(err, new(os.PathError)); perr != nil {
			if os.IsNotExist(perr) {
				code = exit.NoInput
			}
		}
		return code, fmt.Errorf("failed to open container: %v", err)
	}
	defer c.Close()

	for _, s := range c.Signatures() {
		pattern := regexp.MustCompile("[0-9]+")
		if pattern.FindString(s.Signer.Subject.CommonName) == "" {
			personalCode := strings.TrimPrefix(s.Signer.Subject.SerialNumber, "PNOEE-")
			fmt.Println(s.Signer.Subject.CommonName+","+personalCode, s.SigningTime.Format(time.RFC3339))
		} else {
			fmt.Println(s.Signer.Subject.CommonName, s.SigningTime.Format(time.RFC3339))
		}
	}

	return exit.OK, nil
}
