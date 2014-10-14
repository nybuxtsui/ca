package cmd

import (
	"fmt"
	"os"

	"github.com/nybuxtsui/etcd-ca/third_party/github.com/codegangsta/cli"

	"github.com/nybuxtsui/etcd-ca/depot"
	"github.com/nybuxtsui/etcd-ca/pkix"
)

func NewSignCommand() cli.Command {
	return cli.Command{
		Name:        "sign",
		Usage:       "Sign certificate request",
		Description: "Sign certificate request with CA, and generate certificate for the host.",
		Flags: []cli.Flag{
			cli.StringFlag{"passphrase", "", "Passphrase to decrypt private-key PEM block of CA"},
		},
		Action: newSignAction,
	}
}

func newSignAction(c *cli.Context) {
	if len(c.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "One host name must be provided.")
		os.Exit(1)
	}
	name := c.Args()[0]

	if depot.CheckCertificateHost(d, name) {
		fmt.Fprintln(os.Stderr, "Certificate has existed!")
		os.Exit(1)
	}

	csr, err := depot.GetCertificateSigningRequest(d, name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get certificate request error:", err)
		os.Exit(1)
	}
	crt, err := depot.GetCertificateAuthority(d)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get CA certificate error:", err)
		os.Exit(1)
	}
	info, err := depot.GetCertificateAuthorityInfo(d)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get CA certificate info error:", err)
		os.Exit(1)
	}
	key, err := depot.GetEncryptedPrivateKeyAuthority(d, getPassPhrase(c, "CA key"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Get CA key error:", err)
		os.Exit(1)
	}

	crtHost, err := pkix.CreateCertificateHost(crt, info, key, csr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Create certificate error:", err)
		os.Exit(1)
	} else {
		fmt.Printf("Created %s/crt from %s/csr signed by ca/key\n", name, name)
	}

	if err = depot.PutCertificateHost(d, name, crtHost); err != nil {
		fmt.Fprintln(os.Stderr, "Save certificate error:", err)
	}
	if err = depot.UpdateCertificateAuthorityInfo(d, info); err != nil {
		fmt.Fprintln(os.Stderr, "Update CA info error:", err)
	}
}
