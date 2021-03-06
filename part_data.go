package gen

import (
	"bufio"
	"strconv"
	"strings"
)

type PartData struct {
	Thickness          float64
	Quality            string
	Name               string
	PosNo              string
	PartnameLong       string
	ShipNo             string
	AssemblyLow        string
	AssemblyHigh       string
	BlockNo            string
	PartArea           float64
	Mirrored           bool
	ExtensionU         float64
	ExtensionV         float64
	PlateSide          string
	PartCog            Point
	NoIntervalsExcess1 int
	LengthExcess1      float64
	NoIntervalsExcess2 int
	LengthExcess2      float64
	NoIntervalsExcess3 int
	LengthExcess3      float64
	StringData         []StringData
	BurningData        []BurningData
	MarkingData        []MarkingData
	GeometryData       []GeometryData
	Quantity           int
}

func readPartData(s *bufio.Scanner) PartData {
	var p PartData
	var partCog Point

next:
	for s.Scan() {
		k, v, ok := strings.Cut(s.Text(), "=")

		if k == "TRANSFORMATION_DATA" {
			for s.Scan() {
				switch s.Text() {
				case "END_OF_TRANSFORMATION_DATA":
					continue next
				}
			}
		}

		if !ok {
			break
		}

		switch k {
		case "NAME":
			p.Name = v
		case "POSNO":
			p.PosNo = v
		case "PARTNAME_LONG", "PARTNAME_LONG_SB":
			p.PartnameLong = v
		case "SHIP_NO":
			p.ShipNo = v
		case "ASSEMBLY_LOW", "ASSEMBLY_LOW_SB":
			p.AssemblyLow = v
		case "ASSEMBLY_HIGH", "ASSEMBLY_HIGH_SB":
			p.AssemblyHigh = v
		case "BLOCK_NO":
			p.BlockNo = v
		case "PART_AREA":
			p.PartArea, _ = strconv.ParseFloat(v, 64)
		case "MIRRORED":
			p.Mirrored, _ = strconv.ParseBool(v)
		case "EXTENSION_U":
			p.ExtensionU, _ = strconv.ParseFloat(v, 64)
		case "EXTENSION_V":
			p.ExtensionV, _ = strconv.ParseFloat(v, 64)
		case "PLATE_SIDE":
			p.PlateSide = v
		case "PART_COG_U":
			partCog.X, _ = strconv.ParseFloat(v, 64)
		case "PART_COG_V":
			partCog.Y, _ = strconv.ParseFloat(v, 64)
		case "NO_INTERVALS_EXCESS_1":
			p.NoIntervalsExcess1, _ = strconv.Atoi(v)
		case "LENGTH_EXCESS_1":
			p.LengthExcess1, _ = strconv.ParseFloat(v, 64)
		case "NO_INTERVALS_EXCESS_2":
			p.NoIntervalsExcess2, _ = strconv.Atoi(v)
		case "LENGTH_EXCESS_2":
			p.LengthExcess2, _ = strconv.ParseFloat(v, 64)
		case "NO_INTERVALS_EXCESS_3":
			p.NoIntervalsExcess3, _ = strconv.Atoi(v)
		case "LENGTH_EXCESS_3":
			p.LengthExcess3, _ = strconv.ParseFloat(v, 64)
		}
	}

	p.PartCog = partCog
	p.Quantity = 1
	return p
}

func (p *PartData) Mirror() {
	p.Mirrored = !p.Mirrored

	p.PartCog.X = -p.PartCog.X

	for _, s := range p.StringData {
		s.Pos.X = -s.Pos.X
	}

	for _, b := range p.BurningData {
		b.Contour.invert()
		if b.GeometryData != nil {
			b.GeometryData.Contour.invert()
		}
	}

	for _, m := range p.MarkingData {
		if m.MarkingSide == "TS" {
			m.MarkingSide = "OS"
		} else {
			m.MarkingSide = "TS"
		}
		m.Contour.invert()
	}

	for _, g := range p.GeometryData {
		g.Contour.invert()
	}
}
