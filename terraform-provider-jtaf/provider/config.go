// Copyright (c) Juniper Networks, Inc., 2024-2024.
// All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"context"
	"fmt"
	"math"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nemith/netconf"
	ncssh "github.com/nemith/netconf/transport/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

const (
	envPrefix             = "JTAF_"
	knownHostsFile        = ".ssh/known_hosts"
	attrValidationSummary = "%s attribute is unset"
	attrValidationDetail  = "%[1]s attribute must be set in the provider configuration block, or by setting the " + envPrefix + "%[1]s environment variable"

	tfsdkHost                   = "host"
	tfsdkPort                   = "port"
	tfsdkUsername               = "username"
	tfsdkPassword               = "password"
	tfsdkSshPrivateKeyFile      = "ssh_private_key_file"
	tfsdkSshPrivateKeyString    = "ssh_private_key_string"
	tfsdkSshSkipKnownHostsCheck = "ssh_skip_known_host_check"
)

type config struct {
	Host                   types.String `tfsdk:"host"`
	Port                   types.Int64  `tfsdk:"port"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	SshPrivateKeyFile      types.String `tfsdk:"ssh_private_key_file"`
	SshPrivateKeyString    types.String `tfsdk:"ssh_private_key_string"`
	SshSkipKnownHostsCheck types.Bool   `tfsdk:"ssh_skip_known_host_check"`
}

func (p *config) attributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		tfsdkHost: schema.StringAttribute{
			Optional:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		tfsdkPort: schema.Int64Attribute{
			Optional:   true,
			Validators: []validator.Int64{int64validator.Between(1, math.MaxUint16)},
		},
		tfsdkUsername: schema.StringAttribute{
			Optional:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		tfsdkPassword: schema.StringAttribute{
			Optional:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		tfsdkSshPrivateKeyFile: schema.StringAttribute{
			Optional:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		tfsdkSshPrivateKeyString: schema.StringAttribute{
			Optional:   true,
			Validators: []validator.String{stringvalidator.LengthAtLeast(1)},
		},
		tfsdkSshSkipKnownHostsCheck: schema.BoolAttribute{
			Optional: true,
		},
	}
}

func (p *config) loadEnv(diags *diag.Diagnostics) {
	if s, ok := os.LookupEnv(envPrefix + tfsdkHost); ok && p.Host.IsNull() {
		p.Host = types.StringValue(s)
	}

	if s, ok := os.LookupEnv(envPrefix + tfsdkPort); ok && p.Port.IsNull() {
		i, err := strconv.Atoi(s)
		if err != nil {
			diags.AddError(fmt.Sprintf("environment variable %s (%q) failed to parse as integer", envPrefix+tfsdkPort, s), err.Error())
		}
		p.Port = types.Int64Value(int64(i))
	}

	if s, ok := os.LookupEnv(envPrefix + tfsdkUsername); ok && p.Username.IsNull() {
		p.Username = types.StringValue(s)
	}

	if s, ok := os.LookupEnv(envPrefix + tfsdkPassword); ok && p.Password.IsNull() {
		p.Password = types.StringValue(s)
	}

	if s, ok := os.LookupEnv(envPrefix + tfsdkSshPrivateKeyFile); ok && p.SshPrivateKeyFile.IsNull() {
		p.SshPrivateKeyFile = types.StringValue(s)
	}

	if s, ok := os.LookupEnv(envPrefix + tfsdkSshPrivateKeyString); ok && p.SshPrivateKeyString.IsNull() {
		p.SshPrivateKeyString = types.StringValue(s)
	}

	if s, ok := os.LookupEnv(envPrefix + tfsdkSshSkipKnownHostsCheck); ok && p.SshSkipKnownHostsCheck.IsNull() {
		b, err := strconv.ParseBool(s)
		if err != nil {
			diags.AddError(fmt.Sprintf("environment variable %s (%q) failed to parse as boolean", envPrefix+tfsdkSshSkipKnownHostsCheck, s), err.Error())
		}
		p.SshSkipKnownHostsCheck = types.BoolValue(b)
	}
}

func (p *config) validate(diags *diag.Diagnostics) {
	if p.Host.IsNull() || p.Host.ValueString() == "" {
		diags.AddError(fmt.Sprintf(attrValidationSummary, tfsdkHost), fmt.Sprintf(attrValidationDetail, tfsdkHost))
	}

	if p.Username.IsNull() || p.Username.ValueString() == "" {
		diags.AddError(fmt.Sprintf(attrValidationSummary, tfsdkUsername), fmt.Sprintf(attrValidationDetail, tfsdkUsername))
	}

	if p.Username.IsNull() || p.Username.ValueString() == "" {
		diags.AddError(fmt.Sprintf(attrValidationSummary, tfsdkUsername), fmt.Sprintf(attrValidationDetail, tfsdkUsername))
	}

	if (p.Password.IsNull() || p.Password.ValueString() == "") &&
		(p.SshPrivateKeyFile.IsNull() || p.SshPrivateKeyFile.ValueString() == "") &&
		(p.SshPrivateKeyString.IsNull() || p.SshPrivateKeyString.ValueString() == "") {
		diags.AddError(
			fmt.Sprintf("missing account authorization secret"),
			fmt.Sprintf("account authorization secret must be set by configuring at least one of: %[1]s, %[2]s or %[3]s "+
				"in the provider configuration block, or by setting at least one of the following environment variables: "+
				"%[4]s%[1]s, %[4]s%[2]s or %[4]s%[3]s ", tfsdkPassword, tfsdkSshPrivateKeyFile, tfsdkSshPrivateKeyString, envPrefix),
		)
	}
}

func (p *config) hostKeyCallback(diags *diag.Diagnostics) ssh.HostKeyCallback {
	if p.SshSkipKnownHostsCheck.ValueBool() {
		return ssh.InsecureIgnoreHostKey()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		diags.AddError("failed to determine user homedir", err.Error())
		return nil
	}

	hostKeyCallBack, err := knownhosts.New(path.Join(homeDir, knownHostsFile))
	if err != nil {
		diags.AddError("failed to create ssh host key callback function", err.Error())
		return nil
	}

	return hostKeyCallBack
}

func (p *config) session(ctx context.Context, diags *diag.Diagnostics) *netconf.Session {
	var authMethods []ssh.AuthMethod
	if !p.Password.IsNull() {
		authMethods = append(authMethods, ssh.Password(p.Password.ValueString()))
	}
	if !p.SshPrivateKeyString.IsNull() {
		signer, err := ssh.ParsePrivateKey([]byte(p.SshPrivateKeyString.ValueString()))
		if err != nil {
			diags.AddError("failed parsing configured private key string from environment or provider block", err.Error())
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	if !p.SshPrivateKeyFile.IsNull() {
		b, err := os.ReadFile(p.SshPrivateKeyFile.ValueString())
		if err != nil {
			diags.AddError(fmt.Sprintf("failed reading private key file %s", p.SshPrivateKeyFile), err.Error())
			return nil
		}

		signer, err := ssh.ParsePrivateKey(b)
		if err != nil {
			diags.AddError(fmt.Sprintf("failed parsing private key data from file %s", p.SshPrivateKeyFile), err.Error())
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	hostKeyCallback := p.hostKeyCallback(diags)
	if diags.HasError() {
		return nil
	}

	sshConfig := ssh.ClientConfig{
		User:            p.Username.ValueString(),
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
	}

	sshAddr := strings.Join([]string{p.Host.ValueString(), strconv.Itoa(int(p.Port.ValueInt64()))}, ":")

	transport, err := ncssh.Dial(ctx, "tcp", sshAddr, &sshConfig)
	if err != nil {
		diags.AddError("failed to set up ssh transport with "+sshAddr, err.Error())
		return nil
	}

	session, err := netconf.Open(transport)
	if err != nil {
		diags.AddError("failed to open netconf session with "+sshAddr, err.Error())
		return nil
	}

	return session
}
