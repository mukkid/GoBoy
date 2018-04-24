package main

type GBMem struct {
	/* Work RAM at 0xc000 - 0xd000 */
	wram [8 * 1024]uint8
	vram [8 * 1024]uint8
	/* HRAM: 0xff80 - 0xfffe */
	hram [127]uint8
	/* ROM bank 0, nonswitchable - I believe this means this bank is static */
	cartridge GBCartridge
}

/* Cartridge type specified at 0x0147 */
type GBCartridgeType uint8

const (
	GBCartridgeROM GBCartridgeType = 0x00
	GBCartridgeMBC1
	GBCartridgeMBC1RAM
	GBCartridgeMBC1RAMBattery
	/* 0x04 unused */
	GBCartridgeMBC2 = 0x05
	GBCartridgeMBC2RAMBattery
	/* 0x07 unused */
	GBCartridgeROMRAM = 0x08
	GBCartridgeROMRAMBattery
	/* 0x0a unused */
	GBCartridgeMMM01 = 0x0b
	GBCartridgeMMM01RAM
	GBCartridgeMMM01RAMBattery
	/* 0x0e unused */
	GBCartridgeMBC3TimerBattery = 0x0f
	GBCartridgeMBC3RAMTimerBattery
	GBCartridgeMBC3
	GBCartridgeMBC3RAM
	GBCartridgeMBC3RAMBattery
	/* 0x14 - 0x18 unused */
	GBCartridgeMBC5 = 0x19
	GBCartridgeMBC5RAM
	GBCartridgeMBC5RAMBattery
	GBCartridgeMBC5Rumble
	GBCartridgeMBC5RAMRumble
	GBCartridgeMBC5RAMBatteryRumble
	/* 0x1f unused */
	GBCartridgeMBC6RAMBattery = 0x20
	/* 0x21 unused */
	GBCartridgeMBC7RAMBatAccel = 0x22
	/* 0x23 - 0xfb unused */
	GBCartridgePocketCamera = 0xfc
	GBCartridgeBandaiTAMA5
	GBCartridgeHuC3
	GBCartridgeHuC1RAMBattery
)

/*
 * ROM Size
 * 0x0148 indicates the size of ROM, computed as 32KB << n
 * Or 0x8000 < n bytes
 */
func (m *GBMem) read(addr uint16) uint8 {
	if addr >= 0x0000 && addr < 0x8000 {
		return m.cartridge.readROM(addr)
	} else if addr >= 0x8000 && addr < 0xa000 {
		/* VRAM */
		return m.vram[addr-0x8000]
	} else if addr >= 0xa000 && addr < 0xc000 {
		/* SRAM External RAM in cartridge, often battery buffered */
		return m.cartridge.readRAM(addr)
	} else if addr >= 0xc000 && addr < 0xd000 {
		/* WRAM0 Work RAM */
		return m.wram[addr-0xc000]
	} else if addr >= 0xd000 && addr < 0xe000 {
		/*
		 * WRAMX, switchable (1-7) in GBC mode
		 * TODO: Implement GBC mode
		 */
		return m.wram[addr-0xc000]
	} else if addr >= 0xe000 && addr < 0xfe00 {
		/* ECHO of 0xc000 - 0xde00 */
		return m.wram[addr-0x2000-0xc000]
	} else if addr >= 0xfe00 && addr < 0xfea0 {
		/* OAM (Object Attribute Table) Sprite information table */
		return uint8(0x00)
	} else if addr >= 0xfea0 && addr < 0xff00 {
		/* Unused */
		return uint8(0x00)
	} else if addr >= 0xff00 && addr < 0xff80 {
		/* I/O Registers I/O registers are mapped here */
		return uint8(0x00)
	} else if addr >= 0xff80 && addr < 0xffff {
		/* HRAM Internal CPU RAM */
		return m.hram[addr-0xff80]
	} else {
		/* 0xffff - IE Register Interrupt enable flags */
		return uint8(0x00)
	}
}

func (m *GBMem) write(addr uint16, value uint8) {
	if addr >= 0x0000 && addr < 0x8000 {
		/*
		 * Both non-switchable and switchable ROM Bank.
		 */
		m.cartridge.writeROM(addr, value)
	} else if addr >= 0x8000 && addr < 0xa000 {
		/* VRAM */
		m.vram[addr-0x8000] = value
	} else if addr >= 0xa000 && addr < 0xc000 {
		/* SRAM External RAM in cartridge, often battery buffered */
		m.cartridge.writeRAM(addr, value)
	} else if addr >= 0xc000 && addr < 0xd000 {
		/* WRAM0 Work RAM */
		m.wram[addr-0xc000] = value
	} else if addr >= 0xd000 && addr < 0xe000 {
		/*
		 * TODO: Implement GBC mode
		 * WRAMX, switchable (1-7) in GBC mode
		 */
		m.wram[addr-0xc000] = value
	} else if addr >= 0xe000 && addr < 0xfe00 {
		/* ECHO of 0xc000 - 0xde00 */
		m.wram[addr-0x2000-0xc000] = value
	} else if addr >= 0xfe00 && addr < 0xfea0 {
		/* OAM (Object Attribute Table) Sprite information table */
	} else if addr >= 0xfea0 && addr < 0xff00 {
		/* Unused */
	} else if addr >= 0xff00 && addr < 0xff80 {
		/* I/O Registers I/O registers are mapped here */
	} else if addr >= 0xff80 && addr < 0xffff {
		/* HRAM Internal CPU RAM */
		m.hram[addr-0xff80] = value
	} else {
		/* 0xffff - IE Register Interrupt enable flags */
	}
}

func (m *GBMem) loadROM(data []uint8) {
	m.cartridge.loadROM(data)
}

/* Interface for the different cartridges */
type GBCartridge interface {
	/*
	 * ROM  - [0x4000, 0x8000)
	 * VRAM - [0x8000, 0x9FFF)
	 * RAM  - [0xa000, 0xc000)
	 */
	loadROM(data []uint8) error
	readROM(addr uint16) uint8
	readRAM(addr uint16) uint8
	writeROM(addr uint16, data uint8)
	writeRAM(addr uint16, data uint8)
}
