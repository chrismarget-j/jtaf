package providerdata

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"sync"

	"github.com/ChrisTrenkamp/xsel"
	"github.com/antchfx/xmlquery"
	"github.com/chrismarget-j/jtaf/terraform-provider-jtaf/common"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/nemith/netconf"
)

type Client struct {
	session       *netconf.Session
	sessionMutext *sync.Mutex
	cacheOk       bool
	cache         []byte
}

func (c *Client) SetConfig(ctx context.Context, parentPath types.String, v any, diags *diag.Diagnostics) {
	ds := netconf.Candidate

	pp, err := xsel.BuildExpr(parentPath.ValueString())
	if err != nil {
		diags.AddAttributeError(path.Root("parent_path"), "failed to parse xpath", err.Error())
		return
	}

	header, prefix, footer := common.XmlWrappersFromPathGrammar(pp, "", common.XmlIndent, diags)
	if diags.HasError() {
		return
	}

	xmlBytes, err := xml.MarshalIndent(v, prefix, common.XmlIndent)
	if err != nil {
		diags.AddError("failed marshaling config xml", err.Error())
		return
	}

	payload := header + string(xmlBytes) + "\n" + footer

	c.sessionMutext.Lock()
	defer c.sessionMutext.Unlock()

	c.cacheOk = false

	err = c.session.EditConfig(ctx, ds, []byte(payload))
	if err != nil {
		diags.AddError(fmt.Sprintf("failed while editing device %s config", ds), err.Error())
		return
	}
}

func (c *Client) GetConfig(ctx context.Context, path types.String, diags *diag.Diagnostics) []byte {
	ds := netconf.Candidate

	c.sessionMutext.Lock()
	defer c.sessionMutext.Unlock()

	if !c.cacheOk {
		var err error
		c.cache, err = c.session.GetConfig(ctx, ds)
		if err != nil {
			diags.AddError(fmt.Sprintf("failed while reading device %s config", ds), err.Error())
			return nil
		}

		c.cacheOk = true
	}

	cfg, err := xmlquery.Parse(bytes.NewReader(c.cache))
	if err != nil {
		diags.AddError("failed to parse device config", err.Error())
		return nil
	}

	nodes := xmlquery.Find(cfg, path.ValueString())
	switch len(nodes) {
	case 0:
		return nil
	case 1:
	default:
		diags.AddError(
			fmt.Sprintf("Expected 0-1 matches from device config, found %d", len(nodes)),
			fmt.Sprintf("xpath query: %q", path.ValueString()),
		)
	}

	return []byte(nodes[0].OutputXML(true))
}

func newClient(session *netconf.Session) Client {
	return Client{
		session:       session,
		sessionMutext: new(sync.Mutex),
	}
}
