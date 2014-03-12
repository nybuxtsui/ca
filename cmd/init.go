package cmd

import (
	"fmt"
	"os"

	"github.com/coreos/etcd-ca/third_party/github.com/codegangsta/cli"

	"github.com/coreos/etcd-ca/depot"
	"github.com/coreos/etcd-ca/pkix"
)

func NewInitCommand() cli.Command {
	return cli.Command{
		Name:        "init",
		Usage:       "Create Certificate Authority",
		Description: "Create Certificate Authority, including certificate, key and extra information file.",
		Action:      initAction,
	}
}

func initAction(c *cli.Context) {
	if depot.CheckCertificateAuthority(d) || depot.CheckCertificateAuthorityInfo(d) || depot.CheckPrivateKeyAuthority(d) {
		fmt.Fprintln(os.Stderr, "CA has existed!")
		os.Exit(1)
	}

	key, err := pkix.CreateRSAKey()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Created RSA Key error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Created ca/key")
	}

	crt, info, err := pkix.CreateCertificateAuthority(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Created certificate error:", err)
		os.Exit(1)
	} else {
		fmt.Println("Created ca/crt")
	}

	if err = depot.PutCertificateAuthority(d, crt); err != nil {
		fmt.Fprintln(os.Stderr, "Saved certificate error:", err)
	}
	if err = depot.PutCertificateAuthorityInfo(d, info); err != nil {
		fmt.Fprintln(os.Stderr, "Saved certificate info error:", err)
	}
	if err = depot.PutPrivateKeyAuthority(d, key); err != nil {
		fmt.Fprintln(os.Stderr, "Saved key error:", err)
	}
}
