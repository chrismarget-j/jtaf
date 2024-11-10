package providerdata

import (
	"context"
	"encoding/xml"
	"fmt"
	"sync"

	"github.com/ChrisTrenkamp/xsel"
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

func (c Client) SetConfig(ctx context.Context, parentPath types.String, v any, diags *diag.Diagnostics) {
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

//func (c Client) GetConfig(ctx context.Context, path common.Path, diags *diag.Diagnostics) []byte {
//	ds := netconf.Candidate
//
//	c.sessionMutext.Lock()
//	defer c.sessionMutext.Unlock()
//
//	var err error
//	if !c.cacheOk {
//		c.cache, err = c.session.GetConfig(ctx, ds)
//		if err != nil {
//			diags.AddError(fmt.Sprintf("failed while reading device %s config", ds), err.Error())
//			return nil
//		}
//
//		c.cacheOk = true
//	}
//
//	var cbf exml.TagCallback
//	cbf = func(attrs exml.Attrs) {
//
//	}
//	_ = cbf
//
//	var tcb exml.TextCallback
//	tcb = func(data exml.CharData) {
//		x := data
//		_ = x
//	}
//	_ = tcb
//
//	slashes := path.DelimitedString("/")
//
//	decoder := exml.NewDecoder(bytes.NewReader(c.cache))
//	//decoder.On(slashes, func(attrs exml.Attrs) {
//	//	decoder.OnText(tcb)
//	//})
//	//decoder.OnText()
//	decoder.OnTextOf(slashes, func(data exml.CharData) {
//		a := data
//		_ = a
//	})
//	decoder.Run()
//
//	return nil
//}

func newClient(session *netconf.Session) Client {
	return Client{
		session:       session,
		sessionMutext: new(sync.Mutex),
	}
}
