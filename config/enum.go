package jtafCfg

import "github.com/orsinium-labs/enum"

type junosPlatform enum.Member[string]

type osFamily enum.Member[string]

var (
	osFamilyJunos    = osFamily{Value: "junos"}
	osFamilyJunosES  = osFamily{Value: "junos-es"}
	osFamilyJunosQFX = osFamily{Value: "junos-qfx"}
	osFamilies       = enum.New(osFamilyJunos, osFamilyJunosES, osFamilyJunosQFX)
)
